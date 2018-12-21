package menu

import (
	"fmt"
	"log"
	"model"
)

func ReadMenu() (model.Menu, error) {
	var menu model.Menu

	cursor, err := model.MongoDB.Collection("menu").Find(model.Ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("ReadMenu: couldn't not read menu: %v", err)
	}
	defer cursor.Close(model.Ctx)
	for cursor.Next(model.Ctx) {
		item := model.Href{}
		err := cursor.Decode(&item)
		if err != nil {
			log.Printf("Error reading menu: %v", err)
		}
		menu = append(menu, item)
	}
	return menu, nil
}
