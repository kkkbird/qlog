# qlog

Loggigng system use logrus and logstash hook with some predefined setting

Depends on

* Logrus: https://github.com/sirupsen/logrus
* logrus-logstash-hook: https://github.com/bshuster-repo/logrus-logstash-hook (source code has been integrated to `hook_logstash.go`)

## how to use

``` shell
go get -u -v github.com/sirupsen/logrus
go get -u -v github.com/kkkbird/qlog
```

## example

``` go
package main

import (
	"flag"

	log "github.com/kkkbird/qlog"
)

func main() {
	flag.Parse()
	defer log.Flush()

	log.Debug("This is a DEBUG message")
	log.Info("This is a INFO message")
	log.Warn("This is a WARN message")
}

```

build the program and run with `<program_name> -h` to see the flags

## how to use

for example:

``` shell
go run example/demo.go -loglevel=debug -loghooks='[{"type":"logstash","param":"udp://192.168.0.151:5020"},{"type":"file","level":"info","fmt":"classic","param":"."},{"type":"stderr","level":"error"}]' -logoutput=discard -logoptions=enableruntime,forcecolors
```

explaination:
1. root log level set to debug
1. add logstash hook and send log by udp, at same log level <DEBUG> with root
1. add file hook and write to file, at <INFO> level, with <CLASSIC> format
1. add stderr hook, at <ERROR> level
