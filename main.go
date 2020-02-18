package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	b := getAPI(all.Artists)
	json.Unmarshal(b, &artists)

	t, errParse := template.ParseFiles("assets/templates/index.html")
	if errParse != nil {
		log.Fatal((errParse))
	}
	t.Execute(w, &artists)
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("artist"))
	if err != nil {
		log.Fatal(err)
	}

	for index, v := range artists {
		if index == id-1 {
			a := getAPI(v.Relations)
			var rel RelationJSON
			json.Unmarshal(a, &rel)
			v.RelationsData = rel
			t, errParse := template.ParseFiles("assets/templates/artist.html")
			if errParse != nil {
				log.Fatal(errParse)
			}
			t.Execute(w, v)
		} else {
			// TODO: Handle 404
		}
	}

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

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", rootHandle)
	http.ListenAndServe(":8080", nil)
}
