package qlog

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

//copy from old logrus code
func prefixFieldClashes(data logrus.Fields, hasCaller bool) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}

	if hasCaller {
		if fc, ok := data["func"]; ok {
			data["fields.func"] = fc
		}

		if f, ok := data["file"]; ok {
			data["fields.file"] = f
		}
	}
}

// ShortLevel is a simple char to indicate log level
type ShortLevel uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level ShortLevel) String() string {
	_level := logrus.Level(level)
	switch _level {
	case logrus.TraceLevel:
		return "T"
	case logrus.DebugLevel:
		return "D"
	case logrus.InfoLevel:
		return "I"
	case logrus.WarnLevel:
		return "W"
	case logrus.ErrorLevel:
		return "E"
	case logrus.FatalLevel:
		return "F"
	case logrus.PanicLevel:
		return "P"
	}

	return "X"
}

// ClassicFormatter formats logs into parsable json
type ClassicFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
}

// Format renders a single log entry
func (f *ClassicFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	prefixFieldClashes(entry.Data, entry.HasCaller())

	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = longTimeStamp
	}

	if !f.DisableTimestamp {
		b.WriteString(entry.Time.Format(timestampFormat))
	}

	//reportcaller is enabled
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}

		if len(fileVal) > 0 {
			f.appendValueOnly(b, fileVal)
		}

		if len(funcVal) > 0 {
			f.appendValueOnly(b, funcVal)
		}
	}

	f.appendValueOnly(b, fmt.Sprintf("[%s]", ShortLevel(entry.Level).String()))

	if len(entry.Message) > 0 {
		f.appendValueOnly(b, entry.Message)
	}

	for k, v := range entry.Data {
		f.appendKeyValue(b, k, v)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *ClassicFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *ClassicFormatter) appendValueOnly(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	f.appendValue(b, value)
}

func (f *ClassicFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	b.WriteString(fmt.Sprintf("%s", stringVal))
}

var _InitClassicFormatter = func() interface{} {
	registeFormatter("classic", reflect.TypeOf(ClassicFormatter{}))
	return nil
}()
