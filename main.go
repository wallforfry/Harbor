package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"gitlab.com/wallforfry/harbor/registry"
	"html/template"
	"log"
	"net/http"
)

type client struct {
	reg           *registry.Registry
	configuration Configuration
	language      Language
}

type PageBase struct {
	Configuration Configuration
	Lang          Language
}

type ErrorPage struct {
	Base    PageBase
	Code    int
	Message string
	Header  bool
}

type ImagesPage struct {
	Base   PageBase
	Images []registry.Repository
	Header bool
}

type ImageDetailsPage struct {
	Base   PageBase
	Image  registry.Image
	Header bool
}

func PrettifySize(size int) string {
	units := []string{"B", "KB", "MB", "GB"}
	i := 0
	for size > 1024 && i < len(units) {
		size = size / 1024
		i = i + 1
	}
	return fmt.Sprintf("%.*d %s", 0, size, units[i])
}

func (c *client) viewRepositories(w http.ResponseWriter, r *http.Request) {

	base := PageBase{c.configuration, c.language}
	p := ImagesPage{base, nil, true}

	catalog := c.reg.GetCatalog()

	for _, element := range catalog.Repositories {
		p.Images = append(p.Images, c.reg.GetTags(element))
	}

	t := template.New("Main")

	t = template.Must(t.ParseFiles("templates/layout.html", "templates/header.html", "templates/image_list.html"))

	err := t.ExecuteTemplate(w, "layout", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func (c *client) viewImage(w http.ResponseWriter, r *http.Request) {

	base := PageBase{c.configuration, c.language}
	tag, err := c.reg.GetTagsInfo(r.FormValue("image"), r.FormValue("tag"))
	if err != nil {
		//http.Redirect(w, r, "/error/", http.StatusTemporaryRedirect)
		c.viewError(w, r, http.StatusNotFound, err.Error())
		return
	}

	p := ImageDetailsPage{base, tag, false}

	t := template.New("Main")

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"prettifySize": func(size int) string {
			return PrettifySize(size)
		},
	}

	t.Funcs(funcMap)
	t = template.Must(t.ParseFiles("templates/layout.html", "templates/header.html", "templates/image_details.html"))

	err = t.ExecuteTemplate(w, "layout", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func (c *client) viewError(w http.ResponseWriter, r *http.Request, code int, message string) {
	base := PageBase{c.configuration, c.language}
	p := ErrorPage{base, code, message, false}

	t := template.New("Error")
	t = template.Must(t.ParseFiles("templates/layout.html", "templates/header.html", "templates/error.html"))

	err := t.ExecuteTemplate(w, "layout", p)
	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func main() {
	var (
		c client
	)

	err := gonfig.GetConf("config.json", &c.configuration)
	if err != nil {
		panic(err)
	}

	err = gonfig.GetConf(fmt.Sprintf("locales/%s.lang", c.configuration.Language), c.language)
	if err != nil {
		err = gonfig.GetConf("locales/en.lang", &c.language)
		if err != nil {
			panic(err)
		}
	}

	c.reg = registry.New(c.configuration.RegistryUrl, c.configuration.CheckTLS)

	c.reg.GetTagsInfo("golang/harbor", "latest")

	r := mux.NewRouter()
	r.HandleFunc("/", c.viewRepositories)
	r.HandleFunc("/details", c.viewImage).Queries("image", "{image}").Queries("tag", "{tag}")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Launching Server on port :", c.configuration.Port)
	http.ListenAndServe(fmt.Sprintf("%s%d", ":", c.configuration.Port), r)
}
