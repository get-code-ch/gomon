package model

import "github.com/mongodb/mongo-go-driver/bson/objectid"

type Href struct {
	Id      objectid.ObjectID `json:"id" bson:"_id"`
	Key     string            `json:"key" bson:"Key"`
	Link    string            `json:"link" bson:"Link"`
	Text    string            `json:"text" bson:"Text"`
	Visible bool              `json:"visible" bson:"Visible"`
}

type Menu []Href

type Msg struct {
	Msg string `json:"msg"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
	Msg   string `json:"msg"`
}

type Exception struct {
	Msg string `json:"msg"`
}
