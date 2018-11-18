package view

import (
	"controller/authorize"
	"net/http"
	"text/template"
)

func Root(w http.ResponseWriter, r *http.Request) {
	var view []string

	session, _ := authorize.Store.Get(r, authorize.UserContext)
	logout := viewData.Menu["9000LOGOUT"]
	logout.Visible = false
	viewData.Menu["9000LOGOUT"] = logout

	if msg, ok := session.Values["message"].(string); ok {
		viewData.Footer["Message"] = msg
	} else {
		viewData.Footer["Message"] = ""
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		view = append(templateLayout, "view/login.html")
		t, err := template.ParseFiles(view...)
		if err != nil {
			http.Error(w, "root() - login template - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ViewData ViewData
		}{ViewData: viewData}
		t.ExecuteTemplate(w, "layout", data)
		return
	}

	logout = viewData.Menu["9000LOGOUT"]
	logout.Visible = true
	viewData.Menu["9000LOGOUT"] = logout

	view = append(templateLayout, "view/main.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		http.Error(w, "root() - main template - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		ViewData ViewData
	}{ViewData: viewData}
	t.ExecuteTemplate(w, "layout", data)
	return

}
