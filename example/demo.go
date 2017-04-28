package main

import (
	"flag"

	log "github.com/kkkbird/qlog"
)

func main() {
	flag.Parse()
	defer log.Flush()

	// fields := make(logrus.Fields)
	// flag.VisitAll(func(f *flag.Flag) {
	// 	fields[f.Name] = f.Value
	// })

	// if formater, isOk := log.Logger().Formatter.(*logrus.TextFormatter); isOk {
	// 	formater.ForceColors = true
	// }
	// log.WithFields(fields).Infoln("Flags")

	log.Debug("This is a DEBUG message")
	log.Info("This is a INFO message")
	log.Warn("This is a WARN message")
	log.Error("This is a ERROR message")
	log.Fatal("This is a FATAL message")
	log.Panic("This is a PANIC message")
}
