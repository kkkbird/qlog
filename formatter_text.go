package qlog

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

func init() {
	registeFormatter("text", reflect.TypeOf(logrus.TextFormatter{}))
}
