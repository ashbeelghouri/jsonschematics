package jsonschematics

import (
	"errors"
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
	em.Messages = append(em.Messages, ErrorMessage{Message: err, Validator: validator, Target: target, Value: value, ID: nil})
}

func (em *ErrorMessages) AddErrorsForArray(validator string, target string, err string, value interface{}, id interface{}) {
	em.Messages = append(em.Messages, ErrorMessage{Message: err, Validator: validator, Target: target, Value: value, ID: id})
}

func (em *ErrorMessages) HaveErrors() bool {
	if len(em.Messages) > 0 {
		return true
	}
	return false
}

/*
	format: "validation error %message for %target with validating with %validation, provided: %value"
*/

func (em *ErrorMessages) HaveSingleError(format string, appendWith string) error {
	if format == "" {
		format = "validation error %message for %target with validation on %validator, provided: %value"
	}
	if em == nil {
		return nil
	}

	if !(len(em.Messages) > 1) {
		for _, msg := range em.Messages {
			var id *string
			if msg.ID != nil {
				msgID := msg.ID.(string)
				id = &msgID
			} else {
				id = nil
			}
			return errors.New(FormatError(id, msg.Message, msg.Target, msg.Validator, msg.Value.(string), format))
		}
		return nil
	} else {
		if appendWith == "" {
			appendWith = ","
		}
		errorMessages := make([]string, len(em.Messages))
		for _, msg := range em.Messages {
			var id *string
			if msg.ID != nil {
				msgID := msg.ID.(string)
				id = &msgID
			} else {
				id = nil
			}
			errorMessages = append(errorMessages, FormatError(id, msg.Message, msg.Target, msg.Validator, msg.Value.(string), format))
		}
		return errors.New(strings.Join(errorMessages, appendWith))
	}
}
