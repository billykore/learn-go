package learning

import "testing"

type adder[T int | float64] interface {
	add(x T, y T) T
}

type floatAdder struct{}

func (a floatAdder) add(x, y float64) float64 {
	return x + y
}

type intAdder struct{}

func (a intAdder) add(x, y int) int {
	return x + y
}

type calculator[T int | float64] struct {
	adder adder[T]
}

func (c calculator[T]) setAdder(adder adder[T]) {
	c.adder = adder
}

func (c calculator[T]) add(x, y T) T {
	return c.adder.add(x, y)
}

func TestAdder(t *testing.T) {
	c1 := calculator[int]{}
	a1 := intAdder{}
	c1.setAdder(a1)
	t.Log(a1.add(1, 2))

	c2 := calculator[float64]{}
	a2 := floatAdder{}
	c2.setAdder(a2)
	t.Log(c2.add(1.5, 1.5))
}
