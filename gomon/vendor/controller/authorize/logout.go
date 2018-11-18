package authorize

import (
	"controller"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, UserContext)
	if auth, ok := session.Values["authenticated"].(bool); ok || auth {
		session.Values["authenticated"] = false
		session.Values["username"] = ""
		session.Values["message"] = "Logged out..."
		session.Save(r, w)
	}

	// Go to home page
	if controller.Config.Ssl {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
	}
}
