package equals

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStruct(t *testing.T) {
	type Bar struct {
		Arr []int
	}
	type Foo struct {
		A Bar
	}
	ElementsMatchRec(t,
		Foo{A: Bar{Arr: []int{1, 2, 3}}},
		Foo{A: Bar{Arr: []int{3, 2, 1}}})
}

func TestPointer(t *testing.T) {
	type Bar struct {
		Arr []int
	}
	type Foo struct {
		A *Bar
	}
	ElementsMatchRec(t,
		&Foo{A: &Bar{Arr: []int{1, 2, 3}}},
		&Foo{A: &Bar{Arr: []int{3, 2, 1}}})
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
	ElementsMatchRec(t,
		Foo{A: Bar{Arr: []Foobar{{1}, {2}, {3}}}},
		Foo{A: Bar{Arr: []Foobar{{2}, {1}, {3}}}})
}

func TestStructArray2(t *testing.T) {
	type Foobar struct {
		I int
	}
	type Bar struct {
		Arr []Foobar
	}
	type Foo struct {
		A Bar
	}
	ElementsMatchRec(t,
		Foo{A: Bar{Arr: []Foobar{{I: 1}, {I: 2}, {I: 3}}}},
		Foo{A: Bar{Arr: []Foobar{{I: 2}, {I: 1}, {I: 3}}}})
}

func TestStructArray3(t *testing.T) {
	type Foobar struct {
		I int
	}
	type Bar struct {
		Arr []*Foobar
	}
	type Foo struct {
		A Bar
	}
	failt := &RequireFail{}
	ElementsMatchRec(failt,
		Foo{A: Bar{Arr: []*Foobar{{I: 3}, {I: 1}, {I: 3}}}},
		Foo{A: Bar{Arr: []*Foobar{{I: 2}, {I: 1}, {I: 3}}}})
	require.True(t, failt.HasErrored)
}

func TestStructArray4(t *testing.T) {
	type Foobar2 struct {
		M string
	}
	type Foobar struct {
		I []Foobar2
	}
	type Bar struct {
		Arr []*Foobar
	}
	type Foo struct {
		A Bar
	}
	failt := &RequireFail{}
	ElementsMatchRec(failt,
		Foo{A: Bar{Arr: []*Foobar{{I: []Foobar2{{M: "AA"}}}}}},
		Foo{A: Bar{Arr: []*Foobar{{I: []Foobar2{{M: "A"}}}}}})
	require.True(t, failt.HasErrored)
}

func TestStructArrayInterface(t *testing.T) {
	type Foobar2 struct {
		M interface{}
	}
	type Foobar1 struct {
		M interface{}
	}
	type Foobar struct {
		I []interface{}
	}
	type Bar struct {
		Arr []*Foobar
	}
	type Foo struct {
		A Bar
	}
	ElementsMatchRec(t,
		Foo{A: Bar{Arr: []*Foobar{{I: []interface{}{Foobar2{"AA"}}}}}},
		Foo{A: Bar{Arr: []*Foobar{{I: []interface{}{Foobar1{"AA"}}}}}})
}
