package equals

import (
	"github.com/stretchr/testify/require"
	"reflect"
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

func ElementsMatchRec(t *testing.T, want interface{}, got interface{}) {
	w := reflect.ValueOf(want)
	g := reflect.ValueOf(got)
	switch w.Type().Kind() {
	case reflect.Array, reflect.Slice:
		require.ElementsMatch(t, want, got)
	case reflect.Struct:
		for i := 0; i < w.NumField(); i++ {
			ElementsMatchRec(t, w.Field(i).Interface(), g.Field(i).Interface())
		}
	case reflect.Ptr:
		for i := 0; i < w.Elem().NumField(); i++ {
			ElementsMatchRec(t, w.Elem().Field(i).Interface(), g.Elem().Field(i).Interface())
		}
	default:
		require.Equal(t, want, got)
	}
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
