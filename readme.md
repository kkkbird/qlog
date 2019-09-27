# qlog

Logging system based on [sirupsen/logrus](https://github.com/sirupsen/logrus)

## Simple example

``` go
package main

import (
  "flag"

  _ "github.com/kkkbird/qlog" // call qlog hijack
  log "github.com/sirupsen/logrus"
)

func main() {
  log.Debug("This is a DEBUG message")
  log.Info("This is a INFO message")
  log.Warn("This is a WARN message")
}
```

## Configurations

`qlog` use [spf13/viper](https://github.com/spf13/viper) for  configuration

### config via FILE

example for config file

``` yaml
logger:
  level: debug # default log level, used if hooks not specify a log level
  reportcaller: true # show reportcaller, will affect all hooks
  formatter: # default formatter, used if hooks not specify a formatter
    name: text # default formatter name
    opts: # default formatter opts
      forcecolors: false
      prettycaller: truncated

  stdout:
    enabled: true
    level: error

  file:
    enabled: true
    path: ./log
    level: trace
    formatter:
      name: classic
      prettycaller: omitfunc

```

you can set config file by following config flags

* logger.config.path: paths for configuration, default search paths are `.,./conf,/etc/qlog`
* logger.config.name: name for configuration file, default `logger`
* logger.config.type: type for configuration file, default `yaml`
* logger.config.file: file for configuration file, default is empty, if it is set, above configuration will be ignored

for example

``` shell
<app> --logger.config.file=./conf/qlog.yml
```

logger configuration file will be watched, if it is changed in runtime, `qlog` wil reload the module.

### config via ENV

you can set environment to override configuration in config file, for example

``` shell
export LOGGER_REPORTCALLER=true
export LOGGER_STDOUT_LEVEL=trace

<app>
```

### config via FLAGS

base on [spf13/pflag](https://github.com/spf13/pflag), you can add flags to override configuration in config file, for example

``` shell
<app> --logger.level=debug --logger.formatter.name=classic
```

### precedence from high to low

* flag
* env
* config
* default

## Formatters

all formatters will have a `name` field and several `opts` fields, example:

``` yaml
formatter:
  name: text
  opts:
    forcecolors: false
```

formatters can be set

* via File: yes
* via ENV: yes
* via FLAG: no, only support set default `logger.formatter.name`, hooks' formatter cannot be set via flag

### special formaater options

#### prettycaller

if set and formatter has field `CallerPrettyfier`, caller info will call predefined caller pretty function to output caller info. Only available when `reportcaller` is true

valid options are listed below

* omitfunc
* truncated

### NullFormatter

logger won't formatted and output nothing

### TextFormatter

as `logrus` TextFormatter

### JSONFormatter

as `logrus` JSONFormatter

### ClassicFormatter

a classic logger formatter

``` shell
2018/12/24 19:15:28.307516+08:00 [D] This is a DEBUG message
2018/12/24 19:15:28.307563+08:00 [I] This is a INFO message
2018/12/24 19:15:28.307801+08:00 [W] This is a WARN message
2018/12/24 19:15:28.307807+08:00 [E] This is a ERROR message
2018/12/24 19:15:28.307819+08:00 [W] This is a WithField WARN message foo=bar
```

## Hooks

All hooks will have fields following fields

* enabled
* level
* formatter

A hook will not be used if `enabled` is false. And it will used default setting in top level if level or formatter is no set

Deferent hook can have its own configration field, for example

``` yaml
logger:
  file:
    enabled: true
    level: trace
    formatter:
      name: classic
    path:
    - ./log
```

### StdoutHook

* logger.stdout.enabled
* logger.stdout.level

### StderrHook

* logger.stderr.enabled
* logger.stderr.level

### FileHook

* logger.file.enabled
* logger.file.level
* logger.file.path: log file paths
* logger.file.name: log file name
* logger.file.rotate.time: rotate duration, default is "24h", set 0 to disable rotate
* logger.file.rotate.maxage: time to wait until old logs are purged, default 168h(7 days), set 0 to disable
* logger.file.rotate.count: the number of files should be kept, default 0 and count is disabled

### UDPHook

* logger.udp.enabled
* logger.udp.level
* logger.udp.host
* logger.udp.uuid

## HOWTO

### Common use

`qlog` hijack the `logrus` default StandardLogger(), so if your project use `logrus` without init logger object your self, you can just add one line code in your main package as example and keep other codes untouched

### Use flags

If you want use goflag in the app with qlog, you must set your own flag parser now and call logger flag at the end of cmd line. 

For example,

your app.go

```go
package main

import (
  "flag"
  "os"

  _ "github.com/kkkbird/qlog"
  log "github.com/sirupsen/logrus"
)

var cli = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

var appHost string

func init() {
  cli.StringVar(&appHost, "apphost", "https://apphost.cn", "App host url")
}

func main() {
  cli.Parse(log.FilterFlags(os.Args[1:]))
  log.Debug("app host:", appHost)
}
```
