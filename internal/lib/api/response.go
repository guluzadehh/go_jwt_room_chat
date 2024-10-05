package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	*Error `json:"error,omitempty"`
}

type Error struct {
	Message string      `json:"message"`
	Details []ErrDetail `json:"details,omitempty"`
}

type ErrDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
		Error:  nil,
	}
}

func Err(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  &Error{Message: msg},
	}
}

func ErrD(msg string, details []ErrDetail) Response {
	res := Err(msg)
	res.Error.Details = details
	return res
}

func alias(name string) string {
	switch name {
	case "Username":
		return "username"
	case "Password":
		return "password"
	case "ConfPassword":
		return "confirm password"
	default:
		return name
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errDetails []ErrDetail

	for _, err := range errs {
		field := alias(err.Field())

		var msg string

		switch err.ActualTag() {
		case "required":
			msg = fmt.Sprintf("field %s is required.", field)
		case "min":
			msg = fmt.Sprintf("field %s min length must be %s.", field, alias(err.Param()))
		case "max":
			msg = fmt.Sprintf("field %s max length must be %s.", field, alias(err.Param()))
		case "alphanum":
			msg = fmt.Sprintf("field %s must contain both letters and numbers.", field)
		case "containsany":
			msg = fmt.Sprintf("field %s must contain on of the following characters: %s.", field, alias(err.Param()))
		case "eqfield":
			msg = fmt.Sprintf("field %s is not equal to %s field.", field, alias(err.Param()))
		case "passwordpattern":
			msg = "field password must contain at least one letter, one number, and one special character."
		default:
			msg = fmt.Sprintf("field %s is invalid", field)
		}

		errDetails = append(errDetails, ErrDetail{Field: field, Message: msg})
	}

	return ErrD("validation error", errDetails)
}
