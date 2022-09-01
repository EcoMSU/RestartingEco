package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

// {
//     "name": "Николай Терещенко",
//     "img": "tereshenko.jpg",
//     "desc": [
//         "Руководитель департамента регионального развития",
//         "Советник генерального директора, Аргументы и Факты"
//     ]
// }

// {
//     "name": "Федеральное агентство по делам молодежи (Росмолодежь)",
//     "img": "rosmol.svg",
//     "url": "https://myrosmol.ru/"
// },

type Speaker struct {
	Name string   `json:"name"`
	Img  string   `json:"img"`
	Desc []string `json:"desc"`
}

type NewSpeaker struct {
	Name []string `json:"name"`
	Img  []string `json:"img"`
	Desc []string `json:"desc"`
}

type OldSpeaker struct {
	Id   string
	Name int
	Img  int
	Desc []int
}

type Partner struct {
	Name string `json:"name"`
	Img  string `json:"img"`
	Url  string `json:"url"`
}

type NewPartner struct {
	Name []string `json:"name"`
	Img  []string `json:"img"`
	Url  []string `json:"url"`
}

type OldPartner struct {
	Id   string
	Name int
	Img  int
	Url  int
}

func main() {
	partner()
	speaker()
}

func speaker() {
	var s []Speaker
	var ns map[string]NewSpeaker = make(map[string]NewSpeaker)
	var olds []OldSpeaker
	var err error

	if err = get("./speakers.json", &s); err != nil {
		fmt.Println(err)
	}

	for _, i := range s {
		ext := strings.Split(i.Img, ".")
		img := strings.Builder{}
		img.WriteString(uuid.New().String())
		img.Write([]byte("."))
		img.WriteString(ext[1])

		copy(path.Join("../static/img/speakers/", i.Img), path.Join("./img/speakers", img.String()))
		iId := uuid.New().String()

		ni := NewSpeaker{
			Name: []string{i.Name},
			Img:  []string{img.String()},
			Desc: i.Desc,
		}
		ns[iId] = ni

		var oldDesc []int

		for j := 0; j < len(i.Desc); j++ {
			oldDesc = append(oldDesc, j)
		}

		oi := OldSpeaker{
			Id:   iId,
			Name: 0,
			Img:  0,
			Desc: oldDesc,
		}
		olds = append(olds, oi)
	}

	saveJson("newspeakers.json", ns)
	saveJson("oldspeakers.json", olds)
}

func partner() {
	var p []Partner
	var np map[string]NewPartner = make(map[string]NewPartner)
	var op []OldPartner
	var err error

	if err = get("./partners.json", &p); err != nil {
		fmt.Println(err)
	}

	for _, i := range p {
		ext := strings.Split(i.Img, ".")
		img := strings.Builder{}
		img.WriteString(uuid.New().String())
		img.Write([]byte("."))
		img.WriteString(ext[1])

		copy(path.Join("../static/img/partners/", i.Img), path.Join("./img/partners", img.String()))
		iId := uuid.New().String()

		ni := NewPartner{
			Name: []string{i.Name},
			Img:  []string{img.String()},
			Url:  []string{i.Url},
		}
		np[iId] = ni

		oi := OldPartner{
			Id:   iId,
			Name: 0,
			Img:  0,
			Url:  0,
		}
		op = append(op, oi)
	}

	saveJson("newpartners.json", np)
	saveJson("oldpartners.json", op)
}

func get(filename string, data any) (err error) {
	var jsonData []byte
	if jsonData, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(jsonData, data)
	}
	return
}

func saveJson(filename string, data any) (err error) {
	var newData []byte
	var jsonFile *os.File
	if newData, err = json.MarshalIndent(data, "", "\t"); err == nil {
		if jsonFile, err = os.Create(filename); err == nil {
			_, err = jsonFile.Write(newData)
			if err == nil {
				err = jsonFile.Close()
			} else {
				jsonFile.Close()
			}
		}
	}
	return
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
