package probe

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"model"
	"time"
)

type Probe struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	HostId      primitive.ObjectID `json:"host_id" bson:"host_id"`
	CommandId   primitive.ObjectID `json:"command_id" bson:"command_id"`
	Name        string             `json:"name" bson:"Name"`
	Description string             `json:"description" bson:"Description"`
	Interval    int64              `json:"interval" bson:"Interval"`
	Next        time.Time          `json:"next" bson:"Next"`
	Last        time.Time          `json:"last" bson:"Last"`
	Result      string             `json:"result" bson:"Result"`
	State       string             `json:"state" bson:"State"`
	Locked      bool               `json:"locked,omitempty"`
}

type Controller interface {
	Get()
	Post()
	Delete()
	Put()
}

func (p *Probe) Get() ([]Probe, error) {

	var (
		probes []Probe
		filter bson.M
	)

	if p.Id != primitive.NilObjectID {
		filter = bson.M{"_id": p.Id}
	} else {
		filter = nil
	}

	cursor, err := model.MongoDB.Collection("probe").Find(model.Ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(model.Ctx)
	probes = []Probe{}
	for cursor.Next(model.Ctx) {
		item := Probe{}
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		probes = append(probes, item)
	}
	return probes, nil
}

func (p *Probe) Post() error {
	p.Id = primitive.NewObjectID()
	_, err := model.MongoDB.Collection("probe").InsertOne(model.Ctx, p)
	return err

}
func (p *Probe) Delete() error {
	filter := bson.M{"_id": p.Id}
	_, err := model.MongoDB.Collection("probe").DeleteOne(model.Ctx, filter)
	return err
}
func (p *Probe) Put() error {
	filter := bson.M{"_id": p.Id}
	_, err := model.MongoDB.Collection("probe").UpdateOne(model.Ctx, filter, bson.D{
		{"$set", p},
	})
	return err
}
