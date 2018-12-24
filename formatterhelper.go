package qlog

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

const (
	longTimeStamp  = "2006/01/02 15:04:05.000000Z07:00"
	shortTimeStamp = "06/01/02 15:04:05.000"
)

var (
	gRegisteredFormatters map[string]reflect.Type
)

func registeFormatter(name string, typ reflect.Type) {
	if gRegisteredFormatters == nil {
		gRegisteredFormatters = make(map[string]reflect.Type)
	}

	gRegisteredFormatters[name] = typ
}

func newFormatter(name string, key string) (logrus.Formatter, error) {
	var err error
	var typ reflect.Type
	var ok bool

	if typ, ok = gRegisteredFormatters[name]; !ok {
		return nil, fmt.Errorf("[qlog] formatter name(%s) not registered", name)
	}

	f := reflect.New(typ)

	if err = v.UnmarshalKey(key, f.Interface()); err != nil {
		return nil, err
	}

	return f.Interface().(logrus.Formatter), nil
}
