package main

import (
	"database/sql"
	"fmt"
	"github.com/alexquar/U-Watchlist/handlers"
	"github.com/alexquar/U-Watchlist/middleware"
	"github.com/alexquar/U-Watchlist/models"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

func main() {
	fmt.Println("Starting server on :8080")
	var err error
	models.DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer models.DB.Close()
	_, err = models.DB.Exec("CREATE TABLE IF NOT EXISTS films (ID INTEGER PRIMARY KEY AUTOINCREMENT, Title TEXT, Director TEXT, Year INTEGER, User TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /new", handlers.New)
	mux.HandleFunc("DELETE /delete/{ID}", handlers.Delete)
	mux.HandleFunc("GET /update/{ID}", handlers.UpdateTemplate)
	mux.HandleFunc("PUT /update/{ID}", handlers.Update)
	mux.HandleFunc("/", handlers.Home)
	wrapped := middleware.UUIDCookieMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", wrapped))

}
