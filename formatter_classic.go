package qlog

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"

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

type ShortLevel uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level ShortLevel) String() string {
	_level := logrus.Level(level)
	switch _level {
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

	TruncateCallerPath bool
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
		if f.TruncateCallerPath {
			f.appendValueOnly(b, fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line))
		} else {
			f.appendValueOnly(b, fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line))
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
