package qlog

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

var _initTextFormatter = func() interface{} {
	registeFormatter("text", reflect.TypeOf(logrus.TextFormatter{}))

	return nil
}()
