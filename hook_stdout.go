package qlog

import (
	"os"
	"reflect"
)

const (
	keyStdoutEnabled = "logger.stdout.enabled"
	keyStdoutLevel   = "logger.stdout.level"
)

type StdoutHook struct {
	BaseHook
}

func (h *StdoutHook) Setup() error {
	h.baseSetup()

	h.writer = os.Stdout

	return nil
}

var _InitStdoutHook = func() interface{} {
	gCommandLine.Bool(keyStdoutEnabled, true, "logger.stdout.enabled")
	gCommandLine.String(keyStdoutLevel, "error", "logger.stdout.level")

	registerHook("stdout", reflect.TypeOf(StdoutHook{}))
	return nil
}()
