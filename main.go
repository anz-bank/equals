package equals

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

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