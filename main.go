package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

type Info struct {
	Title string
	Desc string
	URL string
	TitleAdd string
}

type PageData struct {
	Lang string
	Share Info
}

func SiteShare(TitleAdd string) Info {
	return Info {
		Title: "Экология: Перезагрузка",
		TitleAdd: TitleAdd,
		Desc: "Всероссийская онлайн-конференция Студенческого совета МГУ с международным участием о разрушении популярных экологических мифов.",
		URL: "https://RestartingEco.ru/",
	}
}

func DataInit(TitleAdd string) PageData {
	return PageData{
		Lang: "ru",
		Share: SiteShare(TitleAdd),
	}
}

func build() {
	var tmpl bytes.Buffer

	logger.Println("Start build")

	IndexPage := IndexInit()
	tmpl = templatingIndex(&IndexPage)

	index_html, _ := os.Create("static/index.html")
	defer index_html.Close()
	index_html.Write(tmpl.Bytes())

	LincensesPage := LicensesInit()
	tmpl = templatingLicenses(&LincensesPage)

	licenses_html, _ := os.Create("static/licenses.html")
	defer licenses_html.Close()
	licenses_html.Write(tmpl.Bytes())
}

func run() {
	logger.Println("Running...")

	fs_css := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs_css))
	fs_fonts := http.FileServer(http.Dir("./static/fonts"))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fs_fonts))
	fs_img := http.FileServer(http.Dir("./static/img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs_img))
	fs_icon := http.FileServer(http.Dir("./static/icon"))
	http.Handle("/icon/", http.StripPrefix("/icon/", fs_icon))

	http.HandleFunc("/licenses.html", handleLicenses)
	http.HandleFunc("/", handleIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
}

func main() {
	if len(os.Args) != 1 {
		switch os.Args[1] {
		case "build":
			build()
		default:
			run()
		}
	} else {
		run()
	}
}
