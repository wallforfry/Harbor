package main

type Configuration struct {
	Port        int
	RegistryUrl string
	CheckTLS    bool
	AppTitle    string
	Language    string
}

type Language struct {
	AppSubTitle  string
	About        string
	Settings     string
	Developed    string
	Golang       string
	Credits      string
	Informations string
	Search       string
}
