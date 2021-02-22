package equals

import (
	"github.com/davecgh/go-spew/spew"
)

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
	DisableMethods:          true,
	MaxDepth:                10,
}

type RequireFail struct {
	HasErrored bool
	HasFailed  bool
}

func (s *RequireFail) Errorf(format string, args ...interface{}) {
	s.HasErrored = true
}

func (s *RequireFail) FailNow() {
	s.HasFailed = true
}

type RequireNull struct {
}

func (RequireNull) Errorf(format string, args ...interface{}) {
}

func (RequireNull) FailNow() {
}
