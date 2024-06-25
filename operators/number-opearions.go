package operators

func Add(i interface{}, attr map[string]interface{}) *interface{} {
	num := i.(float64)
	add := attr["add_with"].(float64)
	var result interface{} = num + add
	return &result
}

func Subtract(i interface{}, attr map[string]interface{}) *interface{} {
	num := i.(float64)
	sub := attr["subtract_with"].(float64)
	var result interface{} = num - sub
	return &result
}

func Multiply(i interface{}, attr map[string]interface{}) *interface{} {
	num := i.(float64)
	mul := attr["multiply_with"].(float64)
	var result interface{} = num * mul
	return &result
}

func Divide(i interface{}, attr map[string]interface{}) *interface{} {
	num := i.(float64)
	divide := attr["divide_with"].(float64)
	var result interface{} = num / divide
	return &result
}
