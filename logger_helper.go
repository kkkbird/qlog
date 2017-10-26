package qlog

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	colorable "github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
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

func logLevels(baseLevel logrus.Level) (level []logrus.Level) {
	level = make([]logrus.Level, 0)
	for i := baseLevel; i > logrus.PanicLevel; i-- {
		level = append(level, i)
	}
	return
}

func createQLogFormatter(fmtType string) (formatter logrus.Formatter) {
	switch strings.ToLower(fmtType) {
	case "classic":
		formatter = &ClassicFormatter{
			TimestampFormat:  fTimestampFormat,
			DisableTimestamp: fDisableTimeStamp,
		}
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat:  fTimestampFormat,
			DisableTimestamp: fDisableTimeStamp,
		}
	case "kv":
		fallthrough
	default:
		formatter = &logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  fTimestampFormat,
			DisableTimestamp: fDisableTimeStamp,
			ForceColors:      fForceColors,
			DisableColors:    fDisableColors,
		}
	}
	return
}

func createQLogStdWriter(writerType string) (writer io.Writer) {
	switch strings.ToLower(writerType) {
	case "stderr":
		writer = colorable.NewColorableStderr()
	case "discard":
		writer = ioutil.Discard
	case "stdout":
		fallthrough
	default:
		writer = colorable.NewColorableStdout()
	}
	return
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	e := logger().WithError(err) //call Logger() first to make sure prepare could be called first
	return &Entry{
		e:          e,
		logruntime: fLogRuntime,
	}
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	e := logger().WithField(key, value) //call Logger() first to make sure prepare could be called first
	return &Entry{
		e:          e,
		logruntime: fLogRuntime,
	}
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *Entry {
	e := logger().WithFields(fields) //call Logger() first to make sure prepare could be called first
	return &Entry{
		e:          e,
		logruntime: fLogRuntime,
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first

	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Debug(args...)
	} else {
		l.Debug(args...)
	}
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Info(args...)
	} else {
		l.Info(args...)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Info(args...)
	} else {
		l.Info(args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warn(args...)
	} else {
		l.Warn(args...)
	}
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warn(args...)
	} else {
		l.Warn(args...)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Error(args...)
	} else {
		l.Error(args...)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Panic(args...)
	} else {
		l.Panic(args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Fatal(args...)
	} else {
		l.Fatal(args...)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Debugf(format, args...)
	} else {
		l.Debugf(format, args...)
	}
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		l.Infof(format, args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		l.Infof(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		l.Warnf(format, args...)
	}
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		l.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Errorf(format, args...)
	} else {
		l.Errorf(format, args...)
	}
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Panicf(format, args...)
	} else {
		l.Panicf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Fatalf(format, args...)
	} else {
		l.Fatalf(format, args...)
	}
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Debugln(args...)
	} else {
		l.Debugln(args...)
	}
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		l.Infoln(args...)
	}
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		l.Infoln(args...)
	}
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		l.Warnln(args...)
	}
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		l.Warnln(args...)
	}
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Errorln(args...)
	} else {
		l.Errorln(args...)
	}
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Panicln(args...)
	} else {
		l.Panicln(args...)
	}
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	l := logger() //call Logger() first to make sure prepare could be called first
	if fLogRuntime {
		l.WithFields(runtimeFields(2)).Fatalln(args...)
	} else {
		l.Fatalln(args...)
	}
}
