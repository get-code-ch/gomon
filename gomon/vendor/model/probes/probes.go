package probes

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"log"
	"model"
)

type Probe struct {
	Id          objectid.ObjectID `bson:"_id"`
	Key         string            `bson:"Key"`
	Name        string            `bson:"Name"`
	Description string            `bson:"Description"`
	Command     string            `bson:"Command"`
}

func ReadProbes() ([]Probe, error) {
	var probes []Probe

	cursor, err := model.MongoDB.Collection("probes").Find(model.Ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("ReadProbes: couldn't not read menu: %v", err)
	}
	defer cursor.Close(model.Ctx)
	for cursor.Next(model.Ctx) {
		item := Probe{}
		err := cursor.Decode(&item)
		if err != nil {
			log.Printf("Error reading menu: %v", err)
		}
		probes = append(probes, item)
	}
	return probes, nil
}
