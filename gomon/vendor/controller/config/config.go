package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type TConfig struct {
	Server       string
	Port         string
	Key          string
	Cert         string
	Ssl          bool
	StaticFolder string `json:"static_folder"`
	Users        string
	Title        string
	Timer        int
	Db           string
	RedisSrv     string `json:"Redis_server"`
	RedisPwd     string `json:"Redis_password"`
	MongoSrv     string `json:"Mongo_server"`
	MongoUsr     string `json:"Mongo_username"`
	MongoPwd     string `json:"Mongo_password"`
	MongoDb      string `json:"Mongo_database"`
}

const envFileConfig = "GOMON_CONFIG"

var Config TConfig

func init() {
	// default config file
	var fileConfig = "./config/config.json"
	var err error

	fc := os.Getenv(envFileConfig)
	if fc != "" {
		fileConfig = fc
	}

	Config, err = getConfiguration(fileConfig)

	if err != nil {
		log.Fatal("gomon init error: ", err)
	} else {
		log.Printf("gomon init ok, config: %v", Config)
	}

	return

}

func getConfiguration(file string) (TConfig, error) {
	var configuration TConfig
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return TConfig{}, err
	}
	err = json.Unmarshal(buffer, &configuration)
	return configuration, err
}
