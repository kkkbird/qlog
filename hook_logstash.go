package qlog

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
)

// Hook represents a Logstash hook.
// It has two fields: writer to write the entry to Logstash and
// formatter to format the entry to a Logstash format before sending.
//
// To initialize it use the `New` function.
//
type LogstashHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	logLevels []logrus.Level
}

// New returns a new logrus.Hook for Logstash.
//
// To create a new hook that sends logs to `tcp://logstash.corp.io:9999`:
//
// conn, _ := net.Dial("tcp", "logstash.corp.io:9999")
// hook := logrustash.New(conn, logrustash.DefaultFormatter())
func NewLogstashHook(cnf LogHookConfig, rootLevel logrus.Level) (logrus.Hook, error) {
	url, err := url.Parse(cnf.Param)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial(url.Scheme, url.Host)

	if err != nil {
		return nil, err
	}

	hookLevel := rootLevel
	if len(cnf.Level) > 0 {
		hookLevel, err = logrus.ParseLevel(cnf.Level)

		if err != nil {
			return nil, fmt.Errorf("LogHook %s: log level [%s] cannot be parsed", cnf.Type, cnf.Level)
		}
	}

	logstashType := program
	if len(cnf.ExtraParam) > 0 {
		logstashType = cnf.ExtraParam
	}

	return &LogstashHook{
		writer: conn,
		formatter: LogstashDefaultFormatter(logrus.Fields{
			"type": logstashType,
		}),
		logLevels: logLevels(hookLevel),
	}, nil
}

// Fire takes, formats and sends the entry to Logstash.
// Hook's formatter is used to format the entry into Logstash format
// and Hook's writer is used to write the formatted entry to the Logstash instance.
func (h *LogstashHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)

	return err
}

// Levels returns all logrus levels.
func (h *LogstashHook) Levels() []logrus.Level {
	return h.logLevels
}

// Using a pool to re-use of old entries when formatting Logstash messages.
// It is used in the Fire function.
var entryPool = sync.Pool{
	New: func() interface{} {
		return &logrus.Entry{}
	},
}

// copyEntry copies the entry `e` to a new entry and then adds all the fields in `fields` that are missing in the new entry data.
// It uses `entryPool` to re-use allocated entries.
func copyEntry(e *logrus.Entry, fields logrus.Fields) *logrus.Entry {
	ne := entryPool.Get().(*logrus.Entry)
	ne.Message = e.Message
	ne.Level = e.Level
	ne.Time = e.Time
	ne.Data = logrus.Fields{}
	for k, v := range fields {
		ne.Data[k] = v
	}
	for k, v := range e.Data {
		ne.Data[k] = v
	}
	return ne
}

// releaseEntry puts the given entry back to `entryPool`. It must be called if copyEntry is called.
func releaseEntry(e *logrus.Entry) {
	entryPool.Put(e)
}

// LogstashFormatter represents a Logstash format.
// It has logrus.Formatter which formats the entry and logrus.Fields which
// are added to the JSON message if not given in the entry data.
//
// Note: use the `DefaultFormatter` function to set a default Logstash formatter.
type LogstashFormatter struct {
	logrus.Formatter
	logrus.Fields
}

var (
	logstashFields   = logrus.Fields{"@version": "1", "type": "log"}
	logstashFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}
)

// DefaultFormatter returns a default Logstash formatter:
// A JSON format with "@version" set to "1" (unless set differently in `fields`,
// "type" to "log" (unless set differently in `fields`),
// "@timestamp" to the log time and "message" to the log message.
//
// Note: to set a different configuration use the `LogstashFormatter` structure.
func LogstashDefaultFormatter(fields logrus.Fields) logrus.Formatter {
	for k, v := range logstashFields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}

	return &LogstashFormatter{
		Formatter: &logrus.JSONFormatter{FieldMap: logstashFieldMap},
		Fields:    fields,
	}
}

// Format formats an entry to a Logstash format according to the given Formatter and Fields.
//
// Note: the given entry is copied and not changed during the formatting process.
func (f *LogstashFormatter) Format(e *logrus.Entry) ([]byte, error) {
	ne := copyEntry(e, f.Fields)
	dataBytes, err := f.Formatter.Format(ne)
	releaseEntry(ne)
	return dataBytes, err
}
