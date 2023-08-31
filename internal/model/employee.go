package model

import validation "github.com/go-ozzo/ozzo-validation"

type Employee struct {
	Id int
}

func (employee *Employee) Validate() error {
	return validation.ValidateStruct(employee, validation.Field(&employee.Id, validation.Required))
}
