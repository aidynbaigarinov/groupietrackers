package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate string   `json:"creationDate"`
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

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}

	artists1 := artists{}

	json.Unmarshal(body, &artists1)

	fmt.Println(artists1)
}

func main() {
	// http.HandleFunc("/", rootHandle)
}
