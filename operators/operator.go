package operators

type Operators struct {
	OpFunctions map[string]Op
}

type Op func(interface{}, map[string]interface{}) *interface{}

func (op *Operators) RegisterOperation(name string, fn Op) {
	if op.OpFunctions == nil {
		op.OpFunctions = make(map[string]Op)
	}
	op.OpFunctions[name] = fn
}

func (op *Operators) LoadBasicOperations() {
	op.RegisterOperation("Capitalize", Capitalize)
	op.RegisterOperation("UpperCase", UpperCase)
	op.RegisterOperation("LowerCase", LowerCase)
}
