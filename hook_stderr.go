package qlog

import (
	"os"
	"reflect"
)

type StderrHook struct {
	BaseHook
}

func (s *StderrHook) Setup() error {
	var err error

	if err = v.UnmarshalKey("logger.stderr", s); err != nil {
		return err
	}

	s.baseSetup()

	s.writer = os.Stderr

	return nil
}

func init() {
	registerHook("stderr", reflect.TypeOf(StdoutHook{}))
}
