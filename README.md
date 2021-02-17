# Equals

Functions for asserting structs are the same

+ Good for working with proto messages
+ Good for working with nested non-deterministic slices

- Bad for working with unexported fields (they are ignored)

This package uses functions from testify/require but unmarshals into Json before asserting equality, so it works well
with asserting REST/GRPC objects are the same.

Example:

```go
package equals

import (
	"github.com/joshcarp/equals"
	"testing"
)

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
	equals.ElementsMatchRec(t, Foo{A: Bar{Arr: []Foobar{{I: 1}, {I: 2}, {I: 3}}}}, Foo{A: Bar{Arr: []Foobar{{I: 2}, {I: 1}, {I: 3}}}})
}
```