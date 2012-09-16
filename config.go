package main

import (
	"encoding/json"
	"log"
	"os"
)

var config map[string]string

func GetConfig() map[string]string {
	if len(config) == 0 {
		file, err := os.Open("config.json")
		defer file.Close()
		if err != nil {
			config["host"] = "127.0.0.1"
			config["port"] = "8081"
			config["database"] = "iwebdb"
			log.Println("GetConfig err, use default setting.")
			log.Println(err)
			return config
		}

		setting := json.NewDecoder(file)
		setting.Decode(&config)
		log.Println("GetConfig OK!")
	}

	return config
}
