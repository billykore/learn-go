package learning

import (
	"fmt"
	"runtime"
	"slices"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
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

func TestNumeric(t *testing.T) {
	validate := validator.New()
	num := ""
	err := validate.Var(num, "number")
	assert.Error(t, err)
}

func TestOmitEmpty(t *testing.T) {
	validate := validator.New()
	num := ""
	err := validate.Var(num, "omitempty,number")
	assert.NoError(t, err)
}

func IsValidStartDateIf(input validator.FieldLevel) bool {
	fieldValue, _, _, ok := input.GetStructFieldOKAdvanced2(input.Top(), "StartDate")
	if !ok {
		return false
	}
	fmt.Println(fieldValue)

	return false
}

type ss struct {
	Frequency string
	StartDate string `validate:"isStartDateIf=Frequency monthly"`
}

func TestIsValidStartDateIf(t *testing.T) {
	v := validator.New()

	err := v.RegisterValidation("isStartDateIf", IsValidStartDateIf)
	assert.NoError(t, err)

	err = v.Struct(ss{
		Frequency: "monthly",
		StartDate: "2024-10-25",
	})
	assert.NoError(t, err)
}

func TestContains(t *testing.T) {
	s := []int{1, 2, 3}
	assert.False(t, slices.Contains(s, 0))

	s2 := []string{"asdf", "zxcvb", "qwerty"}
	assert.True(t, slices.Contains(s2, "qwerty"))
}
