package main

import (
	"net/http"
	"text/template"
	// "html/template"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
