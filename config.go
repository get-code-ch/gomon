package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Server       string
	Port         string
	Key          string
	Cert         string
	Ssl          bool
	StaticFolder string `json:"static_folder"`
	Users        string
	Title        string
	Menu         []Href
}

const fileConfig = "./config/config.json"

var config Config

func init() {
	var err error

	config, err = getConfiguration(fileConfig)

	if err != nil {
		log.Fatal("gomon init error: ", err)
	}
	return

}

func getConfiguration(file string) (Config, error) {
	var configuration Config
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(buffer, &configuration)
	return configuration, err
}
