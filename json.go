package equals

import (
	"encoding/json"
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
