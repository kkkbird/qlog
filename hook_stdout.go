package qlog

import (
	"os"
	"reflect"
)

type StdoutHook struct {
	BaseHook
}

func (s *StdoutHook) Setup() error {
	var err error

	if err = v.UnmarshalKey("logger.stdout", s); err != nil {
		return err
	}

	s.baseSetup()

	s.writer = os.Stdout

	return nil
}

func init() {
	registerHook("stdout", reflect.TypeOf(StdoutHook{}))
}
