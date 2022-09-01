package main

import (
	"io"
	"net/http"
)

type Info struct {
	Title    string
	Desc     string
	URL      string
	TitleAdd string
}

type PageData struct {
	Lang  string
	Share Info
}

func SiteShare(TitleAdd string) Info {
	return Info{
		Title:    "Экология: Перезагрузка",
		TitleAdd: TitleAdd,
		Desc:     "Всероссийская онлайн-конференция Студенческого совета МГУ с международным участием о разрушении популярных экологических мифов.",
		URL:      "https://RestartingEco.ru/",
	}
}

func DataInit(TitleAdd string) PageData {
	return PageData{
		Lang:  "ru",
		Share: SiteShare(TitleAdd),
	}
}

type site struct {
	pages    map[string]page
	patterns []string
}

func NewSite() *site {
	var s site
	s.init()
	return &s
}

func (s *site) init() {
	s.pages = make(map[string]page)
}

func (s *site) Add(pattern string, p page) {
	s.patterns = append(s.patterns, pattern)
	s.pages[pattern] = p
}

func (s *site) ServeTo(sm *http.ServeMux) {
	for _, p := range s.patterns {
		sm.HandleFunc(p, s.pages[p].Handle)
	}
}

type page interface {
	Build() (io.Reader, error)
	Handle(http.ResponseWriter, *http.Request)
	Watch() func()
}
