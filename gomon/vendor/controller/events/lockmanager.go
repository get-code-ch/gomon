package events

import (
	"controller/command"
	"controller/host"
	"controller/probe"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Lock struct {
	Timestamp time.Time
	Object    string
	Id        primitive.ObjectID
	ClientId  int64
}

var lockArray map[primitive.ObjectID]Lock

func init() {
	lockArray = make(map[primitive.ObjectID]Lock)
}

func AddLock(l Lock) {
	lockArray[l.Id] = l
}

func RemoveLock(id primitive.ObjectID) {
	delete(lockArray, id)
}

func IsLocked(id primitive.ObjectID) bool {
	_, exist := lockArray[id]
	return exist
}

func GetByClientId(cid int64) map[primitive.ObjectID]bool {
	result := make(map[primitive.ObjectID]bool)

	for key, value := range lockArray {
		if value.ClientId == cid {
			result[key] = true
		}
	}

	return result
}

func UnlockByClientId(cid int64) {
	var j []byte
	for key, value := range lockArray {
		if value.ClientId == cid {
			switch value.Object {
			case "PROBE":
				o := probe.Probe{Id: value.Id}
				j, _ = json.Marshal(o)
			case "COMMAND":
				o := command.Command{Id: value.Id}
				j, _ = json.Marshal(o)
			case "HOST":
				o := host.Host{Id: value.Id}
				j, _ = json.Marshal(o)
			}
			Broadcast <- SocketMessage{Object: value.Object, Action: "UNLOCK", Data: string(j), ErrorCode: 0}
			delete(lockArray, key)
		}
	}
}
