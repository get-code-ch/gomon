package layout

import (
	"controller/config"
	"encoding/json"
	"model"
	"model/menu"
)

// Custom type definition
type Layout struct {
	Authenticated bool
	Config        map[string]string
	Footer        map[string]string
	Menu          map[string]model.Href
}

// Constant and Global variables definition
var TemplateLayout []string

var LayoutData Layout

func init() {

	// Define layout template
	TemplateLayout = []string{"vendor/view/layout.html", "vendor/view/header.html", "vendor/view/menu.html", "vendor/view/footer.html"}

	// Menu definition - Converting map[string]interface{} -> Href
	LayoutData.Menu = make(map[string]model.Href)
	menu, _ := menu.ReadMenu()

	for _, item := range menu {
		m := model.Href{}
		bytes, _ := json.Marshal(item)
		json.Unmarshal(bytes, &m)
		LayoutData.Menu[m.Key] = m
	}

	LayoutData.Config = make(map[string]string)
	LayoutData.Config["Title"] = config.Config.Title
	LayoutData.Config["Powered"] = "Powered by gomon - claude@get-code.ch"

	LayoutData.Footer = make(map[string]string)
	LayoutData.Footer["Message"] = ""

}
