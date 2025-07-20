package main

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/mohits-git/xss-lab/internal/auth"
)

func (cfg *apiConfig) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	errMsg := q.Get("error")
	email := q.Get("email")

	data := struct {
		Email string
		Error error
	}{
		Email: email,
		Error: (func() error {
			if errMsg == "" {
				return nil
			}
			return errors.New(errMsg)
		})(),
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Error loading login template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, "Error executing login template", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Redirect(w, r, "/login?error=Email and password are required", http.StatusSeeOther)
		return
	}

	user, err := cfg.db.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("Error fetching user by email: %v\n", err)
		http.Redirect(w, r, "/login?error=Invalid email or password", http.StatusSeeOther)
		return
	}

	fmt.Printf("User found: %+v\n", user)

	if user == nil || auth.ComparePassword(user.PasswordHash, password) != nil {
		http.Redirect(w, r, "/login?error=Invalid email or password", http.StatusSeeOther)
		return
	}

	token, err := auth.MakeJwt(fmt.Sprintf("%d", user.ID), cfg.jwtSecret, time.Hour)
	if err != nil {
		http.Error(w, "Error creating JWT token", http.StatusInternalServerError)
		return
	}

	auth.SetAuthHeader(w.Header(), token)

	w.Write([]byte("Login successful! You can now access protected resources."))
}

func (cfg *apiConfig) registerPageHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	errMsg := q.Get("error")
	email := q.Get("email")
	name := q.Get("name")

	data := struct {
		Email string
		Name  string
		Error error
	}{
		Email: email,
		Name:  name,
		Error: (func() error {
			if errMsg == "" {
				return nil
			}
			return errors.New(errMsg)
		})(),
	}

	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Error loading register template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, "Error executing register template", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) registerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		http.Redirect(w, r, "/register?error=All fields are required", http.StatusSeeOther)
		return
	}

	if len(password) < 8 {
		http.Redirect(w, r, "/register?error=Password must be at least 8 characters long", http.StatusSeeOther)
		return
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(name, email, passwordHash)
	fmt.Printf("User created: %+v\n", user)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJwt(fmt.Sprintf("%d", user.ID), cfg.jwtSecret, time.Hour)
	if err != nil {
		http.Error(w, "Error creating JWT token", http.StatusInternalServerError)
		return
	}

	auth.SetAuthHeader(w.Header(), token)
	w.WriteHeader(http.StatusCreated)
}
