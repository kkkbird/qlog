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

// type Logger = logrus.Logger
type Logger struct {
	*logrus.Logger
}

var (
	v *viper.Viper

	// DefaultLogger : default logger object
	gDefaultLogger *Logger

	// root logger settings
	//rootCfg rootConfig

	// flagset
	cli = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
)

const (
	keyRootReportCaller = "logger.reportcaller"
	keyRootLevel        = "logger.level"
)

func setDefault() {
	v.SetDefault(keyRootLevel, "debug")
	v.SetDefault(keyRootReportCaller, false)
}

func initRootFlags() error {
	return nil
}

func initViper() error {
	v = viper.New()

	// read from flags
	cli.Parse(os.Args[1:])
	v.BindPFlags(cli)

	// read from env
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// set default
	setDefault()

	// read from logger.yaml
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/qlog")
	v.AddConfigPath("./conf/")
	v.SetConfigName("logger")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		// no config file
		// return err
	} else {
		// watch configs changes
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("[qlog] config changed: ", e.Name)
			resetLogger()
		})
	}

	return nil
}

func resetLogger() {
	if err := initLogger(); err != nil {

		gDefaultLogger = &Logger{
			Logger: logrus.StandardLogger(),
		}

		fmt.Println("[qlog] reload config fail:", err)
	}
}

func getDefaultFormatter() (logrus.Formatter, error) {
	var err error
	var formatter logrus.Formatter
	defaultFormatterName := v.GetString("logger.formatter.name")

	if defaultFormatterName == "" {
		formatter = &logrus.TextFormatter{}
	} else {
		if formatter, err = newFormatter(defaultFormatterName, "logger.formatter.opts"); err != nil {
			return nil, err
		}
	}

	return formatter, nil
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

func initLogger() error {
	var err error

	reportCaller := v.GetBool(keyRootReportCaller)

	level, err := logrus.ParseLevel(v.GetString(keyRootLevel))
	if err != nil {
		return fmt.Errorf("get default log level error: %s", err)
	}

	formatter, err := getDefaultFormatter()
	if err != nil {
		return fmt.Errorf("get default formatters error: %s", err)
	}

	gDefaultLogger = &Logger{
		Logger: &logrus.Logger{
			Out:          ioutil.Discard,
			Formatter:    formatter,
			Hooks:        nil,
			Level:        level,
			ExitFunc:     os.Exit,
			ReportCaller: reportCaller,
		},
	}

	hooks, err := getActivateHooks()

	if err != nil {
		fmt.Printf("[qlog] get hooks error: %s\n", err)
		gDefaultLogger.SetOutput(os.Stderr)
		return nil
	}

	gDefaultLogger.ReplaceHooks(hooks)
	return nil
}

func init() {
	var err error

	if err = initSysParams(); err != nil {
		panic(fmt.Sprint("[qlog] init system param error:", err))
	}

	if err = initViper(); err != nil {
		panic(fmt.Sprint("[qlog] init viper error:", err))
	}

	if err = initLogger(); err != nil {
		panic(fmt.Sprint("[qlog] initLogger fail:", err))
	}
}
