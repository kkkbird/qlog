package main

import (
	"errors"

	"github.com/kkkbird/qlog"
	"github.com/sirupsen/logrus"
)

var log = qlog.DefaultLogger

func main() {
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
	//log.Fatal("This is a FATAL message")
	//log.Panic("This is a PANIC message")
	log.WithField("foo", "bar").Warn("This is a WithField WARN message")
	log.WithFields(logrus.Fields{
		"hello":  "world",
		"goobye": "moon",
	}).Info("This is a WithFields INFO message")
	log.WithError(errors.New("An error")).Warn("with error warning")

	entry := log.WithField("test", "1")
	entry.Debug("This is a DEBUG message from entry")
	entry.Info("This is a INFO message from entry")

}
