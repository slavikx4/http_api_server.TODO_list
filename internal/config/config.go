package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var Config *Configuration

type Configuration struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPassword string
	PgBase     string
}

func init() {
	Config = &Configuration{}

	file, err := os.Open("./internal/config/config.conf")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(data, Config); err != nil {
		log.Fatalln(err)
	}
}
