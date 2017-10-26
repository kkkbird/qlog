package qlog

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
	"sync"

	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	longTimeStamp  = "2006/01/02 15:04:05.000000Z07:00"
	shortTimeStamp = "06/01/02 15:04:05.000"
)

type LogHookConfig struct {
	Type       string `json:"type"`
	Level      string `json:"level,omitempty"`
	Fmt        string `json:"fmt,omitempty"`
	Options    bool   `json:"options,omitempty"`
	Param      string `json:"param,omitempty"`
	ExtraParam string `json:"extra_param,omitempty"`
}

var (
	fLogLevel         string
	fLogFmt           string
	fLogOutput        string
	fLogHookConfigs   string
	fLogOptions       string
	fLogRuntime       bool           //from fLogOptions
	fTimestampFormat  string         //from fLogOptions
	fDisableTimeStamp bool           //from fLogOptions
	fForceColors      bool           //from fLogOptions
	fDisableColors    bool           //from fLogOptions
	_logger           *logrus.Logger // do not call _logger directed, call Logger to prepare
	initOnce          sync.Once
	initMutex         sync.Mutex
	gFlushed          bool
)

func logger() *logrus.Logger {
	initOnce.Do(func() {
		err := prepare()
		if err != nil {
			panic(err)
		}
		_logger.SetNoLock()
	})
	return _logger
}

func flush() {
	initMutex.Lock()
	defer initMutex.Unlock()

	gFlushed = true

	if _logger == nil {
		return
	}

	//TODO, close hook if file
	if file, isOk := _logger.Out.(*os.File); isOk {
		file.Close()
	}
}

func parseOptions(options string) error {
	opts := strings.Split(strings.ToLower(options), ",")

	for _, opt := range opts {
		if len(opt) == 0 {
			continue
		}
		switch strings.ToLower(opt) {
		case "shorttime":
			fTimestampFormat = shortTimeStamp
		case "longtime":
			fTimestampFormat = longTimeStamp
		case "notime":
			fDisableTimeStamp = true
		case "enableruntime":
			fLogRuntime = true
		case "forcecolors":
			fForceColors = true
		case "disablecolors":
			fDisableColors = true
		default:
			return fmt.Errorf("not a valid logoption:%s", opt)
		}
	}
	return nil
}

func prepare() (err error) {
	if !flag.Parsed() {
		return fmt.Errorf("flag not Parsed, call flag.Parse() first")
	}

	initMutex.Lock()
	defer initMutex.Unlock()

	if gFlushed {
		return errors.New("already flushed ")
	}
	var lvl logrus.Level

	if lvl, err = logrus.ParseLevel(fLogLevel); err != nil {
		return err
	}

	if err = parseOptions(fLogOptions); err != nil {
		return err
	}

	_logger = &logrus.Logger{
		Out:       createQLogStdWriter(fLogOutput),
		Formatter: createQLogFormatter(fLogFmt),
		Hooks:     make(logrus.LevelHooks),
		Level:     lvl,
	}

	//init hooks
	if len(fLogHookConfigs) > 0 {
		var configs []LogHookConfig
		err = json.Unmarshal([]byte(fLogHookConfigs), &configs)
		if err != nil {
			return err
		}

		var hook logrus.Hook
		for _, cnf := range configs {
			switch strings.ToLower(cnf.Type) {
			case "stderr":
				fallthrough
			case "stdout":
				hook, err = NewStdoutHook(cnf, _logger.Formatter, lvl)
			case "file":
				hook, err = NewFileHook(cnf, lvl)
			case "logstash":
				hook, err = NewLogstashHook(cnf, lvl)
			default:
				err = fmt.Errorf("Unknown hook type:%s", cnf.Type)
			}
			if err != nil {
				return err
			}
			_logger.Hooks.Add(hook)
		}
	}

	return nil
}

const (
	logHookHelper = `in JSON format: 
		type: logstash|file|stdout|stderr
		level: same as root debug level if omitted, and must above root debug level
	example:
		'[{"type":"logstash", "level":"debug", "param":"udp://<logstash-ip>:<port>","extra_param":"<logstash type>"},
		{"type":"file", "level":"debug", "fmt":"<log format>", "param":"<logdir,should be OS temp dir if not specifed>"},
		{"type":"stdout", "level":"debug"},
		{"type":"stderr", "level":"debug"}]'`
)

func init() {
	flag.StringVar(&fLogFmt, "logfmt", "kv", "kv,json,classic")
	flag.StringVar(&fLogLevel, "loglevel", "info", "debug,info,waring,fatal,panic")
	flag.StringVar(&fLogOptions, "logoptions", "longtime,enableruntime", "string seperated by comma, longtime|shorttime|notime,enableruntime,forcecolors,disablecolors")
	flag.StringVar(&fLogOutput, "logoutput", "stdout", "stdout, stderr, discard")
	flag.StringVar(&fLogHookConfigs, "loghooks", "", logHookHelper)
}

//get logrus logger
// func Logger() *logrus.Logger {
// 	return logger()
// }

func Flush() {
	flush()
}

func InitQLog() {
	_ = logger()
}
