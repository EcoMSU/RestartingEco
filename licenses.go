package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

type (
	LicensesData struct {
		PageData
	}
	Licenses struct {
		built bool
		raw   []byte
		data  LicensesData
	}
)

func NewLicenses() (page *Licenses) {
	page = new(Licenses)
	page.data.PageData = DataInit(" | Licenses")
	return page
}

func (l *Licenses) reBuild() (err error) {
	w := new(bytes.Buffer)
	err = template.Must(
		template.ParseFiles("tmpl/main.tmpl", "tmpl/licenses.tmpl"),
	).Execute(w, &l.data)

	if err == nil {
		l.raw = w.Bytes()
		l.built = true
	}
	return err
}

func (l *Licenses) Build() (io.Reader, error) {
	var err error
	if !l.built {
		err = l.reBuild()
	}
	return bytes.NewReader(l.raw), err
}

func (l *Licenses) Handle(w http.ResponseWriter, _ *http.Request) {
	var err error
	if !l.built {
		err = l.reBuild()
	}
	if err == nil {
		logger.Println("Licenses page loaded")
		w.Write(l.raw)
	}
	if err != nil {
		logger.Panicln(err)
	}
}

func (l *Licenses) Watch() func() {
	logger.Fatalln("Not implemented yet")
	return nil
}
