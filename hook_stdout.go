package qlog

import (
	"os"
	"reflect"
)

type StdoutHook struct {
	BaseHook
}

func (h *StdoutHook) Setup() error {
	var err error

	if err = v.UnmarshalKey("logger.stdout", h); err != nil {
		return err
	}

	h.baseSetup()

	h.writer = os.Stdout

	return nil
}

func init() {
	registerHook("stdout", reflect.TypeOf(StdoutHook{}))
}
