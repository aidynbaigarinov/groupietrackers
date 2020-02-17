package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type ArtistsIndex struct {
}

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

}

func artistsHandle(w http.ResponseWriter, r *http.Request) {
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

	var artists []Artists
	json.Unmarshal(body, &artists)

	// fmt.Println(artists[2])

	t, err := template.ParseFiles("static/index.html")
	t.Execute(w, artists)
}

func main() {
	log.Println("starting localhost:8080...")
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/artists", artistsHandle)
	http.ListenAndServe(":8080", nil)
}
