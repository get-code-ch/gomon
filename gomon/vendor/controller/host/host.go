package host

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"model"
)

type Host struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Key    string             `json:"key" bson:"Key"`
	Name   string             `json:"name" bson:"Name"`
	FQDN   string             `json:"fqdn" bson:"FQDN"`
	IP     string             `json:"ip" bson:"IP"`
	Locked bool               `json:"locked,omitempty"`
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

	if h.Id != primitive.NilObjectID {
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
	h.Id = primitive.NewObjectID()
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
