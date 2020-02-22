package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

/*
type LocationsJSON struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type DatesJSON struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
*/
type RelationJSON struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var rel []RelationJSON

func rootHandle(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, http.StatusNotFound)
		// notFound(w)
	} else {
		url := "https://groupietrackers.herokuapp.com/api"
		a, err := getAPI(url)
		if err != nil {
			fmt.Println("1")
			errorHandler(w, http.StatusInternalServerError)
		} else {
			json.Unmarshal(a, &all)
			parseIndex(w)
		}
	}
}

func artistHandle(w http.ResponseWriter, r *http.Request) {
	artistBtn := r.FormValue("artist-btn")
	searchBar := r.FormValue("search-bar")
	found := false
	if artistBtn == "" && searchBar == "" {
		rand.Seed(time.Now().UnixNano())
		min := 2
		max := 53
		tmp := rand.Intn(max-min+1) + min
		artistBtn = artists[tmp-1].Name
	}
	if searchBar == "" && artistBtn != "" {
		for _, v := range artists {
			if v.Name == artistBtn {
				parseArtist(w, &v)
				found = true
			}
		}
		if !found {
			errorHandler(w, http.StatusBadRequest)
		}
	}
}

func searchHandle(w http.ResponseWriter, r *http.Request) {
	var searchResult []ArtistsJSON
	input := r.FormValue("search-bar")
	for _, v := range artists {

		if strings.ToLower(v.Name) == strings.ToLower(input) {
			searchResult = append(searchResult, v)
		} else {
			for _, i := range v.Members {
				if strings.ToLower(i) == strings.ToLower(input) {
					searchResult = append(searchResult, v)
				}
			}
		}

		if strings.ToLower(strconv.Itoa(v.CreationDate)) == strings.ToLower(input) {
			searchResult = append(searchResult, v)
		}

		if strings.ToLower(v.FirstAlbum) == strings.ToLower(input) {
			searchResult = append(searchResult, v)
		}

		for key, value := range v.RelationsData.DatesLocations {
			if key == input {
				searchResult = append(searchResult, v)
			}
			for _, i := range value {
				if i == input {
					searchResult = append(searchResult, v)
				}
			}
		}
	}
	if len(searchResult) == 0 {
		errorHandler(w, http.StatusBadRequest)
	} else {
		parseSearch(w, &searchResult)
	}

}

// func search(str, input string) string {
// 	fmt.Println(str)
// 	fmt.Println(input)
// 	lenInput := len(input)
// 	if len(input) > len(str) {
// 		return input
// 	}
// 	if len(input) == 0 {
// 		return input
// 	}
// 	for i, l := 0, len(str); i < l; i++ {
// 		if input[0] == str[i] && i > i+lenInput {
// 			fmt.Println(string(str[i]), i, l-lenInput)
// 			fmt.Println(str[i : i+len(input[1:])])
// 			if input[1:] == str[i:i+len(input[1:])] {
// 				return ""
// 			}
// 		}
// 	}
// 	return input
// }

func parseIndex(w http.ResponseWriter) {
	a, err := getAPI(all.Artists)
	if err != nil {
		fmt.Println("2")
		errorHandler(w, http.StatusInternalServerError)
	} else {
		json.Unmarshal(a, &artists)

		b, err := getAPI(all.Relation)
		if err != nil {
			fmt.Println("3")
			errorHandler(w, http.StatusInternalServerError)
		} else {

			b = b[9 : len(b)-2]
			json.Unmarshal(b, &rel)

			for i := range artists {
				artists[i].RelationsData = rel[i]
			}

			t, errParse := template.ParseFiles("assets/templates/index.html")
			if errParse != nil {
				errorHandler(w, http.StatusInternalServerError)
			} else {
				t.Execute(w, &artists)
			}
		}
	}
}

func parseArtist(w http.ResponseWriter, v *ArtistsJSON) {
	t, errParse := template.ParseFiles("assets/templates/artist.html")
	if errParse != nil {
		fmt.Println("6")
		errorHandler(w, http.StatusInternalServerError)
	} else {
		t.Execute(w, v)
	}
}

func parseSearch(w http.ResponseWriter, v *[]ArtistsJSON) {
	t, errParse := template.ParseFiles("assets/templates/search.html")
	if errParse != nil {
		fmt.Println(errParse)
		errorHandler(w, http.StatusInternalServerError)
	} else {
		t.Execute(w, v)
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
	http.HandleFunc("/search", searchHandle)
	http.ListenAndServe(":8080", nil)
}
