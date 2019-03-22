package qlog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Logger is an alias of logrus.Logger
type Logger = logrus.Logger

// Fields is an alias of logrus.Fields
type Fields = logrus.Fields

// type Logger struct {
// 	*logrus.Logger
// }

type loggerConfig struct {
	Path []string
	Name string
	Typ  string
	File string
}

var (
	// flagset
	cli = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// viper
	v = viper.New()

	// DefaultLogger : default logger object
	qLogger = StandardLogger()

	qLoggerConfig = &loggerConfig{
		Path: make([]string, 0),
	}
)

const (
	keyConfigPath = "logger.config.path"
	keyConfigName = "logger.config.name"
	keyConfigType = "logger.config.type"
	keyConfigFile = "logger.config.file"

	keyReportCaller         = "logger.reportcaller"
	keyDefaultLevel         = "logger.level"
	keyDefaultFormatterName = "logger.formatter.name"
	keyDefaultFormatterOpts = "logger.formatter.opts"
)

func setDefault() {
	v.SetDefault(keyReportCaller, false)
	v.SetDefault(keyDefaultLevel, "debug")
	v.SetDefault(keyDefaultFormatterName, "text")
}

func initFlags() error {
	cli.Bool(keyReportCaller, false, "logger.reportcaller")
	cli.String(keyDefaultLevel, "error", "logger.level")
	cli.String(keyDefaultFormatterName, "text", "logger.formatter.name")

	cli.StringSliceVar(&qLoggerConfig.Path, keyConfigPath, []string{".", "./conf", "/etc/qlog"}, "logger.config.path")
	cli.StringVar(&qLoggerConfig.Name, keyConfigName, "logger", "logger.config.name")
	cli.StringVar(&qLoggerConfig.Typ, keyConfigType, "yaml", "logger.config.type")

	cli.StringVar(&qLoggerConfig.File, keyConfigFile, "", "logger.config.file")

	return nil
}

func initViper() error {
	// read from flags
	cli.Parse(filterLoggerFlags(os.Args[1:], true))
	v.BindPFlags(cli)

	// read from env
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// set default
	setDefault()

	// read from config file
	if len(qLoggerConfig.File) > 0 {
		v.SetConfigFile(qLoggerConfig.File)
	} else {
		for _, p := range qLoggerConfig.Path {
			v.AddConfigPath(p)
		}
		v.SetConfigName(qLoggerConfig.Name)
	}
	v.SetConfigType(qLoggerConfig.Typ)

	if err := v.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// no config file
		default:
			return err
		}
	} else {
		// watch configs changes
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("[qlog] config changed: ", e.Name)
			resetLogger()
		})
	}

	//v.Debug()

	return nil
}

func resetLogger() {
	if err := configLogger(); err != nil {
		fmt.Printf("[qlog] reload config fail:%s, changes may not take effect!", err)
	}
}

func getDefaultFormatter() (logrus.Formatter, error) {
	return newFormatter(v.GetString(keyDefaultFormatterName), keyDefaultFormatterOpts)
}

func getActivateHooks() (logrus.LevelHooks, error) {
	var err error
	var hook logrus.Hook
	var activateHooks = make(logrus.LevelHooks)

	for name := range gRegisteredHooks {
		n := strings.Join([]string{"logger", name, "enabled"}, ".")
		//fmt.Println("initHooks", n, v.GetBool(n))
		if v.GetBool(n) == true {
			if hook, err = newHook(name); err != nil {
				fmt.Printf("[qlog] init hook(%s) error:%s\n", name, err)
				continue
			}
			activateHooks.Add(hook)
		}
	}

	if len(activateHooks) == 0 {
		return nil, errors.New("no activate log hook")
	}

	return activateHooks, nil
}

func configLogger() error {
	var err error

	qLogger.SetReportCaller(v.GetBool(keyReportCaller))

	level, err := logrus.ParseLevel(v.GetString(keyDefaultLevel))

	if err != nil {
		return fmt.Errorf("get default log level error: %s", err)
	}
	qLogger.SetLevel(level)

	formatter, err := getDefaultFormatter()
	if err != nil {
		return fmt.Errorf("get default formatters error: %s", err)
	}
	qLogger.SetFormatter(formatter)

	// SetLevel and SetFormatter must be called before getActivateHooks.
	hooks, err := getActivateHooks()

	if err != nil {
		fmt.Printf("[qlog] get hooks error: %s\n", err)
		qLogger.ReplaceHooks(nil)
		qLogger.SetOutput(os.Stderr)
		return nil
	}

	qLogger.SetOutput(ioutil.Discard)
	qLogger.ReplaceHooks(hooks)
	return nil
}

func init() {
	var err error

	if err = initSysParams(); err != nil {
		panic(fmt.Sprint("[qlog] init system param error:", err))
	}

	if err = initFlags(); err != nil {
		panic(fmt.Sprint("[qlog] init flags error:", err))
	}

	if err = initViper(); err != nil {
		panic(fmt.Sprint("[qlog] init viper error:", err))
	}

	if err = configLogger(); err != nil {
		panic(fmt.Sprint("[qlog] configLogger fail:", err))
	}
}
