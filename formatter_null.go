package qlog

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

type NullFormatter struct {
}

func (NullFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
}

var _InitNullFormatter = func() interface{} {
	registeFormatter("null", reflect.TypeOf(NullFormatter{}))
	return nil
}()
