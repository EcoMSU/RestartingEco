package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Header struct {
	Title string
	Desc string
}

type Speaker struct {
	Name string
	Img string
	Desc []string
}
type SpeakersData struct {
	SpeakersPrefix string
	SpeakersList []Speaker
}

type Partner struct {
	Name string
	Img string
	Url string
	Style template.CSS
}
type PartnersData struct {
	PartnersPrefix string
	PartnersList []Partner
}

type Index struct {
	PageData PageData
	Header Header
	SpeakersData SpeakersData
	PartnersData PartnersData
}

func IndexInit() Index {
	return Index{
		PageData: DataInit(""),
		Header: Header{
			Title: "Экология:",
			Desc: "Перезагрузка",
		},
		SpeakersData: SpeakersData{
			SpeakersPrefix: "img/speakers/",
		},
		PartnersData: PartnersData{
			PartnersPrefix: "img/partners/",
		},
	}
}

func templatingIndex(data *Index) bytes.Buffer {
	var buf bytes.Buffer
	var err error
	
	speakersJson, _ := ioutil.ReadFile("json/speakers.json")
	err = json.Unmarshal([]byte(speakersJson), &data.SpeakersData.SpeakersList)
	if err != nil {
		logger.Println(err)
	}

	partnersJson, _ := ioutil.ReadFile("json/partners.json")
	err = json.Unmarshal([]byte(partnersJson), &data.PartnersData.PartnersList)
	if err != nil {
		logger.Println(err)
	}

	t, err := template.ParseFiles("tmpl/main.tmpl", "tmpl/index.tmpl", "tmpl/speakers.tmpl", "tmpl/partners.tmpl")
	if err != nil {
		logger.Println(err)
	}
	ruTmpl := template.Must(t, err)
	err = ruTmpl.Execute(&buf, &data)
	if err != nil {
		logger.Println(err)
	}
	return buf
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	Page := IndexInit()
	tmpl := templatingIndex(&Page)
	w.Write(tmpl.Bytes())
}