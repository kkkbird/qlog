package qlog

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

func init() {
	registeFormatter("json", reflect.TypeOf(logrus.JSONFormatter{}))
}
