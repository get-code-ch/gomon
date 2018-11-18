package view

import (
	"controller"
	"database/sql"
	"encoding/json"
	"types"
)

// Custom type definition
type ViewData struct {
	Config map[string]string
	Footer map[string]string
	Menu   map[string]types.Href
}

// Constant and Global variables definition
var templateLayout []string

var viewData ViewData

func init() {
	//var err error

	// Define layout template
	templateLayout = []string{"view/layout.html", "view/header.html", "view/menu.html", "view/footer.html"}

	// Menu definition - Converting map[string]interface{} -> Href
	viewData.Menu = make(map[string]types.Href)
	menu, _ := getMenu()

	for _, item := range menu {
		m := types.Href{}
		bytes, _ := json.Marshal(item)
		json.Unmarshal(bytes, &m)
		viewData.Menu[m.Key] = m
	}

	viewData.Config = make(map[string]string)
	viewData.Config["Title"] = controller.Config.Title
	viewData.Config["Powered"] = "Powered by gomon - claude@get-code.ch"

	viewData.Footer = make(map[string]string)
	viewData.Footer["Message"] = ""

}

func getMenu() ([]map[string]interface{}, error) {
	var (
		menu      []map[string]interface{}
		data      string
		jsonArray string
	)

	db, err := sql.Open("postgres", controller.Config.Database)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select data from menu")
	defer rows.Close()

	jsonArray = "["
	row := rows.Next()
	for row {
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		row = rows.Next()
		if row {
			jsonArray += data + ","
		} else {
			jsonArray += data + "]"
		}
	}
	json.Unmarshal([]byte(jsonArray), &menu)

	if err != nil {
		menu = nil
	}
	return menu, err
}
