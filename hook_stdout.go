package qlog

import (
	"os"
	"reflect"
)

type StdoutHook struct {
	BaseHook
}

func (s *StdoutHook) Setup() {
	s.baseSetup()
	s.writer = os.Stdout
}

func init() {
	registerHook("stdout", reflect.TypeOf(StdoutHook{}))
}
