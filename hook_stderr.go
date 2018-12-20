package qlog

import (
	"os"
	"reflect"
)

type StderrHook struct {
	BaseHook
}

func (h *StderrHook) Setup() error {
	var err error

	if err = v.UnmarshalKey("logger.stderr", h); err != nil {
		return err
	}

	h.baseSetup()

	h.writer = os.Stderr

	return nil
}

func init() {
	registerHook("stderr", reflect.TypeOf(StdoutHook{}))
}
