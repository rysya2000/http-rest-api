package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/rysya2000/http-rest-api/internal/app/apiserver"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/apiserver.json", "Path to the config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	// json data
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
