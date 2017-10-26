package qlog

// from glog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type FileHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	logLevels []logrus.Level
}

func (h *FileHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)

	return err
}

// Levels returns all logrus levels.
func (h *FileHook) Levels() []logrus.Level {
	return h.logLevels
}

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1800

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d.%s.%s.%d.log",
		program,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		host,
		userName,
		pid,
	)
	return name, fmt.Sprintf("%s.log", program)
}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(cnf LogHookConfig, t time.Time) (f *os.File, filename string, err error) {
	// logDirs lists the candidate directories for new log files.
	var logDirs []string

	if len(cnf.Param) > 0 {
		logDirs = append(logDirs, cnf.Param)
	}
	logDirs = append(logDirs, os.TempDir())

	if len(logDirs) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}

	name, link := logName(t)
	var lastErr error
	for _, dir := range logDirs {
		fname := filepath.Join(dir, name)
		f, err := os.Create(fname)
		if err == nil {
			symlink := filepath.Join(dir, link)
			os.Remove(symlink)        // ignore err
			os.Symlink(name, symlink) // ignore err
			return f, fname, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}

func NewFileHook(cnf LogHookConfig, rootLevel logrus.Level) (logrus.Hook, error) {
	file, _, err := create(cnf, time.Now())
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

	return &FileHook{
		writer:    file,
		formatter: createQLogFormatter(cnf.Fmt),
		logLevels: logLevels(hookLevel),
	}, nil
}
