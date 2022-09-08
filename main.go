package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/EcoMSU/sef"
	"github.com/otiai10/copy"
)

var logger *log.Logger

func prepare() (s sef.Site, err error) {
	var buildPath string
	buildPath, err = filepath.Abs("./build")
	os.MkdirAll(buildPath, os.ModePerm)

	s = sef.NewSite(sef.SiteData{
		Title: "Экология: Перезагрузка",
		Desc:  "Всероссийская онлайн-конференция Студенческого совета МГУ с международным участием о разрушении популярных экологических мифов.",
		URL:   "https://RestartingEco.ru/",
	})
	logger.Println("Running...")

	s.SetStatic([][3]string{
		{"/css/", "/css/", "./static/css"},
		{"/fonts/", "/fonts/", "./static/fonts"},
		{"/img/partners/", "/img/partners/", "./build/img/partners/"},
		{"/img/speakers/", "/img/speakers/", "./build/img/speakers/"},
		{"/img/", "/img/", "./static/img"},
		{"/icon/", "/icon/", "./static/icon"},
	})

	s.AddPage("licenses.html", NewLicenses(s.GetData()))
	s.AddPage("index.html", NewIndex(buildPath, s.GetData()))
	s.AddAlias("", "index.html")
	return
}

func build() (err error) {
	var buildPath string

	logger.Println("Start build")

	s, err := prepare()

	// Move static
	buildPath, err = filepath.Abs("./build")
	if err = copy.Copy("./static", buildPath); err != nil {
		log.Fatalln(err)
		log.Fatalln("Copy static failed")
		return
	}
	err = s.Build(buildPath)
	return
}

func watch() {
	fmt.Println("WIP")
}

func run() {
	var server *http.ServeMux
	var err error

	s, err := prepare()
	if err != nil {
		log.Fatal(err)
		return
	}

	server = http.NewServeMux()
	s.ServeTo(server)

	log.Fatal(s.Run(0))
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
	if len(os.Args) > 1 {
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
