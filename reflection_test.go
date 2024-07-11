package learning

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Sample struct {
	Name string
}

func TestReflection(t *testing.T) {
	sample := Sample{Name: "Oyen"}
	sampleType := reflect.TypeOf(sample)
	structField := sampleType.Field(0)
	fmt.Println(structField)
}

type fooBar struct {
	Foo string
	Bar int
}

func TestReflect(t *testing.T) {
	x := fooBar{"", 2}

	v := reflect.ValueOf(x)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	for _, value := range values {
		if reflect.TypeOf(value).Name() == "string" && value == "" {
			fmt.Println("nil")
		} else {
			fmt.Println(value)
		}
	}
	fmt.Println(values)
}

func TestBitwise(t *testing.T) {
	var x byte = 5
	var y byte = 3
	var z = x ^ y

	fmt.Printf("x: %04b\n", x)
	fmt.Printf("y: %04b\n", y)
	fmt.Printf("z: %04b\n", z)
}

func TestShift(t *testing.T) {
	var x byte = 5
	var y byte = x << 2
	var z byte = x >> 2

	fmt.Printf("x: %08b = %d\n", x, x)
	fmt.Printf("y: %08b = %d\n", y, y)
	fmt.Printf("z: %08b = %d\n", z, z)
}

func isPowerOfTwo(x int) bool {
	return x&(x-1) == 0
}

func TestIsPowerOfTwo(t *testing.T) {
	assert.True(t, isPowerOfTwo(16))
}

type A struct {
	name string
}

func (a *A) setName(name string) *A {
	a.name = name
	return a
}

func TestPassByValue(t *testing.T) {
	a := A{name: "Oyen"}
	b := a.setName("Kore")

	t.Log(b.name)
	t.Log(a.name)
}
