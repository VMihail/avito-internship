package model

import validation "github.com/go-ozzo/ozzo-validation"

type Experiment struct {
	Id   int
	Name string
}

func (experiment *Experiment) Validate() error {
	return validation.ValidateStruct(experiment,
		validation.Field(&experiment.Id, validation.Required),
		validation.Field(&experiment.Name, validation.Length(3, 100)))
}
