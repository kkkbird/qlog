package qlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

var (
	ErrNoLogDir = errors.New("no log dir exists")
)

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d.%s.%s.PID%d.log",
		gProgram,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		gHost,
		gUserName,
		gPid,
	)
	return name, fmt.Sprintf("%s.log", gProgram)
}

type FileHook struct {
	BaseHook

	FileDir    string
	Filename   string
	fileWriter *os.File

	// Rotate at line
	MaxLines         int
	maxLinesCurLines int

	// Rotate at size
	MaxSize        int
	maxSizeCurSize int

	// Rotate daily
	Daily         bool
	MaxDays       int64
	dailyOpenDate int
	dailyOpenTime time.Time

	Rotate bool

	Perm       string
	RotatePerm string
}

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1800

func (h *FileHook) Setup() error {
	var err error

	if err = v.UnmarshalKey("logger.file", h); err != nil {
		return err
	}

	h.baseSetup()

	var f io.Writer
	if f, _, err = h.create(time.Now()); err != nil {
		return err
	}

	h.writer = f

	return nil
}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func (h *FileHook) create(t time.Time) (f *os.File, filename string, err error) {
	// logDirs lists the candidate directories for new log files.
	var logDirs []string

	if len(h.FileDir) > 0 {
		logDirs = append(logDirs, h.FileDir)
	}
	//logDirs = append(logDirs, os.TempDir())

	if len(logDirs) == 0 {
		return nil, "", ErrNoLogDir
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
	return nil, "", fmt.Errorf("cannot create log: %v", lastErr)
}

func init() {
	registerHook("file", reflect.TypeOf(FileHook{}))
}
