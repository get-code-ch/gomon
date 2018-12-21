package host

import (
	"controller/authorize"
	"controller/config"
	"controller/layout"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"html/template"
	"model"
	"net/http"
)

type Host struct {
	Id   objectid.ObjectID `json:"id" bson:"_id"`
	Key  string            `json:"key" bson:"Key"`
	Name string            `json:"name" bson:"Name"`
	FQDN string            `json:"fqdn" bson:"FQDN"`
	IP   string            `json:"ip" bson:"IP"`
}

type Controller interface {
	Get()
	Post()
	Delete()
	Put()
}

func (h *Host) Get() ([]Host, error) {

	var (
		hosts  []Host
		filter bson.M
	)

	if h.Id != objectid.NilObjectID {
		filter = bson.M{"_id": h.Id}
	} else {
		filter = nil
	}

	cursor, err := model.MongoDB.Collection("host").Find(model.Ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(model.Ctx)
	hosts = []Host{}
	for cursor.Next(model.Ctx) {
		item := Host{}
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, item)
	}
	return hosts, nil
}

func (h *Host) Post() error {
	h.Id = objectid.New()
	_, err := model.MongoDB.Collection("host").InsertOne(model.Ctx, h)
	return err

}
func (h *Host) Delete() error {
	filter := bson.M{"_id": h.Id}
	_, err := model.MongoDB.Collection("host").DeleteOne(model.Ctx, filter)
	return err
}
func (h *Host) Put() error {
	filter := bson.M{"_id": h.Id}
	_, err := model.MongoDB.Collection("host").UpdateOne(model.Ctx, filter, bson.D{
		{"$set", h},
	})
	return err
}

func CreateHosts(w http.ResponseWriter, r *http.Request) {
	var h Host

	session, _ := authorize.Store.Get(r, authorize.UserContext)

	err := r.ParseForm()
	if err != nil {
		session.Values["message"] = "CreateHosts() - Error parsing form error: " + err.Error()
	}
	h.Id = objectid.New()
	h.Key = r.Form["Key"][0]
	h.Name = r.Form["Name"][0]
	h.FQDN = r.Form["FQDN"][0]
	h.IP = r.Form["IP"][0]

	err = h.Post()
	if err != nil {
		session.Values["message"] = "CreateHosts() - CreateHost error: " + err.Error()
		session.Save(r, w)
	} else {
		session.Values["message"] = "New hosts successfully created"
	}

	// Go to home page
	if config.Config.Ssl {
		http.Redirect(w, r, "https://"+r.Host+"/admin#hosts", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "http://"+r.Host+"/admin#hosts", http.StatusSeeOther)
	}
}

func ListHosts(w http.ResponseWriter, r *http.Request) {
	var view []string
	var host Host

	session, _ := authorize.Store.Get(r, authorize.UserContext)

	view = append(layout.TemplateLayout, "vendor/view/hosts.html", "vendor/view/hostslist.html", "vendor/view/hostsform.html")
	t, err := template.ParseFiles(view...)
	if err != nil {
		session.Values["message"] = "ListHosts() - loading template - Internal Server Error: " + err.Error()
		session.Save(r, w)
		return
	}

	h, _ := host.Get()
	data := struct {
		ViewData layout.Layout
		Hosts    []Host
	}{ViewData: layout.LayoutData, Hosts: h}
	t.ExecuteTemplate(w, "layout", data)
	return

}
