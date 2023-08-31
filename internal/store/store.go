package store

import (
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	configPath = "/Users/mihildrozdov/GolandProjects/avito-internship/config/dataBase.json"
)

var (
	driverName string
	url        string
	db         *sql.DB
	logger     *logrus.Logger
)

func Init() error {
	logger = logrus.New()
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	url = config.DbUrl
	driverName = config.DriverName
	d, err := sql.Open(driverName, url)
	if err != nil {
		return err
	}
	db = d
	logger.Info("Connected to database")
	return nil
}

func Close() error {
	err := db.Close()
	if err != nil {
		return err
	}
	logger.Info("Closed connection to database")
	return nil
}
