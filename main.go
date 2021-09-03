package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
) 

// dummy user data 
var users = map[string]string{"user1": "password", "user2": "password"}
// creating a cookie session store
var store = sessions.NewCookieStore([]byte("secret_key"))

func main() {
  r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	// modifying http import struct to add an extra property 
	// of timeout (good practice)
	httpServer := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}
	// ParseForm parses the raw query from the URL and updates r.Form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return 
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// check if user exists 
	storedPassword, exists := users[username]
	if exists {
		// Get registers and returns a session for the given name and session store.
		// session.id is the name of the cookie that will be stored in the client's browser
		session, _ := store.Get(r, "session.id")
		if storedPassword == password {
			session.Values["authenticated"] = true 
			// saves all sessions used during the current request
			session.Save(r, w)
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		}
		w.Write([]byte("Login successfully!"))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get registers and returns a session for the given name and session store
	// session.id is the name of the cookie that will be stored in the client's browser
	session, _ := store.Get(r, "session.id")
	// Set the authenticated value on the session to false 
	session.Values["authenticated"] = false 
	session.Save(r, w)
	w.Write([]byte("Logout Successful"))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	if authenticated != nil && authenticated != false {
		w.Write([]byte("Welcome!"))
		return 
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return 
	}
}