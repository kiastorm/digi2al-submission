package main

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	// "time"
)

// TODO: Improve this
type Page struct {
	Title    string
	Errors   map[string][]string
	FormData map[string]string
	Username string
}

var users = map[string]string{
	"admin@example.com": "admin",
}

var loginTokens = map[string]string{}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendLoginEmail(email, token string) error {
	from := "your-email@example.com"
	password := "your-email-password"
	to := email
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Login Link\n\n" +
		"Click the link to login: http://localhost:8080/login?token=" + token

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
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

func loginActionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")

	page := &Page{Title: "Login", Errors: make(map[string][]string), FormData: map[string]string{"email": email}}

	// if email isn't found
	if _, ok := users[email]; !ok {
		page.Errors["email"] = append(page.Errors["email"], "Email not found")
	}

	if len(page.Errors) > 0 {
		t, err := template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/login.html", "views/partials/login-form.html")

		t.Execute(w, page)

		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	token, err := generateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	loginTokens[token] = email

	err = sendLoginEmail(email, token)

	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	page = &Page{Title: "Login", FormData: map[string]string{"email": email}}
	t, err := template.ParseFiles("views/layouts/base.html", "views/layouts/auth.html", "views/login-sent.html")

	t.Execute(w, page)

	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func loginTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, ok := loginTokens[token]

	if !ok {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: email,
	})

	delete(loginTokens, token)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	http.HandleFunc("/", homeRouteHandler)
	http.HandleFunc("/login", loginRouteHandler)
	http.HandleFunc("/login-action", loginActionHandler)
	http.HandleFunc("/login-token", loginTokenHandler)
	http.HandleFunc("/logout", logoutActionHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
