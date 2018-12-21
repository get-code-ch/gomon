package root

import (
	"controller/authorize"
	"controller/layout"
	"net/http"
	"text/template"
)

func init() {

}

func Root(w http.ResponseWriter, r *http.Request) {
	var view []string

	session, _ := authorize.Store.Get(r, authorize.UserContext)

	if msg, ok := session.Values["message"].(string); ok {
		layout.LayoutData.Footer["Message"] = msg
	} else {
		layout.LayoutData.Footer["Message"] = ""
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		layout.LayoutData.Authenticated = false
		view = append(layout.TemplateLayout, "vendor/view/login.html")
		t, err := template.ParseFiles(view...)
		if err != nil {
			http.Error(w, "root() - login template - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ViewData layout.Layout
		}{ViewData: layout.LayoutData}
		t.ExecuteTemplate(w, "layout", data)
		return
	}
	layout.LayoutData.Authenticated = true

	view = append(layout.TemplateLayout, "vendor/view/main.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		http.Error(w, "root() - main template - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		ViewData layout.Layout
	}{ViewData: layout.LayoutData}
	t.ExecuteTemplate(w, "layout", data)
	return

}
