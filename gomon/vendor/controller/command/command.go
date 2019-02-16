package command

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"model"
)

type Command struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Key         string             `json:"key" bson:"Key"`
	Name        string             `json:"name" bson:"Name"`
	Description string             `json:"description" bson:"Description"`
	Command     string             `json:"command" bson:"Command"`
	CommandType string             `json:"command_type" bson:"CommandType"`
	Locked      bool               `json:"locked,omitempty"`
}

func (c *Command) Get() ([]Command, error) {

	var (
		commands []Command
		filter   bson.M
	)

	if c.Id != primitive.NilObjectID {
		filter = bson.M{"_id": c.Id}
	} else {
		filter = nil
	}

	cursor, err := model.MongoDB.Collection("command").Find(model.Ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(model.Ctx)
	commands = []Command{}
	for cursor.Next(model.Ctx) {
		item := Command{}
		err := cursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		commands = append(commands, item)
	}
	return commands, nil
}

func (c *Command) Post() error {
	c.Id = primitive.NewObjectID()
	_, err := model.MongoDB.Collection("command").InsertOne(model.Ctx, c)
	return err

}
func (c *Command) Delete() error {
	filter := bson.M{"_id": c.Id}
	_, err := model.MongoDB.Collection("command").DeleteOne(model.Ctx, filter)
	return err
}
func (c *Command) Put() error {
	filter := bson.M{"_id": c.Id}
	_, err := model.MongoDB.Collection("command").UpdateOne(model.Ctx, filter, bson.D{
		{"$set", c},
	})
	return err
}
