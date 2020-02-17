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

func rootHandle(w http.ResponseWriter, r *http.Request) {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	// urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	// urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	// urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	req, err := http.Get(urlArtists)
	if err != nil {
		log.Fatal((err))
	}

	body, errRead := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if errRead != nil {
		log.Fatal(errRead)
	}

	var artistsList Artists
	json.Unmarshal(body, &artistsList)

	// fmt.Println(artistsList)

	t, errParse := template.ParseFiles("templates/index.html")
	if errParse != nil {
		fmt.Println("errParse")
		log.Fatal((err))
	}
	t.Execute(w, artistsList)
}

func main() {
	log.Println("starting localhost:8080...")
	http.HandleFunc("/", rootHandle)
	http.ListenAndServe(":8080", nil)
}
