package qlog

import (
	"flag"
	"os"
	"runtime"
	"strings"
	"sync"

	"fmt"
	"net/url"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
)

var (
	// qlogger is the name of the standard logger in qloggerlib `log`
	qlogger QLogger
)

type QLogger struct {
	logger       *logrus.Logger
	loglevel     string
	logtype      string
	logdir       string
	logstash     string
	logstashtype string

	initOnce sync.Once
}

func (l *QLogger) Logger() *logrus.Logger {
	l.initOnce.Do(func() {
		err := l.prepare()
		if err != nil {
			panic(err)
		}

		l.logger.SetNoLock()
	})
	return l.logger
}

func (l *QLogger) prepare() (err error) {
	if !flag.Parsed() {
		return fmt.Errorf("flag not Parsed, call flag.Parse() first")
	}
	var loglevel logrus.Level

	if loglevel, err = logrus.ParseLevel(l.loglevel); err != nil {
		return err
	}

	if len(l.logdir) == 0 {
		l.logger = &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			},
			Hooks: make(logrus.LevelHooks),
			Level: loglevel,
		}
	} else {
		file, _, err := create(time.Now())
		if err != nil {
			return err
		}
		l.logger = &logrus.Logger{
			Out: file,
			Formatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			},
			Hooks: make(logrus.LevelHooks),
			Level: loglevel,
		}
	}

	if len(l.logstash) > 0 {
		logstashUrl, err := url.Parse(l.logstash)
		if err != nil {
			return err
		}
		hook, err := logrus_logstash.NewHookWithFieldsAndPrefix(logstashUrl.Scheme, logstashUrl.Host, l.logstashtype, logrus.Fields{
			"PID":       pid,
			"_hostname": host,
			"_user":     userName,
		}, "_")

		if err != nil {
			return err
		}

		l.logger.Hooks.Add(hook)
	}

	return nil
}

func runtimeFields() logrus.Fields {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return logrus.Fields{
		"FILE": file,
		"LINE": line,
	}
}

func withRuntimeFields() *logrus.Entry {
	return qlogger.Logger().WithFields(runtimeFields())
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return WithField(logrus.ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	defaultFields := runtimeFields()
	defaultFields[key] = value
	return qlogger.Logger().WithFields(defaultFields)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	defaultFields := runtimeFields()
	for k, v := range fields {
		defaultFields[k] = v
	}
	return qlogger.Logger().WithFields(defaultFields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	withRuntimeFields().Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	withRuntimeFields().Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	withRuntimeFields().Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	withRuntimeFields().Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	withRuntimeFields().Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	withRuntimeFields().Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	withRuntimeFields().Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	withRuntimeFields().Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	withRuntimeFields().Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	withRuntimeFields().Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	withRuntimeFields().Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	withRuntimeFields().Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	withRuntimeFields().Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	withRuntimeFields().Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	withRuntimeFields().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	withRuntimeFields().Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	withRuntimeFields().Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	withRuntimeFields().Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	withRuntimeFields().Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	withRuntimeFields().Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	withRuntimeFields().Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	withRuntimeFields().Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	withRuntimeFields().Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	withRuntimeFields().Fatalln(args...)
}

func init() {
	flag.StringVar(&qlogger.loglevel, "loglevel", "info", "log level:debug,info,waring,fatal,panic")
	flag.StringVar(&qlogger.logdir, "logdir", "", "log dir, leave empty to log to stderr")
	flag.StringVar(&qlogger.logstash, "logstash", "", "logstash address, also log to logstash, example: udp://192.168.0.92:5000")
	flag.StringVar(&qlogger.logstashtype, "logstashtype", program, "logstash type field, only available when logstash mode")
}

//get logrus logger
func Logger() *logrus.Logger {
	return qlogger.Logger()
}

func Flush() {
	if qlogger.logger == nil {
		return
	}

	if file, isOk := qlogger.logger.Out.(*os.File); isOk {
		file.Close()
	}
}
