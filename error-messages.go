package jsonschematics

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
