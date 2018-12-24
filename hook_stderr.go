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
	cli.Bool(keyStderrEnabled, false, "logger.stderr.enabled")
	cli.String(keyStderrLevel, "", "logger.stderr.level") // DONOT set default level in pflag

	registerHook("stderr", reflect.TypeOf(StderrHook{}))
	return nil
}()
