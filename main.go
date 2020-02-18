package main

import (
	"encoding/json"
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

type ArtistsJSON struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	// Locations     string   `json:"locations"`
	// ConcertDates  string   `json:"concertDates"`
	Relations     string `json:"relations"`
	RelationsData RelationJSON
}

var artists []ArtistsJSON

type LocationsJSON struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type DatesJSON struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type RelationJSON struct {
	// ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		t, errParse := template.ParseFiles("assets/templates/error/404.html")
		if errParse != nil {
			log.Fatal(errParse)
		}
		t.Execute(w, nil)
	} else {
		b := getAPI(all.Artists)
		json.Unmarshal(b, &artists)

		t, errParse := template.ParseFiles("assets/templates/index.html")
		if errParse != nil {
			log.Fatal((errParse))
		}
		t.Execute(w, &artists)
	}
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("artist")
	found := false

	for _, v := range artists {
		if v.Name == name {
			a := getAPI(v.Relations)
			var rel RelationJSON
			json.Unmarshal(a, &rel)
			v.RelationsData = rel
			t, errParse := template.ParseFiles("assets/templates/artist.html")
			if errParse != nil {
				log.Fatal(errParse)
			}
			t.Execute(w, v)
			found = true
		}
	}
	if !found {
		t, errParse := template.ParseFiles("assets/templates/error/404.html")
		if errParse != nil {
			log.Fatal(errParse)
		}
		t.Execute(w, nil)
	}
	// TODO: handle 404
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

	url := "https://groupietrackers.herokuapp.com/api"
	a := getAPI(url)
	json.Unmarshal(a, &all)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/artist", artistHandle)
	http.ListenAndServe(":8080", nil)
}
