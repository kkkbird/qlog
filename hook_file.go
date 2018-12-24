package qlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"
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
	FileName   string
	fileWriter *os.File

	// Rotate at line
	MaxLines         int
	maxLinesCurLines int

	// Rotate at size
	MaxSize        int
	maxSizeCurSize int

	// Rotate daily
	Daily         bool
	MaxDays       int
	dailyOpenDate int
	dailyOpenTime time.Time

	Rotate bool

	Perm       string
	RotatePerm string
}

const (
	keyFileEnabled    = "logger.file.enabled"
	keyFileLevel      = "logger.file.level"
	keyFileDir        = "logger.file.filedir"
	keyFileName       = "logger.file.filename"
	keyFileMaxLines   = "logger.file.maxlines"
	keyFileDaily      = "logger.file.daily"
	keyFileMaxDays    = "logger.file.maxdays"
	keyFileRotate     = "logger.file.rotate"
	keyFilePerm       = "logger.file.perm"
	keyFileRotatePerm = "logger.file.rotateperm"
)

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1800

func (h *FileHook) Setup() error {
	var err error

	h.baseSetup()

	h.FileDir = v.GetString(keyFileDir)
	h.FileName = v.GetString(keyFileName)
	h.MaxLines = v.GetInt(keyFileMaxLines)
	h.Daily = v.GetBool(keyFileDaily)
	h.MaxDays = v.GetInt(keyFileMaxDays)
	h.Rotate = v.GetBool(keyFileRotate)
	h.Perm = v.GetString(keyFilePerm)
	h.RotatePerm = v.GetString(keyFileRotatePerm)

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
		return nil, "", fmt.Errorf("logDirs <%q> not exist", logDirs)
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

var _InitFileHook = func() interface{} {
	cli.Bool(keyFileEnabled, false, "logger.file.enabled")
	cli.String(keyFileLevel, "", "logger.file.level") // DONOT set default level in pflag

	cli.String(keyFileDir, "", "logger.file.filedir")
	cli.String(keyFileName, "", "logger.file.filename")
	cli.Int(keyFileMaxLines, 0, "logger.file.maxlines")
	cli.Bool(keyFileDaily, false, "logger.file.daily")
	cli.Int(keyFileMaxDays, 0, "logger.file.maxdays")
	cli.Bool(keyFileRotate, false, "logger.file.rotate")
	cli.String(keyFilePerm, "0440", "logger.file.perm")
	cli.String(keyFileRotatePerm, "0660", "logger.file.rotateperm")

	registerHook("file", reflect.TypeOf(FileHook{}))

	return nil
}()
