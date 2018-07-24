package configuration

// Configuration : prototype of config.json file
type Configuration struct {
	Port        int
	RegistryUrl string
	CheckTLS    bool
	AppTitle    string
	Language    string
}

// Language : prototype of <locale>.lang file
type Language struct {
	AppSubTitle        string
	About              string
	Settings           string
	Developed          string
	Golang             string
	Credits            string
	Informations       string
	Search             string
	ErrorString        string
	ImageOrTagNotFound string
	Back               string
}
