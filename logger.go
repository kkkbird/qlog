package qlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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

var (
	// default logger object
	DefaultLogger *Logger

	// default logger settings
	gDefaultLogLevel logrus.Level
	gReportCaller    bool

	// system param
	gPid      = os.Getpid()
	gProgram  = filepath.Base(os.Args[0])
	gHost     = "unknownhost"
	gUserName = "unknownuser"
)

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

func initSysParams() error {
	h, err := os.Hostname()
	if err == nil {
		gHost = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		gUserName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	gUserName = strings.Replace(gUserName, `\`, "_", -1)

	return nil
}

func initViper() error {
	v = viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/qlog")
	v.AddConfigPath("./conf/")
	v.SetConfigName("logger")
	v.SetConfigType("yaml")

	v.SetDefault("logger.level", "debug")
	v.SetDefault("logger.reportcaller", false)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func init() {
	var err error

	if err = initSysParams(); err != nil {
		panic(fmt.Sprint("[qlog] init system param error:", err))
	}

	if err = initViper(); err != nil {
		panic(fmt.Sprint("[qlog] init viper error:", err))
	}

	if gDefaultLogLevel, err = logrus.ParseLevel(v.GetString("logger.level")); err != nil {
		panic(fmt.Sprint("[qlog] get default log level error:", err))
	}

	gReportCaller = v.GetBool("logger.reportcaller")

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
				ReportCaller: gReportCaller,
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
				ReportCaller: gReportCaller,
			},
		}
	}
}
