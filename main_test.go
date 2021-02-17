package equals

import (
	"testing"
)

func TestStruct(t *testing.T) {
	type Bar struct {
		Arr []int
	}
	type Foo struct {
		A Bar
	}
	ElementsMatchRec(t, Foo{A: Bar{Arr: []int{1, 2, 3}}}, Foo{A: Bar{Arr: []int{3, 2, 1}}})
}

func TestPointer(t *testing.T) {
	type Bar struct {
		Arr []int
	}
	type Foo struct {
		A *Bar
	}
	ElementsMatchRec(t, &Foo{A: &Bar{Arr: []int{1, 2, 3}}}, &Foo{A: &Bar{Arr: []int{3, 2, 1}}})
}

func TestArray(t *testing.T) {
	ElementsMatchRec(t, []int{1, 2, 3}, []int{3, 2, 1})
}

func TestPrimitive(t *testing.T) {
	ElementsMatchRec(t, 1, 1)
}

func TestStructArray(t *testing.T) {
	type Foobar struct {
		I int
	}
	type Bar struct {
		Arr []Foobar
	}
	type Foo struct {
		A Bar
	}
	ElementsMatchRec(t, Foo{A: Bar{Arr: []Foobar{{1}, {2}, {3}}}}, Foo{A: Bar{Arr: []Foobar{{2}, {1}, {3}}}})
}
