package main

import (
	"html/template"
	"net/http"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../../assets/templates/home.page.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)

}
