package operators

type Operators struct {
	OpFunctions map[string]Op
}

type Op func(interface{}, map[string]interface{}) error

func (v *Operators) RegisterValidator(name string, fn Op) {
	if v.OpFunctions == nil {
		v.OpFunctions = make(map[string]Op)
	}
	v.OpFunctions[name] = fn
}
