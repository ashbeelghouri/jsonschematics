package validators

type Validators struct {
	ValidationFns map[string]Validator
}

type Validator func(interface{}, map[string]interface{}) error

func LoadValidator() *Validators {
	var validators = Validators{}
	validators.BasicValidators()
	return &validators
}

func (v *Validators) RegisterValidator(name string, fn Validator) {
	if v.ValidationFns == nil {
		v.ValidationFns = make(map[string]Validator)
	}
	v.ValidationFns[name] = fn
}

func (v *Validators) BasicValidators() {
	v.RegisterValidator("IsString", IsString)
	v.RegisterValidator("IsInt", IsInt)
}
