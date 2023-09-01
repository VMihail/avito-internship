package store

import (
	"avito-internship/internal/model"
	"database/sql"
	_ "database/sql"
	"fmt"
)

func CreateExperiment(name string) (*model.Experiment, error) {
	logger.Info("method <CreateExperiment> started")
	result := model.Experiment{}
	if err := db.QueryRow(
		"insert into segmentation_service.experiment (experiment_name) values ($1) returning experiment_id", name).
		Scan(&result.Id); err != nil {
		return nil, err
	}
	result.Name = name
	if err := result.Validate(); err != nil {
		err := DeleteExperiment(name)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return &result, nil
}

func GetAllExperiment() ([]model.Experiment, error) {
	logger.Info("method <GetAllExperiment> started")
	rows, err := db.Query("select * from segmentation_service.experiment")
	if err != nil {
		return nil, err
	}
	var result []model.Experiment
	for rows.Next() {
		p := model.Experiment{}
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}
		if err := p.Validate(); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func DeleteExperiment(name string) error {
	logger.Info("method <DeleteExperiment> started")
	_, err := db.Exec("delete from segmentation_service.experiment where experiment_name = $1", name)
	if err != nil {
		return err
	}
	return nil
}

func AddEmployeeToExperiment(experimentNames []string, employeeId int) error {
	logger.Info("method <AddEmployeeToExperiment> started")
	_, err := db.Exec(getTransactionAdd(experimentNames, employeeId))
	if err != nil {
		return err
	}
	return nil
}

func RemoveEmployeeFromExperiment(experimentNames []string, employeeId int) error {
	logger.Info("method <RemoveEmployeeFromExperiment> started")
	_, err := db.Exec(getTransactionRemove(experimentNames, employeeId))
	if err != nil {
		return err
	}
	return nil
}

func getTransactionRemove(experimentNames []string, employeeId int) string {
	logger.Info("method <getTransactionRemove> started")
	return getTransaction(experimentNames, employeeId, func(experimentNames []string, employeeId int) string {
		result := ""
		for i := 0; i < len(experimentNames); i++ {
			result += fmt.Sprintf("update segmentation_service.employee_experiment set dateremoved = now() where employee_id = %d and experiment_id = (select experiment_id from segmentation_service.experiment where experiment.experiment_name = '%s');", employeeId, experimentNames[i])
		}
		return result
	})
}

func getTransactionAdd(experimentNames []string, employeeId int) string {
	logger.Info("method <getTransactionAdd> started")
	return getTransaction(experimentNames, employeeId, func(experimentNames []string, employeeId int) string {
		result := ""
		for i := 0; i < len(experimentNames); i++ {
			exp := model.Experiment{Id: -1, Name: experimentNames[i]}
			if err := exp.Validate(); err != nil {
				return ""
			}
			result += fmt.Sprintf("insert into segmentation_service.employee_experiment (employee_id, experiment_id, dateadded) values (%d, (select experiment_id from segmentation_service.experiment where experiment_name = '%s'), now());\n", employeeId, experimentNames[i])
		}
		return result
	})
}

func GetReportById(employeeId int) ([]model.Report, error) {
	logger.Info("method <GetReportById> started")
	var result []model.Report
	err := getReport(&result, fmt.Sprintf("select employee_id, experiment_name, 'add' as operation, dateadded from segmentation_service.employee_experiment inner join segmentation_service.experiment using(experiment_id) where employee_id = %d;", employeeId))
	if err != nil {
		return nil, err
	}
	err = getReport(&result, fmt.Sprintf("select employee_id, experiment_name, 'remove' as operation, dateremoved from segmentation_service.employee_experiment inner join segmentation_service.experiment using(experiment_id) where employee_id = %d;", employeeId))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AddPercentageOfEmployeeToExperiment(percent float32, experimentNames []string) error {
	logger.Info("method <AddPercentageOfEmployeeToExperiment> started")
	employees, err := GetPercentageOfEmployee(percent)
	if err != nil {
		return err
	}
	for _, p := range employees {
		err := AddEmployeeToExperiment(experimentNames, p.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddEmployeeToExperimentWithTTL(experimentNames []string, employeeId int, dateDelete string) error {
	logger.Info("method <AddEmployeeToExperimentWithTTL> started")
	_, err := db.Exec(getTransactionAddWithTTL(experimentNames, employeeId, dateDelete))
	if err != nil {
		return err
	}
	return nil
}

func getTransactionAddWithTTL(experimentNames []string, employeeId int, dateDelete string) string {
	logger.Info("method <getTransactionAddWithTTL> started")
	return getTransaction(experimentNames, employeeId, func(experimentNames []string, employeeId int) string {
		result := ""
		for i := 0; i < len(experimentNames); i++ {
			exp := model.Experiment{Id: -1, Name: experimentNames[i]}
			if err := exp.Validate(); err != nil {
				return ""
			}
			result += fmt.Sprintf("insert into segmentation_service.employee_experiment (employee_id, experiment_id, dateadded, date_delete) values (%d, (select experiment_id from segmentation_service.experiment where experiment_name = '%s'), now(), now() + %s);\n", employeeId, experimentNames[i], dateDelete)
		}
		return result
	})
}

func getReport(result *[]model.Report, query string) error {
	logger.Info("method <getReport> started")
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		p := model.Report{}
		var date sql.NullString
		err := rows.Scan(&p.EmployeeId, &p.ExperimentName, &p.Operation, &date)
		if err != nil {
			return err
		}
		if !date.Valid {
			p.Date = "null"
		} else {
			p.Date = date.String
		}
		*result = append(*result, p)
	}
	return nil
}

func getTransaction(experimentNames []string,
	employeeId int,
	impl func(experimentNames []string, employeeId int) string) string {
	result := "begin transaction;\n"
	result += impl(experimentNames, employeeId)
	result += "commit transaction;"
	fmt.Println(result)
	return result
}
