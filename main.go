package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type URLs struct {
	Url string
}

func (u *URLs) GetJSON(url string) []byte {

	r, err := http.Get(url)
	if err != nil {
		log.Fatal((err))
	}

	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errRead != nil {
		log.Fatal(errRead)
	}
	return body
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	// u := URLs{}
	var u = &URLs{Url: "https://groupietrackers.herokuapp.com/api/artists"}
	// "https://groupietrackers.herokuapp.com/api/artists",
	// "https://groupietrackers.herokuapp.com/api/locations",
	// "https://groupietrackers.herokuapp.com/api/dates",
	// "https://groupietrackers.herokuapp.com/api/relation",
	// fmt.Println(artistsList)

	var artistsList Artists
	json.Unmarshal(u.GetJSON(u.Url), &artistsList)

	t, errParse := template.ParseFiles("assets/templates/index.html")
	if errParse != nil {
		fmt.Println("errParse")
		log.Fatal((errParse))
	}
	t.Execute(w, artistsList)
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	t, errParse := template.ParseFiles("assets/templates/artist.html")
	if errParse != nil {
		log.Fatal(errParse)
	}
	t.Execute(w, nil)
}

func main() {
	log.Println("starting localhost:8080...")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", rootHandle)
	// http.HandleFunc("/", artistHandle)
	http.ListenAndServe(":8080", nil)
}
