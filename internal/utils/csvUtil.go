package utils

import (
	"avito-internship/internal/model"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

const (
	configPath = "/Users/mihildrozdov/GolandProjects/avito-internship/config/csvFiles.json"
)

var (
	path string
)

func Init() error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	path = config.CsvFilesPath
	return nil
}

func GetCsvFile(reports []model.Report) (string, error) {
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		err := file.Close()
		if err != nil {
			return
		}
	}()
	header := []string{"EmployeeId", "ExperimentName", "Operation", "Date"}
	err = writer.Write(header)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(reports); i++ {
		row := []string{
			strconv.Itoa(reports[i].EmployeeId),
			reports[i].ExperimentName,
			reports[i].Operation,
			reports[i].Date,
		}
		err := writer.Write(row)
		if err != nil {
			return "", err
		}
	}
	return file.Name(), nil
}
