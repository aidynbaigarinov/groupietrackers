package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type APIurls struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

var all *APIurls

type ArtistsJSON []struct {
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

type LocationsJSON struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type DatesJSON struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type RelationJSON struct {
	Index []struct {
		ID            int                 `json:"id"`
		DatesLocation map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	b := getAPI(all.Artists)
	var artists *ArtistsJSON
	json.Unmarshal(b, &artists)

	t, errParse := template.ParseFiles("assets/templates/index.html")
	if errParse != nil {
		fmt.Println("errParse")
		log.Fatal((errParse))
	}
	t.Execute(w, &artists)
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	t, errParse := template.ParseFiles("assets/templates/artist.html")
	if errParse != nil {
		log.Fatal(errParse)
	}
	t.Execute(w, nil)
}

func getAPI(url string) []byte {
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errRead != nil {
		log.Fatal(errRead)
	}
	return body
}

func main() {
	log.Println("starting localhost:8080...")

	http.HandleFunc("/artist", artistHandle)
	a := getAPI("https://groupietrackers.herokuapp.com/api")
	json.Unmarshal(a, &all)

	c := getAPI(all.Relation)
	var relation *RelationJSON
	json.Unmarshal(c, &relation)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", rootHandle)
	http.ListenAndServe(":8080", nil)
}
