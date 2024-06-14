package jsonschematics

import "errors"

type ErrorMessages struct {
	Messages []string
}

func (em *ErrorMessages) AddError(e string) {
	em.Messages = append(em.Messages, e)
}

func (em *ErrorMessages) HaveErrors() bool {
	if len(em.Messages) > 0 {
		return true
	}
	return false
}

func (em *ErrorMessages) Errors() []error {
	var messages []error
	for _, e := range em.Messages {
		messages = append(messages, errors.New(e))
	}
	return messages
}
