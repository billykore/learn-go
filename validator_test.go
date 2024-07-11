package learning

import (
	"runtime"
	"testing"

	"github.com/go-playground/validator/v10"
)

type user struct {
	KTPName  string
	FullName string `validate:"mustEqualIgnoreCase=KTPName"`
}

func mustEqualIgnoreCase(field validator.FieldLevel) bool {
	value, _, _, ok := field.GetStructFieldOK2()
	if !ok {
		panic("field not ok")
	}

	firstValue := field.Field().String()
	secondValue := value.String()

	return firstValue == secondValue
}

func TestCrossValidation(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("mustEqualIgnoreCase", mustEqualIgnoreCase)
	if err != nil {
		panic(err)
	}

	u := &user{
		KTPName:  "Kore",
		FullName: "Kore  ",
	}

	err = validate.Struct(u)
	if err != nil {
		panic(err)
	}
}

const defaultSkipper = 0

func TestRuntime(t *testing.T) {
	pc, file, line, ok := runtime.Caller(defaultSkipper)
	t.Log(pc)
	t.Log(file)
	t.Log(line)
	t.Log(ok)

	fn := runtime.FuncForPC(pc)
	t.Log(fn.Name())
}
