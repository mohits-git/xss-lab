package main

import (
	"net/http"
	"text/template"
	// "html/template" // for escaping HTML
)

func (cfg *apiConfig) landingPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title":   "XSS Lab",
		"Content": "Welcome to the XSS Lab! This is a safe environment to learn about Cross-Site Scripting (XSS) vulnerabilities and how to prevent them.",
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
