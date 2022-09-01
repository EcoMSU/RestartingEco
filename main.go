package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

var logger *log.Logger

func build() (err error) {
	var data io.Reader
	var buildPath string

	logger.Println("Start build")

	// Move static
	buildPath, err = filepath.Abs("./build")
	if err = copy.Copy("./static", buildPath); err != nil {
		log.Fatalln(err)
		log.Fatalln("Copy static failed")
		return
	}

	IndexPage := NewIndex(buildPath)
	data, err = IndexPage.Build()

	index_html, _ := os.Create("build/index.html")
	defer index_html.Close()
	_, err = io.Copy(index_html, data)

	LincensesPage := NewLicenses()
	data, err = LincensesPage.Build()

	licenses_html, _ := os.Create("build/licenses.html")
	defer licenses_html.Close()
	_, err = io.Copy(licenses_html, data)
	return
}

func watch() {
	fmt.Println("WIP")
}

func run() {
	var server *http.ServeMux
	site := NewSite()
	logger.Println("Running...")

	tempbuild, err := os.MkdirTemp("", "")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Temp directory for build", tempbuild)
	defer os.RemoveAll(tempbuild)

	server = http.NewServeMux()

	server.Handle("/css/", http.StripPrefix(
		"/css/",
		http.FileServer(http.Dir("./static/css")),
	))
	server.Handle("/fonts/", http.StripPrefix(
		"/fonts/",
		http.FileServer(http.Dir("./static/fonts")),
	))
	server.Handle("/img/partners/",
		http.FileServer(http.Dir(tempbuild)),
	)
	server.Handle("/img/speakers/",
		http.FileServer(http.Dir(tempbuild)),
	)
	server.Handle("/img/", http.StripPrefix(
		"/img/",
		http.FileServer(http.Dir("./static/img")),
	))
	server.Handle("/icon/", http.StripPrefix(
		"/icon/",
		http.FileServer(http.Dir("./static/icon")),
	))

	site.Add("/licenses.html", NewLicenses())
	site.Add("/", NewIndex(tempbuild))
	site.ServeTo(server)

	log.Fatal(http.ListenAndServe(":8080", server))
}

func empty() {
	logger.Println("No or wrong option found")
	fmt.Print(`
Expected options:
  watch
  	serve site and rebuild every time template file changed
  run
  	serve site with launch template file state
  build
  	build site into folder "build"
`)
}

func init() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
}

func main() {
	if len(os.Args) != 1 {
		logger.Println(os.Args[1], "option found")
		switch os.Args[1] {
		case "watch":
			watch()
		case "run":
			run()
		case "build":
			build()
		default:
			empty()
		}
	} else {
		empty()
	}
}
