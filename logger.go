package qlog

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logger struct {
	*logrus.Logger
}

var v *viper.Viper

func logLevels(baseLevel logrus.Level) (level []logrus.Level) {
	level = make([]logrus.Level, 0)
	for i := baseLevel; i > logrus.PanicLevel; i-- {
		level = append(level, i)
	}
	return
}

var DefaultLogger *Logger

var gDefaultLogLevel logrus.Level

func initViper() error {
	v = viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/qlog")
	v.AddConfigPath("./conf/")
	v.SetConfigName("logger")
	v.SetConfigType("yaml")

	v.SetDefault("logger.level", "debug")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func init() {
	var err error
	if err = initViper(); err != nil {
		panic(fmt.Sprint("[qlog] init viper error:", err))
	}

	if gDefaultLogLevel, err = logrus.ParseLevel(v.GetString("logger.level")); err != nil {
		panic(fmt.Sprint("[qlog] get default log level error:", err))
	}

	if err = initFormatters(); err != nil {
		panic(fmt.Sprint("[qlog] init default formatter:", err))
	}

	if err = initHooks(); err != nil {
		panic(fmt.Sprint("[qlog] init default formatter:", err))
	}

	if len(gActivedHooks) > 0 {
		DefaultLogger = &Logger{
			Logger: &logrus.Logger{
				Out:          ioutil.Discard,
				Formatter:    gDefaultFormatter,
				Hooks:        gActivedHooks,
				Level:        gDefaultLogLevel,
				ExitFunc:     os.Exit,
				ReportCaller: true,
			},
		}
	} else {
		DefaultLogger = &Logger{
			Logger: &logrus.Logger{
				Out:          os.Stderr,
				Formatter:    gDefaultFormatter,
				Hooks:        nil,
				Level:        gDefaultLogLevel,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			},
		}
	}
}
