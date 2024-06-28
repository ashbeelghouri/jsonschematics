package jsonschematics

import (
	"errors"
	"fmt"
	"strings"
)

// make error format same for the arrays as well as objects

type ErrorMessage struct {
	Message   string
	Validator string
	Target    string
	Value     interface{}
	ID        interface{}
}

type ErrorMessages struct {
	Messages []ErrorMessage
}

func (em *ErrorMessages) AddError(validator string, target string, err string, value interface{}) {
	logs.DEBUG("adding new error message", err)
	em.Messages = append(em.Messages, ErrorMessage{Message: err, Validator: validator, Target: target, Value: value, ID: nil})
}

func (em *ErrorMessages) AddErrorsForArray(validator string, target string, err string, value interface{}, id interface{}) {
	logs.DEBUG("adding error for arrays", err, "on id: ", id)
	em.Messages = append(em.Messages, ErrorMessage{Message: err, Validator: validator, Target: target, Value: value, ID: id})
}

func (em *ErrorMessages) HaveErrors() bool {
	return len(em.Messages) > 0
}

func (em *ErrorMessages) ExtractAsStrings(format string) *[]string {
	logs.DEBUG("extracting errors as a string")
	var errs []string
	if !em.HaveErrors() {
		return nil
	}
	if format == "" {
		format = "validation error %message for %target with validation on %validator, provided: %value"
	}

	for _, msg := range em.Messages {
		value := fmt.Sprint(msg.Value)
		var id *string
		if msg.ID != nil {
			msgID := fmt.Sprint(msg.ID)
			id = &msgID
		} else {
			id = nil
		}
		errs = append(errs, FormatError(id, msg.Message, msg.Target, msg.Validator, value, format))
	}

	return &errs
}

func (em *ErrorMessages) ExtractAsErrors(format string) []error {
	logs.DEBUG("extracting errors as array of errors")
	if !em.HaveErrors() {
		return nil
	}
	var errs []error

	messages := em.ExtractAsStrings(format)
	for _, msg := range *messages {
		errs = append(errs, errors.New(msg))
	}

	return errs
}

/*
	format: "validation error %message for %target with validating with %validation, provided: %value"
*/

func (em *ErrorMessages) HaveSingleError(format string, appendWith string) error {
	logs.DEBUG("joining all the errors to represent only one error")

	if !em.HaveErrors() {
		return nil
	}
	err := em.ExtractAsStrings(format)
	if err != nil && !(len(*err) > 1) {
		return errors.New(strings.Join(*err, ""))
	} else if err != nil && len(*err) > 1 {
		if appendWith == "" {
			appendWith = ","
		}
		return errors.New(strings.Join(*err, appendWith))
	} else if err != nil {
		logs.ERROR("[code=1] We are unable to determine the error :::: >>>> ", err)
		return errors.New("unable to determine the error")
	}

	return nil
}
