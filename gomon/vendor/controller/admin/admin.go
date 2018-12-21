package admin

import (
	"controller/authorize"
	"controller/config"
	"controller/host"
	"controller/layout"
	"html/template"
	"model/probes"
	"net/http"
)

func init() {
}

func Admin(w http.ResponseWriter, r *http.Request) {
	// TODO: Test authorization level of user

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

	view = append(layout.TemplateLayout, "vendor/view/admin.html", "vendor/view/hostslist.html", "vendor/view/hostsform.html", "vendor/view/probeslist.html", "vendor/view/probesform.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		session.Values["message"] = "Admin() - loading template - Internal Server Error: " + err.Error()
		session.Save(r, w)
		return
	}

	p, _ := probes.ReadProbes()
	h, _ := new(host.Host).Get()
	data := struct {
		ViewData layout.Layout
		Probes   []probes.Probe
		Hosts    []host.Host
	}{ViewData: layout.LayoutData, Probes: p, Hosts: h}
	t.ExecuteTemplate(w, "layout", data)
	return

}
