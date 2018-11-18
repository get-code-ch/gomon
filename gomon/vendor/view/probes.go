package view

import (
	"controller"
	"controller/authorize"
	"html/template"
	"model/probes"
	"net/http"
)

func Probes(w http.ResponseWriter, r *http.Request) {
	var view []string

	session, _ := authorize.Store.Get(r, authorize.UserContext)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		session.Values["message"] = "Unauthorized access"
		session.Save(r, w)

		if controller.Config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	view = append(templateLayout, "view/probes.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		session.Values["message"] = "probes() - loading template - Internal Server Error: " + err.Error()
		session.Save(r, w)
		return
	}

	p, _ := probes.GetProbes()
	data := struct {
		ViewData ViewData
		Probes   []map[string]interface{}
	}{ViewData: viewData, Probes: p}
	t.ExecuteTemplate(w, "layout", data)
	return

}
