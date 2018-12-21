package probes

import (
	"controller/authorize"
	"controller/config"
	"controller/layout"
	"html/template"
	"model/probes"
	"net/http"
)

func ListProbes(w http.ResponseWriter, r *http.Request) {
	var view []string

	session, _ := authorize.Store.Get(r, authorize.UserContext)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		session.Values["message"] = "Unauthorized access"
		session.Save(r, w)

		if config.Config.Ssl {
			http.Redirect(w, r, "https://"+r.Host+"/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host+"/", http.StatusSeeOther)
		}
		return
	}

	view = append(layout.TemplateLayout, "vendor/view/probes.html", "vendor/view/probeslist.html", "vendor/view/probesform.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		session.Values["message"] = "probes() - loading template - Internal Server Error: " + err.Error()
		session.Save(r, w)
		return
	}

	p, _ := probes.ReadProbes()
	data := struct {
		ViewData layout.Layout
		Probes   []probes.Probe
	}{ViewData: layout.LayoutData, Probes: p}
	t.ExecuteTemplate(w, "layout", data)
	return

}
