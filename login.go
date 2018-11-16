package main

import (
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const userContext = "user-context"

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("1FE3B7IIB05CGBK2F6D17KJ61H36OLJJ")
	store = sessions.NewCookieStore(key)
)

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userContext)

	// Parse login form
	err := r.ParseForm()
	if err != nil {
		session.Values["authenticated"] = false
		session.Save(r, w)
		http.Error(w, "login() - Parse - Internal Server Error", http.StatusInternalServerError)
	}
	usr := r.Form["username"][0]
	pwd := r.Form["password"][0]
	search := regexp.MustCompile(`(?mi)^(` + usr + `):(.*)$`)

	// Check hashPwd authentication
	file, err := ioutil.ReadFile(config.Users)
	if err != nil {
		session.Values["authenticated"] = false
		session.Values["message"] = "Internal server error: " + err.Error()
		session.Save(r, w)
		if config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	hashPwd := search.FindSubmatch(file)
	if hashPwd == nil || hashPwd[2] == nil {
		// Go to home page
		session.Values["authenticated"] = false
		session.Values["message"] = "Invalid Username/Password..."
		session.Save(r, w)
		if config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword(hashPwd[2], []byte(pwd)); err != nil {
		// Go to home page
		session.Values["authenticated"] = false
		session.Values["message"] = "Invalid Username/Password..."
		session.Save(r, w)
		if config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	// Set hashPwd as authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = usr
	session.Values["message"] = ""
	session.Save(r, w)

	// Go to home page
	if config.Ssl {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, userContext)
	if auth, ok := session.Values["authenticated"].(bool); ok || auth {
		session.Values["authenticated"] = false
		session.Values["username"] = ""
		session.Values["message"] = "Logged out..."
		session.Save(r, w)
	}
	// Go to home page
	if config.Ssl {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
	}

}
