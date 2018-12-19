package qlog

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

var (
	gRegisteredFormatters map[string]reflect.Type = make(map[string]reflect.Type)
)

func registeFormatter(name string, typ reflect.Type) {
	gRegisteredFormatters[name] = typ
}

func newFormatter(name string, key string) (logrus.Formatter, error) {
	var err error
	var typ reflect.Type
	var ok bool

	if typ, ok = gRegisteredFormatters[name]; !ok {
		return nil, errors.New(fmt.Sprintf("[qlog] formatter name(%s) not registered", name))
	}

	f := reflect.New(typ)

	if err = v.UnmarshalKey(key, f.Interface()); err != nil {
		return nil, err
	}

	return f.Interface().(logrus.Formatter), nil
}

var gDefaultFormatter logrus.Formatter

func initFormatters() error {
	var err error
	defaultFormatterName := v.GetString("logger.formatter.name")

	if defaultFormatterName == "" {
		gDefaultFormatter = &logrus.TextFormatter{}
	} else {
		if gDefaultFormatter, err = newFormatter(defaultFormatterName, "logger.formatter.opts"); err != nil {
			return err
		}
	}

	return nil
}
