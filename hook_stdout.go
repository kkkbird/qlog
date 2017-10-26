package qlog

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type StdoutHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	logLevels []logrus.Level
}

func NewStdoutHook(cnf LogHookConfig, logfmt logrus.Formatter, rootLevel logrus.Level) (logrus.Hook, error) {
	//only use std hook when root out discarded
	if strings.ToLower(fLogOutput) != "discard" {
		return nil, errors.New("set logoutput=discard if want to use std_hook")
	}

	out := createQLogStdWriter(cnf.Type)

	var err error

	hookLevel := rootLevel
	if len(cnf.Level) > 0 {
		hookLevel, err = logrus.ParseLevel(cnf.Level)

		if err != nil {
			return nil, fmt.Errorf("LogHook %s: log level [%s] cannot be parsed", cnf.Type, cnf.Level)
		}
	}

	return &StdoutHook{
		writer:    out,
		formatter: logfmt,
		logLevels: logLevels(hookLevel),
	}, nil
}

func (h *StdoutHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)

	return err
}

func (h *StdoutHook) Levels() []logrus.Level {
	return h.logLevels
}
