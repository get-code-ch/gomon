package history

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"model"
	"time"
)

type History struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	HostId    primitive.ObjectID `json:"host_id" bson:"host_id"`
	ProbeId   primitive.ObjectID `json:"probe_id" bson:"probe_id"`
	Comment   string             `json:"comment" bson:"Comment"`
	Timestamp time.Time          `json:"timestamp" bson:"Timestamp"`
	State     string             `json:"state" bson:"State"`
	Message   string             `json:"message" bson:"Message"`
	Metric    string             `json:"metric,omitempty" bson:"Metric,omitempty" `
}

type Controller interface {
	Get()
	Post()
	Delete()
	Put()
}

func (h *History) Get() ([]History, error) {

	var (
		probes []History
		filter bson.M
	)

	if h.Id != primitive.NilObjectID {
		filter = bson.M{"_id": h.Id}
	} else {
		filter = nil
	}

	cursor, err := model.MongoDB.Collection("history").Find(model.Ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(model.Ctx)
	probes = []History{}
	for cursor.Next(model.Ctx) {
		item := History{}
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		probes = append(probes, item)
	}
	return probes, nil
}

func (h *History) Post() error {
	h.Id = primitive.NewObjectID()
	_, err := model.MongoDB.Collection("history").InsertOne(model.Ctx, h)
	return err

}
func (h *History) Delete() error {
	filter := bson.M{"_id": h.Id}
	_, err := model.MongoDB.Collection("history").DeleteOne(model.Ctx, filter)
	return err
}
func (h *History) Put() error {
	filter := bson.M{"_id": h.Id}
	_, err := model.MongoDB.Collection("history").UpdateOne(model.Ctx, filter, bson.D{
		{"$set", h},
	})
	return err
}
