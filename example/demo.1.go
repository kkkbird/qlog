package main

import (
	_ "github.com/kkkbird/qlog"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Debug("This is a DEBUG message")
	log.Info("This is a INFO message")
	log.Warn("This is a WARN message")
	log.Error("This is a ERROR message")
}
