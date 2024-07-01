package errorHandler

import (
	"errors"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"log"
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
	Data       map[string]interface{}
}

type Errors struct {
	Messages map[Target]Error
}

func (e *Error) AddMessage(local string, message string) {
	if e.Message == nil {
		e.Message = make(map[Locale]string)
	}
	e.Message[Locale(local)] = message
}
func (e *Error) updateData(target string) Target {
	var t string
	convertedID, ok := e.ID.(string)

	if ok && e.ID != nil {
		t = fmt.Sprintf("%s:%s", convertedID, target)
	} else {
		t = fmt.Sprintf("%s", target)
	}
	e.Data = make(map[string]interface{})
	e.Data["target"] = t
	e.Data["messages"] = e.Message
	e.Data["validator"] = e.Validator
	e.Data["value"] = e.Value
	e.Data["value"] = e.Value
	e.Data["id"] = e.ID
	return Target(t)
}

func (em *Errors) AddError(target string, err Error) {
	if em.Messages == nil {
		em.Messages = make(map[Target]Error)
	}
	t := err.updateData(target)
	em.Messages[t] = err
}

func (em *Errors) HasErrors() bool {
	if em != nil {
		for _, err := range em.Messages {
			if len(err.Message) > 0 {
				return true
			}
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
		format = "validation error %message for %target with validation on %validator, provided: %value: {%data}"
	}

	for target, msg := range em.Messages {
		log.Println(target)
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
		errs = append(errs, utils.FormatError(id, message, string(target), msg.Validator, value, format, &msg.Data))
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
		log.Println(target)
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
		errs = append(errs, errors.New(utils.FormatError(id, message, string(target), msg.Validator, value, format, &msg.Data)))
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
	if !em2.HasErrors() {
		return
	}
	if em.Messages == nil {
		em.Messages = make(map[Target]Error)
	}
	for target, err := range em2.Messages {
		em.Messages[target] = err
	}
}
