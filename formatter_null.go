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

func init() {
	registeFormatter("null", reflect.TypeOf(NullFormatter{}))
}
