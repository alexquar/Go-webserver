package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexquar/U-Watchlist/middleware"
	"github.com/alexquar/U-Watchlist/models"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	rows, err := models.DB.Query(("SELECT * FROM films WHERE User = ?"), middleware.GetUserID(r))
	if err != nil {
		http.Error(w, "Failed to retrieve films", http.StatusInternalServerError)
	}
	defer rows.Close()
	films := make(map[string][]models.Film)
	for rows.Next() {
		var film models.Film
		err := rows.Scan(&film.ID, &film.Title, &film.Director, &film.Year, &film.User)
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

func New(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		director := r.FormValue("director")
		year := r.FormValue("year")
		userID := middleware.GetUserID(r)

		if title == "" || director == "" || year == "" || userID == "" {
			http.Error(w, "Title, Director, Year and User must exist", http.StatusBadRequest)
			return
		}
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			http.Error(w, "Invalid year format", http.StatusBadRequest)
			return
		}
		entry, err := models.DB.Exec("INSERT INTO FILMS (Title, Director, Year, User) VALUES (?, ?, ?, ?)", title, director, yearInt, userID)
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
			ID:       id,
			Year:     &yearInt,
			User:     &userID,
		})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		userID := middleware.GetUserID(r)
		_, err := models.DB.Exec("DELETE FROM films WHERE ID = ? AND User = ?", id, userID)
		if err != nil {
			http.Error(w, "Failed to delete film", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.PathValue("ID"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var film models.Film
		userID := middleware.GetUserID(r)
		err = models.DB.QueryRow("SELECT ID, Title, Director, Year, User FROM films WHERE ID = ? AND User = ?", id, userID).Scan(&film.ID, &film.Title, &film.Director, &film.Year, &film.User)
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

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		id, _ := strconv.Atoi(r.PathValue("ID"))
		title := r.FormValue("title")
		director := r.FormValue("director")
		year := r.FormValue("year")
		userID := middleware.GetUserID(r)
		if title == "" || director == "" || year == "" || userID == "" {
			http.Error(w, "Title, Director, Year and User cannot be empty", http.StatusBadRequest)
		}
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			http.Error(w, "Invalid year format", http.StatusBadRequest)
			return
		}
		_, err = models.DB.Exec(("UPDATE films SET Title = ?, Director = ?, Year = ? WHERE ID = ? AND User = ?"), title, director, yearInt, id, userID)
		if err != nil {
			http.Error(w, "Failed to update film", http.StatusInternalServerError)
		}
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.ExecuteTemplate(w, "filmCard", models.Film{
			Title:    title,
			Director: director,
			Year:     &yearInt,
			ID:       int64(id),
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
