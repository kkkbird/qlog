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
	var level = gDefaultLogLevel
	var err error
	if h.Level = v.GetString(strings.Join([]string{"logger", h.Name, "level"}, ".")); h.Level != "" {
		if level, err = logrus.ParseLevel(h.Level); err != nil {
			fmt.Printf("[qlog] setup hook(%s), parse level fail:%s\n", h.Name, err)
			level = gDefaultLogLevel
		}
	}

	h.logLevels = logLevels(level)

	// setup formatters
	if hookFormatterName := v.GetString(strings.Join([]string{"logger", h.Name, "formatter", "name"}, ".")); hookFormatterName != "" {
		if h.formatter, err = newFormatter(hookFormatterName, strings.Join([]string{"logger", h.Name, "formatter", "opts"}, ".")); err != nil {
			fmt.Printf("[qlog] setup hook(%s) formatter(%s) fail:%s\n", h.Name, hookFormatterName, err)
			h.formatter = gDefaultFormatter
		}
	} else {
		h.formatter = gDefaultFormatter
	}
}

var gRegisteredHooks = make(map[string]reflect.Type)

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

	h.Elem().FieldByName("Name").SetString(name)
	h.Elem().FieldByName("Enabled").SetBool(true)

	hook := h.Interface().(logrus.Hook)

	if setup, ok := hook.(HookSetuper); ok {
		if err = setup.Setup(); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("[qlog] hook(%s) has no Setup() func", name)
	}

	return hook, nil
}

var gActivedHooks = make(logrus.LevelHooks)

func hasActivedHook() bool {
	return true
}

func initHooks() error {
	var err error
	var hook logrus.Hook

	for name := range gRegisteredHooks {
		n := strings.Join([]string{"logger", name, "enabled"}, ".")
		//fmt.Println("initHooks", n, v.GetBool(n))
		if v.GetBool(n) == true {
			if hook, err = newHook(name); err != nil {
				fmt.Printf("[qlog] init hook(%s) error:%s\n", name, err)
				continue
			}
			gActivedHooks.Add(hook)
		}
	}

	return nil
}

func resetHooks() {
	gActivedHooks = make(logrus.LevelHooks)
}
