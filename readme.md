# qlog

Loggigng system use logrus and logstash hook with some predefined setting

Depends on  
* Logrus: https://github.com/sirupsen/logrus
* logrus-logstash-hook: https://github.com/bshuster-repo/logrus-logstash-hook
## how to use
``` shell
go get -u -v github.com/sirupsen/logrus
go get -u -v github.com/kkkbird/logrus-logstash-hook
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



