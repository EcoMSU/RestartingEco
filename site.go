package main

import (
	"errors"
	"io"
	"net/http"
)

var ErrSiteWrongMode = errors.New("Site initialized in unexpected mode")

type page interface {
	Build() (io.Reader, error)
	Handle(http.ResponseWriter, *http.Request)
	Watch() func()
}

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

type SiteMode int8

const (
	BuildMode SiteMode = iota
	RunMode
)

func SiteData(TitleAdd string) Info {
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
		Share: SiteData(TitleAdd),
	}
}

type site struct {
	mode     SiteMode
	server   *http.ServeMux
	data     Info
	pages    map[string]page
	patterns []string
	static   [][3]string
}

func NewSite(mode SiteMode) (s *site) {
	s = new(site)
	s.mode = mode
	s.pages = make(map[string]page)
	return s
}

func (s *site) Run() error {
	if s.mode != RunMode {
		return ErrSiteWrongMode
	}
	if s.server == nil {
		logger.Println("Try to run with http.DefaultServeMux")
		s.server = http.DefaultServeMux
	}
	return http.ListenAndServe(":8080", s.server)
}

// patterns [3]string = pattern string, cut string, realpath string
func (s *site) SetStatic(patterns [][3]string) {
	s.static = patterns
}

func (s *site) AddPage(pattern string, p page) {
	s.patterns = append(s.patterns, pattern)
	s.pages[pattern] = p
}

func (s *site) ServeTo(sm *http.ServeMux) {
	for _, p := range s.static {
		sm.Handle(p[0], http.StripPrefix(p[1],
			http.FileServer(http.Dir(p[2])),
		))
	}
	for _, p := range s.patterns {
		sm.HandleFunc(p, s.pages[p].Handle)
	}
}
