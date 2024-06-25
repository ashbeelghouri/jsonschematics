package operators

import "github.com/ashbeelghouri/jsonschematics/utils"

type Operators struct {
	OpFunctions map[string]Op
	Logger      utils.Logger
}

type Op func(interface{}, map[string]interface{}) *interface{}

func (op *Operators) RegisterOperation(name string, fn Op) {
	op.Logger.DEBUG("registering operation:", name)
	if op.OpFunctions == nil {
		op.OpFunctions = make(map[string]Op)
	}
	op.OpFunctions[name] = fn
}

func (op *Operators) LoadBasicOperations() {
	op.Logger.DEBUG("loading basic operations")
	op.RegisterOperation("Capitalize", Capitalize)
	op.RegisterOperation("UpperCase", UpperCase)
	op.RegisterOperation("LowerCase", LowerCase)

	// number operations
	op.RegisterOperation("Add", Add)
	op.RegisterOperation("Subtract", Subtract)
	op.RegisterOperation("Multiply", Multiply)
	op.RegisterOperation("Divide", Divide)

	op.Logger.DEBUG("basic operations loaded")
}
