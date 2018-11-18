package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"types"
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
	Database     string
	Menu         []types.Href
}

const fileConfig = "./config/config.json"

var Config TConfig

func init() {
	var err error

	Config, err = getConfiguration(fileConfig)

	if err != nil {
		log.Fatal("gomon init error: ", err)
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
