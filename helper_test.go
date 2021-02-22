package equals

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
