package qlog

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

var _InitJSONFormat = func() interface{} {
	registeFormatter("json", reflect.TypeOf(logrus.JSONFormatter{}))

	return nil
}()
