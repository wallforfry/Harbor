package main

import (
    "html/template"
    "log"
    "net/http"
    "encoding/json"
    "fmt"
    "net/url"
    "github.com/gorilla/mux"
    "github.com/tkanos/gonfig"
)


var configuration	Configuration
var language		Language


type PageBase struct {
	Configuration	Configuration
	Lang			Language
}

type ImagesPage struct {
	Base   PageBase
	Images []Image
}

type Catalog struct {
	Repositories	[]string   `json:"repositories"`
}

type Image struct {
	Name	string		`json:"name"`
	Tags	[]string	`json:"tags"`
}

func makeGetHttpRequest(query string) *http.Response {
	log.Print(query)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Fatal("newHttpGetRequest: ", err)
		return nil
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("clientDoHttpGetRequest ", err)
		return nil
	}

	return resp
}

func getTags(baseUrl string, imageName string) Image{
	query := url.QueryEscape(imageName)
	requestUrl := fmt.Sprintf("%s%s%s", baseUrl, query, "/tags/list")
	resp := makeGetHttpRequest(requestUrl)
	defer resp.Body.Close()

	var record Image
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	return record
}

func getCatalog(baseUrl string) Catalog{
	query := url.QueryEscape("_catalog")
	requestUrl := fmt.Sprintf("%s%s", baseUrl, query)
	resp := makeGetHttpRequest(requestUrl)
	defer resp.Body.Close()

	var record Catalog
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	return record
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	baseUrl := configuration.RegistryUrl

	base := PageBase{configuration, language}
	p := ImagesPage{base, nil}

	catalog := getCatalog(baseUrl)

	for _, element := range catalog.Repositories {
		p.Images = append(p.Images, getTags(baseUrl, element))
	}

	t := template.New("Main")

	t = template.Must(t.ParseFiles("templates/layout.html", "templates/image_list.html"))

	err := t.ExecuteTemplate(w, "layout", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func main() {
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	err = gonfig.GetConf(fmt.Sprintf("locales/%s.lang", configuration.Language), &language)
	if err != nil {
		err = gonfig.GetConf("locales/en.lang", &language)
		if err != nil {
			panic(err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", viewHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Launching Server on port :", configuration.Port)
	http.ListenAndServe(fmt.Sprintf("%s%d", ":",configuration.Port), r)
}
