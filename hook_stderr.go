package qlog

import (
	"os"
	"reflect"
)

const (
	keyStderrEnabled = "logger.stderr.enabled"
	keyStderrLevel   = "logger.stderr.level"
)

type StderrHook struct {
	BaseHook
}

func (h *StderrHook) Setup() error {
	h.baseSetup()

	h.writer = os.Stderr

	return nil
}

var _InitStderrHook = func() interface{} {
	gCommandLine.Bool(keyStderrEnabled, true, "logger.stderr.enabled")
	gCommandLine.String(keyStderrLevel, "error", "logger.stderr.level")

	registerHook("stderr", reflect.TypeOf(StderrHook{}))
	return nil
}()
