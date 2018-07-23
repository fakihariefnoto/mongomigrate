package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

type (

	// Pkg is interface to store package object after calling New func
	Pkg interface {
	}
	pkgConfig struct{}

	// Config struct for unite all config data
	Config struct {
		FileName string
		URL      string
		Path     string
		Config   struct {
			Updater  `json:"updater"`
			Receiver `json:"receiver"`
			Date     `json:"date"`
		}
	}

	Updater struct {
		Connection UpdaterConn `json:"connection"`
		Query      string      `json:"query"`
	}

	UpdaterConn struct {
		Name       string `json:"name"`
		ConnString string `json:"connection_string"`
	}

	Receiver struct {
		Connection ReceiverConn `json:"connection"`
	}

	ReceiverConn struct {
		Name       string `json:"name"`
		ConnString string `json:"connection_string"`
		Database   string `json:"database"`
		Collection string `json:"collection"`
	}

	Date struct {
		Start       string `json:"start"`
		End         string `json:"end"`
		LastEndDate string `json:"last_end_date"`
	}
)

var (
	cfg Config
)

/*
	{
		"updater" : {
			"connection" : [
				{
					"name" : "name connection",
					"connection_string" : "connection string"
				},
				{
					"name" : "name connection",
					"connection_string" : "connection string"
				}
			],
			"query" : "query postgre"
		},
		"receiver" : {
			"connection" : [
				{
					"name" : "name connection",
					"connection_string" : "connection string",
					"database" : "database name",
					"collection" : "collection name"
				},
				{
					"name" : "name connection",
					"connection_string" : "connection string",
					"database" : "database name",
					"collection" : "collection name"

				}
			]
		},
		"date" : {
			"start" : "2018-01-01",
			"end" : "2018-07-01",
			"last_end_date" : "2017-12-31"
		}
	}
*/

func Init() {

	cfg = Config{
		FileName: "config",
		Path:     ".",
	}

	if ok := readJSONConfig(&cfg); ok != nil {
		log.Fatalln("failed to read config -> " + ok.Error())
	}
}

// Get Config struct
func Get() Config {
	return cfg
}

// Save to save config to file json
func Save(config Config) (err error) {

	configbyte, err := json.MarshalIndent(config, " ", "    ")
	if err != nil {
		return errors.New("Error when marshaling config struct: " + err.Error())
	}

	filename := cfg.Path + "/" + cfg.FileName + ".json"

	err = ioutil.WriteFile(filename, configbyte, 0644)
	if err != nil {
		return errors.New("Error when writing to file: " + err.Error())
	}

	return
}

func readJSONConfig(config *Config) (err error) {

	fname := config.Path + "/" + config.FileName + ".json"

	bytefile, err := ioutil.ReadFile(fname)
	if err != nil {
		return errors.New("Error when open config err: " + err.Error())
	}

	err = json.Unmarshal(bytefile, &config.Config)
	if err != nil {
		return errors.New("Error when unmarshal config err: " + err.Error())
	}

	return
}
