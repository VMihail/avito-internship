package main

import (
	"avito-internship/internal/apiserver"
	"avito-internship/internal/store"
	"avito-internship/internal/utils"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	configPath = "/Users/mihildrozdov/GolandProjects/avito-internship/config/apiserver.json"
)

func main() {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := apiserver.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	err = store.Init()
	if err != nil {
		return
	}
	err = utils.Init()
	if err != nil {
		return
	}
	s := apiserver.New(&config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	err = store.Close()
	if err != nil {
		return
	}
}
