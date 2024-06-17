package jsonschematics

type ErrorMessage struct {
	Message   string
	Validator string
	Target    string
	Value     interface{}
}

type ErrorMessages struct {
	Messages []ErrorMessage
}

type ArrayOfErrors struct {
	Errors ErrorMessages
	ID     interface{}
}

func (em *ErrorMessages) AddError(validator string, target string, err string, value interface{}) {
	em.Messages = append(em.Messages, ErrorMessage{Message: err, Validator: validator, Target: target, Value: value})
}

func (em *ErrorMessages) HaveErrors() bool {
	if len(em.Messages) > 0 {
		return true
	}
	return false
}
