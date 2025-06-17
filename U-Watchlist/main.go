package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type Film struct {
	Title    string
	Director string
	ID       int
}

var films = map[string][]Film{
	"Films": {
		{Title: "Inception", Director: "Christopher Nolan", ID: 1},
		{Title: "The Matrix", Director: "Lana Wachowski, Lilly Wachowski", ID: 2},
		{Title: "Interstellar", Director: "Christopher Nolan", ID: 3},
	},
}

func main() {
	fmt.Println("Starting server on :8080")

	http.HandleFunc("/", home)
	http.HandleFunc("POST /new", new)
	http.HandleFunc("DELETE /delete/{ID}", delete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, films)
}

func new(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		director := r.FormValue("director")
		films["Films"] = append(films["Films"], Film{Title: title, Director: director, ID: len(films["Films"]) + 1})
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", films["Films"][len(films["Films"])-1])
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodDelete {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		for i, film := range films["Films"] {
			if film.ID == id {
				films["Films"] = slices.Delete(films["Films"], i, i+1)
				tmpl := template.Must(template.ParseFiles("index.html"))
				tmpl.ExecuteTemplate(w, "films", films)
				return
			}
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
