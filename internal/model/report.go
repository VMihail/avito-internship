package model

type Report struct {
	EmployeeId     int
	ExperimentName string
	Operation      string
	Date           string
}

func (report *Report) Validate() error {
	employee := Employee{Id: report.EmployeeId}
	experiment := Experiment{Id: -1, Name: report.ExperimentName}
	employeeResult := employee.Validate()
	experimentResult := experiment.Validate()
	if employeeResult != nil {
		return employeeResult
	}
	if experimentResult != nil {
		return experimentResult
	}
	return nil
}
