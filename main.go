package equals

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ElementsMatchRec(t require.TestingT, want interface{}, got interface{}) (equal bool) {
	if !ElementsMatchRec2(t, want, got){
		assert.Fail(t, fmt.Sprintf("Not equal: %v, %v", want, got))
		return false
	}
	return true
}

func ElementsMatchRec2(t require.TestingT, want interface{}, got interface{}) (equal bool) {
	var wantj, gotj = interfaceToJson(want, got)
	return elementsMatchRecHelper(t, wantj, gotj)
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// require.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func ElementsMatch(t require.TestingT, listA, listB interface{}, msgAndArgs ...interface{}) (ok bool) {
	if isEmpty(listA) && isEmpty(listB) {
		return true
	}

	if !isList(t, listA, msgAndArgs...) || !isList(t, listB, msgAndArgs...) {
		return false
	}

	extraA, extraB := diffLists(listA, listB)

	if len(extraA) == 0 && len(extraB) == 0 {
		return true
	}

	return assert.Fail(t, formatListDiff(listA, listB, extraA, extraB), msgAndArgs...)
}

func elementsMatchRecHelper(t require.TestingT, want interface{}, got interface{}) (equal bool) {
	defer func() {
		if err := recover(); err != nil {
			equal = false
		}
	}()
	w := reflect.ValueOf(want)
	g := reflect.ValueOf(got)
	switch w.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return ElementsMatch(t, want, got)
	case reflect.Map:
		for _, e := range w.MapKeys() {
			if welem, gelem := w.MapIndex(e), g.MapIndex(e); welem.CanInterface() && gelem.CanInterface() {
				if !elementsMatchRecHelper(t, welem.Interface(), gelem.Interface()) {
					return false
				}
			} else {
				if !ElementsMatch(t, welem, gelem) {
					return false
				}
			}
		}
		return w.Len() == g.Len()
	default:
		return AssertJson(t, want, got)
	}
}

// objectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		if !equalJson(expected, actual) {
			return ElementsMatchRec2(RequireNull{}, expected, actual)
		}
		return true
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
