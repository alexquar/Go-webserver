package main

import (
	"database/sql"
	"fmt"
	"github.com/alexquar/U-Watchlist/models"
	"html/template"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
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
