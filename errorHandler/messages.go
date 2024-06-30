package errorHandler

import (
	"errors"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics"
	"strings"
)

type Locale string
type Target string

type Error struct {
	DataTarget string
	Message    map[Locale]string
	Validator  string
	Value      interface{}
	ID         interface{}
}

type Errors struct {
	Messages map[Target]Error
}

func (e *Error) AddMessage(local string, message string) {
	e.Message[Locale(local)] = message
}

func (em *Errors) AddError(target string, err Error) {
	if em.Messages == nil {
		em.Messages = make(map[Target]Error)
	}
	em.Messages[Target(target)] = err
}

func (em *Errors) HasErrors() bool {
	for _, err := range em.Messages {
		if len(err.Message) > 0 {
			return true
		}
	}
	return false
}

func (em *Errors) GetStrings(locale Locale, format string) *[]string {
	var errs []string
	if !em.HasErrors() {
		return nil
	}
	if format == "" {
		format = "validation error %message for %target with validation on %validator, provided: %value"
	}

	for target, msg := range em.Messages {
		message, ok := msg.Message[locale]
		if !ok {
			continue
		}
		value := fmt.Sprint(msg.Value)
		var id *string
		if msg.ID != nil {
			msgID := fmt.Sprint(msg.ID)
			id = &msgID
		} else {
			id = nil
		}
		errs = append(errs, jsonschematics.FormatError(id, message, string(target), msg.Validator, value, format))
	}
	return &errs
}

func (em *Errors) GetErrors(locale Locale, format string) *[]error {
	var errs []error
	if !em.HasErrors() {
		return nil
	}
	if format == "" {
		format = "validation error %message for %target with validation on %validator, provided: %value"
	}

	for target, msg := range em.Messages {
		message, ok := msg.Message[locale]
		if !ok {
			continue
		}
		value := fmt.Sprint(msg.Value)
		var id *string
		if msg.ID != nil {
			msgID := fmt.Sprint(msg.ID)
			id = &msgID
		} else {
			id = nil
		}
		errs = append(errs, errors.New(jsonschematics.FormatError(id, message, string(target), msg.Validator, value, format)))
	}
	return &errs
}

func (em *Errors) GetJoinedError(locale string, singleErrorFormat string, appendWith string) error {
	errorStrings := em.GetStrings(Locale(locale), singleErrorFormat)
	if errorStrings == nil {
		return nil
	}
	if errorStrings != nil && !(len(*errorStrings) > 1) {
		return errors.New(strings.Join(*errorStrings, ""))
	}
	if appendWith == "" {
		appendWith = ","
	}
	return errors.New(strings.Join(*errorStrings, appendWith))
}

func (em *Errors) MergeErrors(em2 *Errors) {
	if !(em.HasErrors() && em2.HasErrors()) {
		return
	}
	if em.Messages == nil {
		em.Messages = make(map[Target]Error)
	}
	for target, err := range em2.Messages {
		em.Messages[target] = err
	}
}
