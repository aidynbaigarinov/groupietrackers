package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
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
		errorHandler(w, http.StatusNotFound)
		// notFound(w)
	} else {
		url := "https://groupietrackers.herokuapp.com/api"
		a, err := getAPI(url)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError)
		} else {
			json.Unmarshal(a, &all)
			parseIndex(w)
		}
	}
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("artist")
	found := false
	fmt.Println(name)
	if name == "" {
		rand.Seed(time.Now().UnixNano())
		min := 2
		max := 53
		tmp := rand.Intn(max-min+1) + min
		name = artists[tmp-1].Name
	}
	for _, v := range artists {
		if v.Name == name {
			parseArtist(w, &v)
			found = true
		}
	}
	if !found {
		errorHandler(w, http.StatusNotFound)
	}

}

func parseIndex(w http.ResponseWriter) {
	b, err := getAPI(all.Artists)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
	} else {
		json.Unmarshal(b, &artists)

		t, errParse := template.ParseFiles("assets/templates/index.html")
		if errParse != nil {
			errorHandler(w, http.StatusInternalServerError)
		} else {
			t.Execute(w, &artists)
		}
	}
}

func parseArtist(w http.ResponseWriter, v *ArtistsJSON) {
	a, err := getAPI(v.Relations)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
	} else {
		var rel RelationJSON
		json.Unmarshal(a, &rel)
		v.RelationsData = rel
		t, errParse := template.ParseFiles("assets/templates/artist.html")
		if errParse != nil {
			errorHandler(w, http.StatusInternalServerError)
		} else {
			t.Execute(w, v)
		}
	}
}

func badRequest(w http.ResponseWriter) {
	t, errParse := template.ParseFiles("assets/templates/error/400.html")
	if errParse != nil {
		errorHandler(w, http.StatusInternalServerError)
	} else {
		t.Execute(w, nil)
	}
}

func notFound(w http.ResponseWriter) {
	t, errParse := template.ParseFiles("assets/templates/error/404.html")
	if errParse != nil {
		errorHandler(w, http.StatusInternalServerError)
	} else {
		t.Execute(w, nil)
	}
}

func internalServerError(w http.ResponseWriter) {
	t, errParse := template.ParseFiles("assets/templates/error/500.html")
	if errParse != nil {
		log.Fatal(errParse)
	} else {
		t.Execute(w, nil)
	}
}

func checkUrl(url string) bool {
	for i, n := 0, len(url); i < n; i++ {
		if i != len(url)-1 && url[i] == '%' {
			if url[i+1] == '%' {
				return false
			} else if url[i] == '{' {
				return false
			}
		}
	}
	return true
}

func errorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	if status == http.StatusBadRequest {
		badRequest(w)
	} else if status == http.StatusNotFound {
		notFound(w)
	} else if status == http.StatusInternalServerError {
		internalServerError(w)
	}
}

func getAPI(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errRead != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	log.Println("starting localhost:8080...")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/artist", artistHandle)
	http.ListenAndServe(":8080", nil)
}
