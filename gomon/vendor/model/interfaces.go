package model

type Controller interface {
	HandleMessage()
	Get()
	Post()
	Delete()
	Put()
}
