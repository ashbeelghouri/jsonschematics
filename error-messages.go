package jsonschematics

type ErrorMessage struct {
	Message   string
	Validator string
	Target    string
}

type ErrorMessages struct {
	Messages []ErrorMessage
}

func (em *ErrorMessages) AddError(validator string, target string, e string) {
	em.Messages = append(em.Messages, ErrorMessage{Message: e, Validator: validator, Target: target})
}

func (em *ErrorMessages) HaveErrors() bool {
	if len(em.Messages) > 0 {
		return true
	}
	return false
}
