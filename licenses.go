package main

import (
	"bytes"
	"html/template"
	"net/http"
)

type Licenses struct {
	PageData PageData
}

func LicensesInit() Licenses {
	return Licenses{
		PageData: DataInit(" | Licenses"),
	}
}

func templatingLicenses(data *Licenses) bytes.Buffer {
	var buf bytes.Buffer

	t, err := template.ParseFiles("tmpl/main.tmpl", "tmpl/licenses.tmpl")
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

func handleLicenses(w http.ResponseWriter, _ *http.Request) {
	Page := LicensesInit()
	tmpl := templatingLicenses(&Page)
	w.Write(tmpl.Bytes())
}