package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexquar/U-Watchlist/models"
	_ "modernc.org/sqlite"
)

var db *sql.DB

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
	mux := http.NewServeMux()
	mux.HandleFunc("POST /new", new)
	mux.HandleFunc("DELETE /delete/{ID}", delete)
	mux.HandleFunc("GET /update/{ID}", updateTemplate)
	mux.HandleFunc("PUT /update/{ID}", update)
	mux.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	rows, err := db.Query(("SELECT * FROM films"))
	if err != nil {
		http.Error(w, "Failed to retrieve films", http.StatusInternalServerError)
	}
	defer rows.Close()
	films := make(map[string][]models.Film)
	for rows.Next() {
		var film models.Film
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
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", models.Film{
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
		id, err := strconv.Atoi(r.PathValue("ID"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var film models.Film
		err = db.QueryRow("SELECT ID, Title, Director FROM films WHERE ID = ?", id).Scan(&film.ID, &film.Title, &film.Director)
		if err != nil {
			http.Error(w, "Film not found", http.StatusNotFound)
			return
		}

		tmpl := template.Must(template.ParseFiles("./templates/update.html"))
		err = tmpl.ExecuteTemplate(w, "updateCard", film)
		if err != nil {
			http.Error(w, "Template execution failed", http.StatusInternalServerError)
			fmt.Println("Template error:", err)
		}
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
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", models.Film{
			Title:    title,
			Director: director,
			ID:       int64(id),
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
