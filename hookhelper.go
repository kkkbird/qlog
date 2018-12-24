package qlog

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

func getLogLevels(baseLevel logrus.Level) (level []logrus.Level) {
	level = make([]logrus.Level, 0)
	for i := baseLevel; i > logrus.PanicLevel; i-- {
		level = append(level, i)
	}
	return
}

type HookSetuper interface {
	Setup() error
}

type BaseHook struct {
	Name  string
	Level string

	formatter logrus.Formatter
	logLevels []logrus.Level
	writer    io.Writer
}

func (h *BaseHook) Fire(e *logrus.Entry) error {
	// fmt.Println("fire:", h.Name)
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
	var level = qLogger.Level
	var err error
	if h.Level = v.GetString(strings.Join([]string{"logger", h.Name, "level"}, ".")); h.Level != "" {
		if level, err = logrus.ParseLevel(h.Level); err != nil {
			fmt.Printf("[qlog] setup hook(%s), parse level fail:%s\n", h.Name, err)
			level = qLogger.Level
		}
	}

	h.logLevels = getLogLevels(level)

	// setup formatters
	if hookFormatterName := v.GetString(strings.Join([]string{"logger", h.Name, "formatter", "name"}, ".")); hookFormatterName != "" {
		if h.formatter, err = newFormatter(hookFormatterName, strings.Join([]string{"logger", h.Name, "formatter", "opts"}, ".")); err != nil {
			fmt.Printf("[qlog] setup hook(%s) formatter(%s) fail:%s\n", h.Name, hookFormatterName, err)
			h.formatter = qLogger.Formatter
		}
	} else {
		h.formatter = qLogger.Formatter
	}
}

var gRegisteredHooks map[string]reflect.Type

func registerHook(name string, typ reflect.Type) {
	if gRegisteredHooks == nil {
		gRegisteredHooks = make(map[string]reflect.Type)
	}

	gRegisteredHooks[name] = typ

	if _, ok := reflect.New(typ).Interface().(HookSetuper); !ok {
		panic(fmt.Sprintf("[qlog] registe hook (%s) fail: must be HookSetuper()", name))
	}
}

func newHook(name string) (logrus.Hook, error) {
	var err error
	var typ reflect.Type
	var ok bool

	if typ, ok = gRegisteredHooks[name]; !ok {
		return nil, fmt.Errorf("[qlog] hook name(%s) not registered", name)
	}

	h := reflect.New(typ)

	h.Elem().FieldByName("Name").SetString(name)

	hook := h.Interface().(logrus.Hook)

	setuper, _ := hook.(HookSetuper)
	if err = setuper.Setup(); err != nil {
		return nil, err
	}

	return hook, nil
}
