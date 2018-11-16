package main

// Custom type definition
type Href struct {
	Key     string
	Link    string
	Text    string
	Visible bool
}

type ViewData struct {
	Config map[string]string
	Footer map[string]string
	Menu   map[string]Href
}

// Constant and Global variables definition
var templateLayout []string

var viewData ViewData

func init() {
	//var err error

	// Define layout template
	templateLayout = []string{"view/layout.html", "view/header.html", "view/menu.html", "view/footer.html"}

	// Menu definition
	viewData.Menu = make(map[string]Href)
	for _, item := range config.Menu {
		viewData.Menu[item.Key] = item
	}

	viewData.Config = make(map[string]string)
	viewData.Config["Title"] = config.Title
	viewData.Config["Powered"] = "Powered by gomon - claude@get-code.ch"

	viewData.Footer = make(map[string]string)
	viewData.Footer["Message"] = ""

}
