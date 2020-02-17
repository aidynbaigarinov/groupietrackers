package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Artists struct {
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

func rootHandle(w http.ResponseWriter, r *http.Request) {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	// urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	// urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	// urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	req, err := http.Get(urlArtists)
	if err != nil {
		log.Fatal((err))
	}

	body, readErr := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}

	artists1 := Artists{}
	json.Unmarshal(body, &artists1)

	t, err := template.ParseFiles("index.html")
	t.Execute(w, artists1)
}

func main() {
	log.Println("starting localhost:8080...")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
