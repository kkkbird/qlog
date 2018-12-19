package qlog

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type HookSetuper interface {
	Setup() error
}

type BaseHook struct {
	Name    string
	Enabled bool
	Level   string

	formatter logrus.Formatter
	logLevels []logrus.Level
	writer    io.Writer
}

func (h *BaseHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)

	return err
}

func (h *BaseHook) Levels() []logrus.Level {
	return h.logLevels
}

func (h *BaseHook) baseSetup() {
	// setup levels
	var level = gDefaultLogLevel
	var err error
	if hookLevel := v.GetString(strings.Join([]string{"logger", h.Name, "level"}, ",")); hookLevel != "" {
		if level, err = logrus.ParseLevel(hookLevel); err != nil {
			fmt.Printf("[qlog] setup hook(%s), parse level fail:%s\n", h.Name, err)
			level = gDefaultLogLevel
		}
	}

	h.logLevels = logLevels(level)

	// setup formatters
	if hookFormatterName := v.GetString(strings.Join([]string{"logger", h.Name, "formatter", "name"}, ",")); hookFormatterName != "" {
		if h.formatter, err = newFormatter(hookFormatterName, strings.Join([]string{"logger", h.Name, "formatter", "opts"}, ",")); err != nil {
			fmt.Printf("[qlog] setup hook(%s) formatter(%s) fail:%s\n", h.Name, hookFormatterName, err)
			h.formatter = gDefaultFormatter
		}
	} else {
		h.formatter = gDefaultFormatter
	}
}

var gRegisteredHooks map[string]reflect.Type = make(map[string]reflect.Type)

func registerHook(name string, typ reflect.Type) {
	gRegisteredHooks[name] = typ
}

func newHook(name string) (logrus.Hook, error) {
	var err error
	var typ reflect.Type
	var ok bool

	if typ, ok = gRegisteredHooks[name]; !ok {
		return nil, fmt.Errorf("[qlog] hook name(%s) not registered", name)
	}

	h := reflect.New(typ)

	if err = v.UnmarshalKey(strings.Join([]string{"log", name}, "."), h.Interface()); err != nil {
		return nil, err
	}

	h.Elem().FieldByName("Name").SetString(name)

	hook := h.Interface().(logrus.Hook)

	if setup, ok := hook.(HookSetuper); ok {
		if err = setup.Setup(); err != nil {
			return nil, err
		}
	}

	return hook, nil
}

var gActivedHooks = make(logrus.LevelHooks)

func initHooks() error {
	var err error
	var hook logrus.Hook
	for name := range gRegisteredHooks {
		if v.GetBool(strings.Join([]string{"log", name, "enabled"}, ".")) == true {
			if hook, err = newHook(name); err != nil {
				fmt.Printf("[qlog] init hook(%s) error:%s\n", name, err)
				continue
			}
			gActivedHooks.Add(hook)
		}
	}

	return nil
}
