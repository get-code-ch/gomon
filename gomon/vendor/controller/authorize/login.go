package authorize

import (
	"controller"
	"controller/events"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, UserContext)

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
	file, err := ioutil.ReadFile(controller.Config.Users)
	if err != nil {
		events.EventsMsg <- "Authentication failed for: " + usr

		session.Values["authenticated"] = false
		session.Values["message"] = "Internal server error: " + err.Error()
		session.Save(r, w)
		if controller.Config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	hashPwd := search.FindSubmatch(file)
	if hashPwd == nil || hashPwd[2] == nil {
		// Go to home page
		events.EventsMsg <- "Authentication failed for: " + usr

		session.Values["authenticated"] = false
		session.Values["message"] = "Invalid Username/Password..."
		session.Save(r, w)
		if controller.Config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword(hashPwd[2], []byte(pwd)); err != nil {
		// Go to home page
		events.EventsMsg <- "Authentication failed for: " + usr

		session.Values["authenticated"] = false
		session.Values["message"] = "Invalid Username/Password..."
		session.Save(r, w)
		if controller.Config.Ssl {
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
	if controller.Config.Ssl {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
	}
}
