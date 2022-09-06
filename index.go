package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"sync"

	"github.com/EcoMSU/emdb"
)

type (
	Header struct {
		Title string
		Desc  string
	}

	SpeakersData struct {
		ImgPrefix string
		List      []emdb.ExportSpeaker
	}

	PartnersData struct {
		ImgPrefix string
		List      []emdb.ExportPartner
	}

	ScheduleLine struct {
		Time    string
		Speaker string
	}

	Schedule struct {
		Saturday []ScheduleLine
		Sunday   []ScheduleLine
	}

	IndexData struct {
		PageData PageData
		Header   Header
		Speakers SpeakersData
		Partners PartnersData
		Schedule Schedule
	}

	Index struct {
		built bool
		raw   []byte
		wg    *sync.WaitGroup
		data  IndexData
	}
)

func NewIndex(build string) *Index {
	page := new(Index)
	page.init()
	page.wg.Add(1)
	go page.load(build)
	return page
}

func (i *Index) init() {
	i.data = IndexData{
		PageData: DataInit(""),
		Header: Header{
			Title: "Экология:",
			Desc:  "Перезагрузка",
		},
		Speakers: SpeakersData{
			ImgPrefix: "img/speakers/",
		},
		Partners: PartnersData{
			ImgPrefix: "img/partners/",
		},
	}
	i.wg = new(sync.WaitGroup)
}

func (i *Index) load(build string) {
	var rawSpeakers []emdb.ImportSpeaker
	var rawPartners []emdb.ImportPartner
	var err error

	saturdayJson, _ := ioutil.ReadFile("json/saturday.json")
	err = json.Unmarshal(saturdayJson, &i.data.Schedule.Saturday)
	if err != nil {
		logger.Println(err)
	}

	sundayJson, _ := ioutil.ReadFile("json/sunday.json")
	err = json.Unmarshal(sundayJson, &i.data.Schedule.Sunday)
	if err != nil {
		logger.Println(err)
	}

	speakersJson, _ := ioutil.ReadFile("json/speakers.json")
	if err = json.Unmarshal(speakersJson, &rawSpeakers); err != nil {
		logger.Println(err)
	}
	if i.data.Speakers.List, err = emdb.GetSpeakers(rawSpeakers, path.Join(build, i.data.Speakers.ImgPrefix), ""); err != nil {
		logger.Println(err)
	}

	partnersJson, _ := ioutil.ReadFile("json/partners.json")
	err = json.Unmarshal(partnersJson, &rawPartners)
	if err != nil {
		logger.Fatalln(err)
	}
	if i.data.Partners.List, err = emdb.GetPartners(rawPartners, path.Join(build, i.data.Partners.ImgPrefix), ""); err != nil {
		logger.Fatalln(err)
	}

	i.wg.Done()
}

func (i *Index) reBuild() (err error) {
	w := new(bytes.Buffer)
	i.wg.Wait()
	err = template.Must(
		template.ParseFiles("tmpl/main.tmpl", "tmpl/index.tmpl", "tmpl/speakers.tmpl", "tmpl/partners.tmpl"),
	).Execute(w, &i.data)
	if err == nil {
		i.raw = w.Bytes()
		i.built = true
	}
	return err
}

func (i *Index) Build() (io.Reader, error) {
	var err error
	if !i.built {
		err = i.reBuild()
	}
	return bytes.NewReader(i.raw), err
}

func (i *Index) Handle(w http.ResponseWriter, _ *http.Request) {
	var err error
	if !i.built {
		err = i.reBuild()
	}
	if err == nil {
		logger.Println("Index page loaded")
		_, err = w.Write(i.raw)
	}
	if err != nil {
		logger.Panicln(err)
	}
}

func (i *Index) Watch() func() {
	logger.Fatalln("Not implemented yet")
	return nil
}
