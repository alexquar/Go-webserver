package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Film struct {
	Title    string
	Director string
	ID       int64
}

func main() {
	fmt.Println("Starting server on :8080")
	var err error
	db, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS films (ID INTEGER PRIMARY KEY AUTOINCREMENT, Title TEXT, Director TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", home)
	http.HandleFunc("POST /new", new)
	http.HandleFunc("DELETE /delete/{ID}", delete)
	http.HandleFunc("GET /update/{ID}", updateTemplate)
	http.HandleFunc("PUT /update/{ID}", update)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	rows, err := db.Query(("SELECT * FROM films"))
	if err != nil {
		http.Error(w, "Failed to retrieve films", http.StatusInternalServerError)
	}
	defer rows.Close()
	films := make(map[string][]Film)
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.ID, &film.Title, &film.Director)
		if err != nil {
			http.Error(w, "Failed to scan film", http.StatusInternalServerError)
		}
		films["Films"] = append(films["Films"], film)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading films", http.StatusInternalServerError)
	}
	tmpl.Execute(w, films)
}

func new(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		director := r.FormValue("director")
		if title == "" || director == "" {
			http.Error(w, "Title and Director cannot be empty", http.StatusBadRequest)
			return
		}
		entry, err := db.Exec("INSERT INTO FILMS (Title, Director) VALUES (?, ?)", title, director)
		if err != nil {
			http.Error(w, "Failed to add film", http.StatusInternalServerError)
			return
		}
		id, err := entry.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		}
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", Film{
			Title:    title,
			Director: director,
			ID:       id})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodDelete {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		_, err := db.Exec("DELETE FROM films WHERE ID = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete film", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateTemplate(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodGet {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		row, _ := db.Query("SELECT * FROM films WHERE ID = ?", id)
		var film Film
		err := row.Scan(&film.ID, &film.Title, &film.Director)
		if err != nil {
			http.Error(w, "Film not found", http.StatusNotFound)
		}
		tmpl := template.Must(template.ParseFiles("update.html"))
		tmpl.ExecuteTemplate(w, "updateCard", film)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	if r.Method == http.MethodPut {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		title := r.FormValue("title")
		director := r.FormValue("director")
		if title == "" || director == "" {
			http.Error(w, "Title and Director cannot be empty", http.StatusBadRequest)
		}
		_, err := db.Exec(("UPDATE films SET Title = ?, Director = ? WHERE ID = ?"), title, director, id)
		if err != nil {
			http.Error(w, "Failed to update film", http.StatusInternalServerError)
		}
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", Film{
			Title:    title,
			Director: director,
			ID:       int64(id),
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
