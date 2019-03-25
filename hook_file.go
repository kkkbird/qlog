package qlog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// FileHook implement file support of logger hook
type FileHook struct {
	BaseHook

	FilePath string
	FileName string

	// Rotate params
	RotateTime   time.Duration // 0 means do not rotate
	RotateMaxAge time.Duration // time to wait until old logs are purged, default 7 days, set 0 to disable
	RotateCount  uint          // the number of files should be kept, default 0 means disabled
}

const (
	keyFileEnabled      = "logger.file.enabled"
	keyFileLevel        = "logger.file.level"
	keyFilePath         = "logger.file.path"
	keyFileName         = "logger.file.name"
	keyFileRotateTime   = "logger.file.rotate.time"
	keyFileRotateMaxAge = "logger.file.rotate.maxage"
	keyFileRotateCount  = "logger.file.rotate.count"
)

// Setup function for FileHook
func (h *FileHook) Setup() error {
	var err error
	var fullPath string

	h.baseSetup()

	h.FilePath = v.GetString(keyFilePath)
	h.FileName = v.GetString(keyFileName)

	rotateTime := v.GetString(keyFileRotateTime)

	if _, err = os.Stat(h.FilePath); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}

	if fullPath, err = filepath.Abs(filepath.Join(h.FilePath, h.FileName)); err != nil {
		return err
	}

	if h.RotateTime, err = time.ParseDuration(rotateTime); err != nil {
		return fmt.Errorf("Parse logger.file.rotate.time fail: %s", err)
	}

	if h.RotateTime > 0 {
		if h.RotateMaxAge, err = time.ParseDuration(v.GetString(keyFileRotateMaxAge)); err != nil {
			return fmt.Errorf("Parse logger.file.rotate.maxage fail: %s", err)
		}

		h.RotateCount = uint(v.GetInt(keyFileRotateCount))

		if h.writer, err = rotatelogs.New(fullPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(fullPath),
			rotatelogs.WithMaxAge(h.RotateMaxAge),
			rotatelogs.WithRotationTime(h.RotateTime),
			rotatelogs.WithRotationCount(h.RotateCount),
		); err != nil {
			return fmt.Errorf("Create rotate log fail: %s", err)
		}
	} else {
		if h.writer, err = os.Create(fullPath); err != nil {
			return fmt.Errorf("Create log fail: %s", err)
		}
	}

	return nil
}

var _InitFileHook = func() interface{} {
	cli.Bool(keyFileEnabled, false, "logger.file.enabled")
	cli.String(keyFileLevel, "", "logger.file.level") // DONOT set default level in pflag

	cli.String(keyFilePath, ".", "logger.file.path")
	cli.String(keyFileName, "qlog.log", "logger.file.name")
	cli.String(keyFileRotateTime, "1d", "logger.file.rotate.time")
	cli.String(keyFileRotateMaxAge, "7d", "logger.file.rotate.maxag")
	cli.String(keyFileRotateCount, "0", "logger.file.rotate.count")

	registerHook("file", reflect.TypeOf(FileHook{}))

	return nil
}()
