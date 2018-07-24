# Harbor

Harbor is a lightweight Web UI designed in Golang for your private Docker Registry.

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/wallforfry/Harbor)](https://goreportcard.com/report/gitlab.com/wallforfry/Harbor)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Harbor is designed in Golang, if you want to build it on your computer you need Golang.
You can also run in in Docker container.

### Installing

##### With Docker :

```
$ git clone https://github.com/wallforfry/harbor
$ cd harbor
$ docker built -t golang/harbor .
$ docker run -it --name harbor -p 3008:3008 golang/harbor
```

Your project is now running in a container named _harbor_ and expose port 3008 to your host.
Access to http://localhost:3008/ to see the Web UI

##### On your computer :

```
$ go get wallforfry.fr/harbor
$ cd $GOPATH/src/wallforfry.fr/harbor
$ go get -d -v ./...
$ go install -v ./...
$ $GOPATH/bin/harbor
```
Your project is now running and expose port 3008 to your host.
Access to http://localhost:3008/ to see the Web UI

## Configuration

You can change some values in the config file _config.json_ like :
- __Port__ : The port of the webserver | _default 3008_
- __RegistryUrl__ : Your private registry URL | _default ""_
- __AppTitle__ : The name of your organisation | _default "WebUI"_
- __Language__ : The language of the UI | _default "en"_

##### Language
You can choose language value between __fr__ and __en__ but you can also create your own language file in the locales directory.

## Deployment

You can deploy Harbor inside Docker or has a service.

## Built With

* [Mux](https://github.com/gorilla/mux) - A powerful URL router and dispatcher for Golang
* [Gonfig](https://github.com/tkanos/gonfig) - Manage Configuration file and environment in GO
* [Materialize](https://materializecss.com/) - A modern responsive front-end framework based on Material Design
* [GoRequest](https://github.com/parnurzeal/gorequest/) - Simplified HTTP client
* [GJSON](https://github.com/tidwall/gjson/) - Get JSON values quickly - JSON Parser for Go

## Contributing

Please read [CONTRIBUTING.md](#) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/wallforfry/harbor/tags).

## Authors

* **Wallerand Delevacq** - *Initial work* - [wallforfry](https://github.com/wallforfry)

See also the list of [contributors](https://github.com/wallforfry/harbor/contributors) who participated in this project.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.