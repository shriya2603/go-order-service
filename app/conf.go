package app

import (
	"encoding/json"
	"log"
	"os"
)

const ConfigFilePath string = "./config.json"

var appConfig *Configuration

// Configuration: consist of configs from config.json file
type Configuration struct {
	DBHost     string `json:"host"`
	DBPort     int    `json:"port"`
	DBUser     string `json:"user"`
	DBPassword string `json:"password"`
	DBName     string `json:"dbname"`
}

func loadConfig() {
	configFile, err := os.Open(ConfigFilePath)
	printError(err, "Error while reading config file: ")
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
}

func printError(err error, errMsg string) {
	if err != nil {
		log.Fatalln(errMsg, err)
	}
}

func GetConfiguration() *Configuration {
	if appConfig == nil {
		loadConfig()
	}
	return appConfig
}
