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
	cli.Bool(keyStdoutEnabled, false, "logger.stdout.enabled")
	cli.String(keyStdoutLevel, "", "logger.stdout.level") // DONOT set default level in pflag

	registerHook("stdout", reflect.TypeOf(StdoutHook{}))
	return nil
}()
