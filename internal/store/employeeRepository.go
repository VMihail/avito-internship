package store

import (
	"avito-internship/internal/model"
	"fmt"
	"log"
)

func GetAllEmployee() ([]model.Employee, error) {
	logger.Info("method <GetAllEmployee> started")
	rows, err := db.Query("select * from segmentation_service.employee")
	if err != nil {
		return nil, err
	}
	var result []model.Employee
	for rows.Next() {
		p := model.Employee{}
		err := rows.Scan(&p.Id)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, p)
	}
	return result, nil
}

func CreateEmployee() (*model.Employee, error) {
	logger.Info("method <CreateEmployee> started")
	result := model.Employee{}
	if err := db.QueryRow(
		"insert into segmentation_service.employee default values returning employee_id").
		Scan(&result.Id); err != nil {
		return nil, err
	}
	if err := result.Validate(); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetById(id int) (*model.Employee, error) {
	logger.Info("method <GetById> started")
	result := model.Employee{}
	if err := db.QueryRow(
		"select * from segmentation_service.employee where employee_id = $1", id).Scan(&result.Id); err != nil {
		return nil, err
	}
	if err := result.Validate(); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetExperimentsById(id int) ([]model.Experiment, error) {
	logger.Info("method <GetExperimentsById> started")
	rows, err := db.Query("select * from segmentation_service.experiment where experiment_id in (select experiment_id from segmentation_service.employee_experiment where employee_id = $1);", id)
	if err != nil {
		return nil, err
	}
	var result []model.Experiment
	for rows.Next() {
		p := model.Experiment{}
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func GetPercentageOfEmployee(percent float32) ([]model.Employee, error) {
	logger.Info("method <GetPercentageOfEmployee> started")
	rows, err := db.Query(fmt.Sprintf("select * from segmentation_service.employee limit (%f / 100.0) * (select count(*) from segmentation_service.employee);", percent))
	if err != nil {
		return nil, err
	}
	var result []model.Employee
	for rows.Next() {
		p := model.Employee{}
		err := rows.Scan(&p.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, err
}
