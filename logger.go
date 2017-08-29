package qlog

import (
	"flag"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	"fmt"
	"net/url"

	"time"

	"github.com/kkkbird/logrus-logstash-hook"
	colorable "github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var (
	// qlogger is the name of the standard logger in qloggerlib `log`
	qlogger QLogger
)

const (
	longTimeStamp  = "2006/01/02 15:04:05.000000Z07:00"
	shortTimeStamp = "06/01/02 15:04:05.000"
)

type QLogger struct {
	logger            *logrus.Logger
	loglevel          string
	logfmt            string
	logdir            string
	logstash          string
	logstashtype      string
	logoptions        string
	withRuntimeFields bool
	initOnce          sync.Once
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

	l.withRuntimeFields = true
	timestampFormat := ""
	disableTimeStamp := false

	options := strings.Split(strings.ToLower(l.logoptions), ",")

	for _, opt := range options {
		if len(opt) == 0 {
			continue
		}
		switch strings.ToLower(opt) {
		case "shorttime":
			timestampFormat = shortTimeStamp
		case "longtime":
			timestampFormat = longTimeStamp
		case "notime":
			disableTimeStamp = true
		case "disableruntime":
			l.withRuntimeFields = false
		default:
			return fmt.Errorf("not a valid logoption:%s", opt)
		}
	}

	var formatter logrus.Formatter
	switch l.logfmt {
	case "classic":
		formatter = &ClassicFormatter{
			TimestampFormat:  timestampFormat,
			DisableTimestamp: disableTimeStamp,
		}
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat:  timestampFormat,
			DisableTimestamp: disableTimeStamp,
		}
	case "kv":
		fallthrough
	default:
		formatter = &logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  timestampFormat,
			DisableTimestamp: disableTimeStamp,
		}
	}

	var out io.Writer

	if len(l.logdir) == 0 {
		if _formatter, ok := formatter.(*logrus.TextFormatter); ok {
			_formatter.ForceColors = true
		}
		out = colorable.NewColorableStdout()
	} else {
		file, _, err := create(time.Now())
		if err != nil {
			return err
		}
		out = file
	}
	l.logger = &logrus.Logger{
		//Out: os.Stderr,
		Out:       out,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     loglevel,
	}

	if len(l.logstash) > 0 {
		logstashUrl, err := url.Parse(l.logstash)
		if err != nil {
			return err
		}

		conn, err := net.Dial(logstashUrl.Scheme, logstashUrl.Host)

		if err != nil {
			return err
		}

		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{
			"type": l.logstashtype,
		}))

		if err != nil {
			return err
		}

		l.logger.Hooks.Add(hook)
	}

	return nil
}

func runtimeFields(skip int) logrus.Fields {
	_, file, line, ok := runtime.Caller(skip)
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

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return &Entry{
		e:                 qlogger.Logger().WithError(err),
		withRunTimeFields: qlogger.withRuntimeFields,
	}
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return &Entry{
		e:                 qlogger.Logger().WithField(key, value),
		withRunTimeFields: qlogger.withRuntimeFields,
	}
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *Entry {
	return &Entry{
		e:                 qlogger.Logger().WithFields(fields),
		withRunTimeFields: qlogger.withRuntimeFields,
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Debug(args...)
	} else {
		qlogger.Logger().Debug(args...)
	}
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Info(args...)
	} else {
		qlogger.Logger().Info(args...)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Info(args...)
	} else {
		qlogger.Logger().Info(args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warn(args...)
	} else {
		qlogger.Logger().Warn(args...)
	}
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warn(args...)
	} else {
		qlogger.Logger().Warn(args...)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Error(args...)
	} else {
		qlogger.Logger().Error(args...)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Panic(args...)
	} else {
		qlogger.Logger().Panic(args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Fatal(args...)
	} else {
		qlogger.Logger().Fatal(args...)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Debugf(format, args...)
	} else {
		qlogger.Logger().Debugf(format, args...)
	}
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		qlogger.Logger().Infof(format, args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		qlogger.Logger().Infof(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		qlogger.Logger().Warnf(format, args...)
	}
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		qlogger.Logger().Warnf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Errorf(format, args...)
	} else {
		qlogger.Logger().Errorf(format, args...)
	}
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Panicf(format, args...)
	} else {
		qlogger.Logger().Panicf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Fatalf(format, args...)
	} else {
		qlogger.Logger().Fatalf(format, args...)
	}
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Debugln(args...)
	} else {
		qlogger.Logger().Debugln(args...)
	}
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		qlogger.Logger().Infoln(args...)
	}
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		qlogger.Logger().Infoln(args...)
	}
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		qlogger.Logger().Warnln(args...)
	}
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		qlogger.Logger().Warnln(args...)
	}
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Errorln(args...)
	} else {
		qlogger.Logger().Errorln(args...)
	}
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Panicln(args...)
	} else {
		qlogger.Logger().Panicln(args...)
	}
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	if qlogger.withRuntimeFields {
		qlogger.Logger().WithFields(runtimeFields(2)).Fatalln(args...)
	} else {
		qlogger.Logger().Fatalln(args...)
	}
}

func init() {
	flag.StringVar(&qlogger.logfmt, "logfmt", "kv", "logfmt:kv,json,classic")
	flag.StringVar(&qlogger.loglevel, "loglevel", "info", "log level:debug,info,waring,fatal,panic")
	flag.StringVar(&qlogger.logoptions, "logoptions", "longtime", "log options, longtime|shorttime|notime, disableruntime")
	flag.StringVar(&qlogger.logdir, "logdir", "", "log dir, leave empty to log to stderr")
	flag.StringVar(&qlogger.logstash, "logstash", "", "logstash address, also log to logstash, example: udp://192.168.0.92:5000")
	flag.StringVar(&qlogger.logstashtype, "logstashtype", program, "logstash type field, only available when logstash mode")

	//flag.BoolVar(&qlogger.withRuntimeFields, "logruntime", true, "log with runtime fields")
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
