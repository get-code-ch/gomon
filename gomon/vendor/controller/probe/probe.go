package probe

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"model"
	"time"
)

type Probe struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	HostId      primitive.ObjectID `json:"host_id, omitempty" bson:"host_id, omitempty"`
	CommandId   primitive.ObjectID `json:"command_id" bson:"command_id"`
	Name        string             `json:"name" bson:"Name"`
	Description string             `json:"description" bson:"Description"`
	Interval    int64              `json:"interval" bson:"Interval"`
	Next        time.Time          `json:"next" bson:"Next"`
	Last        time.Time          `json:"last" bson:"Last"`
	Result      string             `json:"result" bson:"Result"`
	State       string             `json:"state" bson:"State"`
	Locked      bool               `json:"locked,omitempty"`
	Username    string             `json:"username,omitempty" bson:"Username"`
	Password    string             `json:"password,omitempty" bson:"Password,omitempty"`
	Secret      string             `json:"secret,omitempty" bson:"Secret"`
}

var secret = []byte("1FE3B7IIB05CGBK2F6D17KJ61H36OLJJ")

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
		item.Secret = ""
		item.Password = ""
		probes = append(probes, item)
	}
	return probes, nil
}

func (p *Probe) GetSecret() (string, error) {

	var filter bson.M

	if p.Id != primitive.NilObjectID {
		filter = bson.M{"_id": p.Id}
	} else {
		return "", nil
	}

	cursor, err := model.MongoDB.Collection("probe").Find(model.Ctx, filter, nil)
	if err != nil {
		return "", err
	}
	defer cursor.Close(model.Ctx)
	for cursor.Next(model.Ctx) {
		item := Probe{}
		err := cursor.Decode(&item)
		if err == nil && item.Secret != "" {
			return decrypt([]byte(item.Secret)), err
		}
	}
	return "", err
}

func (p *Probe) Post() error {
	p.Id = primitive.NewObjectID()
	if p.Password != "" {
		p.Secret = encrypt([]byte(p.Password))
		p.Password = ""
	} else {
		p.Secret = ""
	}
	_, err := model.MongoDB.Collection("probe").InsertOne(model.Ctx, p)
	p.Secret = ""
	return err
}

func (p *Probe) Delete() error {
	filter := bson.M{"_id": p.Id}
	_, err := model.MongoDB.Collection("probe").DeleteOne(model.Ctx, filter)
	return err
}
func (p *Probe) Put() error {
	filter := bson.M{"_id": p.Id}
	if p.Password != "" {
		if p.Password == " " {
			p.Secret = ""
		} else {
			p.Secret = encrypt([]byte(p.Password))
		}
		p.Password = ""
	} else {
		pwd, _ := p.GetSecret()
		p.Secret = encrypt([]byte(pwd))
	}
	_, err := model.MongoDB.Collection("probe").UpdateOne(model.Ctx, filter, bson.D{
		{"$set", p},
	})
	p.Secret = ""
	return err
}

func encrypt(data []byte) string {
	block, _ := aes.NewCipher([]byte(secret))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return string(ciphertext)
}

func decrypt(data []byte) string {
	key := []byte(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}
