package registry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"wallforfry.fr/harbor/configuration"
)

type Registry struct {
	url           string
	checkTLS      bool
	request       *gorequest.SuperAgent
	configuration configuration.Configuration
	language      configuration.Language
}

type Catalog struct {
	Repositories []string `json:"repositories"`
}

type Repository struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []Layer `json:"layers"`
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type Image struct {
	Registry     string
	Size         int
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Architecture string `json:"architecture"`
	Digest       string
	TagV2        TagV2
}

func New(configuration configuration.Configuration, language configuration.Language) *Registry {
	uri := configuration.RegistryUrl
	checkTLS := configuration.CheckTLS
	r := &Registry{
		url:           strings.TrimRight(uri, "/"),
		checkTLS:      checkTLS,
		request:       gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: !checkTLS}),
		configuration: configuration,
		language:      language,
	}

	resp, _, err := r.request.Get(r.url).End()
	if len(err) > 0 {
		panic(err)
		return nil
	}

	//Everything ok
	if resp.StatusCode == 200 {
		return r
	} else {
		log.Fatal("Can't create RegistryClient, Status : ", resp.StatusCode)
		return nil
	}
}

func (r *Registry) makeRequest(uri string, version uint) *http.Response {
	headerAccept := fmt.Sprintf("application/vnd.docker.distribution.manifest.v%d+json", version)
	resp, _, err := r.request.Get(r.url+uri).Set("Accept", headerAccept).End()
	if len(err) > 0 {
		panic(err)
	}

	if resp.StatusCode != 200 {
		log.Println("Error during request : ", r.url+uri, " | Status : ", resp.StatusCode)
		return nil
	}

	return resp
}

func (r *Registry) GetCatalog() Catalog {
	query := url.QueryEscape("_catalog")
	requestUrl := fmt.Sprintf("/%s", query)
	resp := r.makeRequest(requestUrl, 2)

	var record Catalog
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	return record
}

func (r *Registry) GetTags(imageName string) Repository {
	query := url.QueryEscape(imageName)
	requestUrl := fmt.Sprintf("/%s%s", query, "/tags/list")
	resp := r.makeRequest(requestUrl, 2)

	var record Repository
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	sort.Strings(record.Tags)
	return record
}

func (r *Registry) GetTagsInfo(imageName, tagName string) (Image, error) {

	imageName = url.QueryEscape(imageName)
	tagName = url.QueryEscape(tagName)

	requestUrl := fmt.Sprintf("/%s/manifests/%s", imageName, tagName)
	resp := r.makeRequest(requestUrl, 2)
	if resp == nil {
		return Image{}, errors.New(r.language.ImageOrTagNotFound)
	}

	var image Image
	image.Digest = resp.Header.Get("Docker-Content-Digest")[7:]
	image.Registry = strings.TrimRight(strings.TrimLeft(r.url, "https://"), "/v2/")

	if err := json.NewDecoder(resp.Body).Decode(&image.TagV2); err != nil {
		log.Println(err)
	}

	resp = r.makeRequest(requestUrl, 1)
	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		log.Println(err)
	}

	for _, element := range image.TagV2.Layers {
		image.Size += element.Size
	}

	image.Size += image.TagV2.Config.Size

	return image, nil
}
