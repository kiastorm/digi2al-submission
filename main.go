package main

import (
	"html/template"
	"log"
	"net/http"
)

// TODO: Improve this
type Page struct {
	Title    string
	Errors   map[string][]string
	FormData map[string]string
	Username string
}


var users = map[string]string{
	"admin": "password",
}

func homeRouteHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	username := ""

	if err == nil {
		username = cookie.Value
	}

	page := &Page{Title: "Home", Username: username}

	t, err := template.ParseFiles("views/layouts/base.html", "views/layouts/app.html", "views/home.html")

	t.Execute(w, page)

	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func loginRouteHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{Title: "Login", Errors: make(map[string][]string)}

	err := template.Must(template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/login.html", "views/partials/login-form.html")).Execute(w, page)

	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func signupRouteHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{Title: "Signup", Errors: make(map[string][]string)}

	err := template.Must(template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/signup.html", "views/partials/signup-form.html")).Execute(w, page)

	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func signupActionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirmPassword")

	page := &Page{Title: "Signup", Errors: make(map[string][]string), FormData: map[string]string{"username": username}}

	// if username is empty
	if username == "" {
		page.Errors["username"] = append(page.Errors["username"], "Username is required")
	}

	// if username already exists
	if _, ok := users[username]; ok {
		page.Errors["username"] = append(page.Errors["username"], "Username already exists")
	}

	// if password is empty
	if password == "" {
		page.Errors["password"] = append(page.Errors["password"], "Password is required")
	}

	// if password is less than 8 characters
	if len(password) < 8 {
		page.Errors["password"] = append(page.Errors["password"], "Password must be at least 8 characters")
	}

	// if password and password confirm do not match
	if password != passwordConfirm {
		page.Errors["confirmPassword"] = append(page.Errors["confirmPassword"], "Passwords do not match")
	}

	if len(page.Errors) > 0 {
		if r.Header.Get("Hx-Request") == "true" {

			t, err := template.ParseFiles("views/partials/signup-form.html")

			t.Execute(w, page)

			if err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		} else {
			t, err := template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/signup.html", "views/partials/signup-form.html")

			t.Execute(w, page)

			if err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
	}

	users[username] = password

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
	})

	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusOK)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func loginActionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	page := &Page{Title: "Login", Errors: make(map[string][]string), FormData: map[string]string{"username": username}}

	// if username isnt found
	if _, ok := users[username]; !ok {
		page.Errors["username"] = append(page.Errors["username"], "Username not found")
	}

	// if username IS found but password is incorrect
	if _, usernameExists := users[username]; usernameExists && users[username] != password {
		page.Errors["password"] = append(page.Errors["password"], "Password is incorrect")
	}

	if len(page.Errors) > 0 {
		if r.Header.Get("Hx-Request") == "true" {

			t, err := template.ParseFiles("views/partials/login-form.html")

			t.Execute(w, page)

			if err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		} else {
			t, err := template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/login.html", "views/partials/login-form.html")

			t.Execute(w, page)

			if err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
	})

	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")

		http.Redirect(w, r, "/", http.StatusOK)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func logoutActionHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		MaxAge: -1,
	})

	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("GET /", homeRouteHandler)
	http.HandleFunc("GET /login", loginRouteHandler)
	http.HandleFunc("GET /signup", signupRouteHandler)

	http.HandleFunc("POST /signup", signupActionHandler)
	http.HandleFunc("POST /login", loginActionHandler)
	http.HandleFunc("POST /logout", logoutActionHandler)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
