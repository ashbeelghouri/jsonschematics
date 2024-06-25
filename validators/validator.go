package validators

import "github.com/ashbeelghouri/jsonschematics/utils"

type Validators struct {
	ValidationFns map[string]Validator
	Logger        utils.Logger
}

type Validator func(interface{}, map[string]interface{}) error

func (v *Validators) RegisterValidator(name string, fn Validator) {
	v.Logger.DEBUG("registering validator:", name)
	if v.ValidationFns == nil {
		v.ValidationFns = make(map[string]Validator)
	}
	v.ValidationFns[name] = fn
}

func (v *Validators) BasicValidators() {
	v.Logger.DEBUG("loading all the basic validators")
	// String Validators
	v.RegisterValidator("IsString", IsString)
	v.RegisterValidator("NotEmpty", NotEmpty)
	v.RegisterValidator("StringInArr", StringInArr)
	v.RegisterValidator("IsEmail", IsEmail)
	v.RegisterValidator("MaxLengthAllowed", MaxLengthAllowed)
	v.RegisterValidator("MinLengthAllowed", MinLengthAllowed)
	v.RegisterValidator("InBetweenLengthAllowed", InBetweenLengthAllowed)
	v.RegisterValidator("NoSpecialCharacters", NoSpecialCharacters)
	v.RegisterValidator("HaveSpecialCharacters", HaveSpecialCharacters)
	v.RegisterValidator("LeastOneUpperCase", LeastOneUpperCase)
	v.RegisterValidator("LeastOneLowerCase", LeastOneLowerCase)
	v.RegisterValidator("LeastOneDigit", LeastOneDigit)
	v.RegisterValidator("IsURL", IsURL)
	v.RegisterValidator("IsNotURL", IsNotURL)
	v.RegisterValidator("HaveURLHostName", HaveURLHostName)
	v.RegisterValidator("HaveQueryParameter", HaveQueryParameter)
	v.RegisterValidator("IsHttps", IsHttps)
	v.RegisterValidator("IsURL", IsValidUuid)
	v.RegisterValidator("LIKE", LIKE)
	v.RegisterValidator("IsValidUuid", Regex)

	// Number Validators
	v.RegisterValidator("IsNumber", IsNumber)
	v.RegisterValidator("MaxAllowed", MaxAllowed)
	v.RegisterValidator("MinAllowed", MinAllowed)
	v.RegisterValidator("InBetween", InBetween)

	// Date Validators
	v.RegisterValidator("IsValidDate", IsValidDate)
	v.RegisterValidator("IsLessThanNow", IsLessThanNow)
	v.RegisterValidator("IsMoreThanNow", IsMoreThanNow)
	v.RegisterValidator("IsBefore", IsBefore)
	v.RegisterValidator("IsAfter", IsAfter)
	v.RegisterValidator("IsInBetweenTime", IsInBetweenTime)

	//Arrays
	v.RegisterValidator("ArrayLengthMax", ArrayLengthMax)
	v.RegisterValidator("ArrayLengthMin", ArrayLengthMin)

	v.Logger.DEBUG("basic validators loaded")
}
