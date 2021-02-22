package equals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
)

/* AssertJson is a function that checks two dereferenced proto objects because require.Equal infinite loops */
func AssertJson(t require.TestingT, want interface{}, got interface{}) bool {
	expectJson, err := json.Marshal(want)
	if !assert.NoError(t, err) {
		return false
	}
	gotJson, err := json.Marshal(got)
	if !assert.NoError(t, err) {
		return false
	}
	return assert.JSONEq(t, string(expectJson), string(gotJson))
}

func ElementsMatchRec2(t require.TestingT, want interface{}, got interface{}) (equal bool) {
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
				if !ElementsMatchRec2(t, welem.Interface(), gelem.Interface()) {
					return false
				}
			} else {
				if !ElementsMatch(t, welem, gelem) {
					return false
				}
			}
		}
		return true
	default:
		return AssertJson(t, want, got)
	}
	return equal
}

func ElementsMatchRec(t require.TestingT, want interface{}, got interface{}) (equal bool) {
	var wantj, gotj = interfaceToJson(want, got)
	return ElementsMatchRec2(t, wantj, gotj)
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

// diffLists diffs two arrays/slices and returns slices of elements that are only in A and only in B.
// If some element is present multiple times, each instance is counted separately (e.g. if something is 2x in A and
// 5x in B, it will be 0x in extraA and 3x in extraB). The order of items in both lists is ignored.
func diffLists(listA, listB interface{}) (extraA, extraB []interface{}) {
	aValue := reflect.ValueOf(listA)
	bValue := reflect.ValueOf(listB)

	aLen := aValue.Len()
	bLen := bValue.Len()

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := aValue.Index(i).Interface()
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			element2 := bValue.Index(j).Interface()
			if objectsAreEqual(element2, element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, bValue.Index(j).Interface())
	}

	return
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
			return ElementsMatchRec(RequireNull{}, expected, actual)
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

// jsoneq asserts that two JSON strings are equivalent.
//
//  require.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func jsoneq(expected string, actual string, msgAndArgs ...interface{}) bool {
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		return false
	}

	return reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface)
}

/* equalJson is a function that checks two dereferenced proto objects because require.Equal infinite loops */
func equalJson(want interface{}, got interface{}) bool {
	expectJson, err := json.Marshal(want)
	if err != nil {
		return false
	}
	gotJson, err := json.Marshal(got)
	if err != nil {
		return false
	}
	return jsoneq(string(expectJson), string(gotJson))
}

/* equalJson is a function that checks two dereferenced proto objects because require.Equal infinite loops */
func interfaceToJson(want interface{}, got interface{}) (interface{}, interface{}) {
	expectJson, _ := json.Marshal(want)
	gotJson, _ := json.Marshal(got)
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal(expectJson, &expectedJSONAsInterface); err != nil {
		return nil, nil
	}

	if err := json.Unmarshal(gotJson, &actualJSONAsInterface); err != nil {
		return nil, nil
	}

	return expectedJSONAsInterface, actualJSONAsInterface
}

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {

	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
		// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
		// for all other types, compare against the zero value
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

// isList checks that the provided value is array or slice.
func isList(t require.TestingT, list interface{}, msgAndArgs ...interface{}) (ok bool) {
	kind := reflect.TypeOf(list).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return assert.Fail(t, fmt.Sprintf("%q has an unsupported type %s, expecting array or slice", list, kind),
			msgAndArgs...)
	}
	return true
}

func formatListDiff(listA, listB interface{}, extraA, extraB []interface{}) string {
	var msg bytes.Buffer

	msg.WriteString("elements differ")
	if len(extraA) > 0 {
		msg.WriteString("\n\nextra elements in list A:\n")
		msg.WriteString(spewConfig.Sdump(extraA))
	}
	if len(extraB) > 0 {
		msg.WriteString("\n\nextra elements in list B:\n")
		msg.WriteString(spewConfig.Sdump(extraB))
	}
	msg.WriteString("\n\nlistA:\n")
	msg.WriteString(spewConfig.Sdump(listA))
	msg.WriteString("\n\nlistB:\n")
	msg.WriteString(spewConfig.Sdump(listB))

	return msg.String()
}

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
	DisableMethods:          true,
	MaxDepth:                10,
}
