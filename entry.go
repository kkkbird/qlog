package qlog

import (
	"github.com/sirupsen/logrus"
)

type Entry struct {
	e                 *logrus.Entry
	withRunTimeFields bool
}

func (entry *Entry) String() (string, error) {
	serialized, err := entry.e.Logger.Formatter.Format(entry.e)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return &Entry{
		e:                 entry.e.WithError(err),
		withRunTimeFields: entry.withRunTimeFields,
	}
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{
		e:                 entry.e.WithField(key, value),
		withRunTimeFields: entry.withRunTimeFields,
	}
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields logrus.Fields) *Entry {
	return &Entry{
		e:                 entry.e.WithFields(fields),
		withRunTimeFields: entry.withRunTimeFields,
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Debug(args...)
	} else {
		entry.e.Debug(args...)
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Info(args...)
	} else {
		entry.e.Info(args...)
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Info(args...)
	} else {
		entry.e.Info(args...)
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warn(args...)
	} else {
		entry.e.Warn(args...)
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warn(args...)
	} else {
		entry.e.Warn(args...)
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Error(args...)
	} else {
		entry.e.Error(args...)
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Fatal(args...)
	} else {
		entry.e.Fatal(args...)
	}
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Panic(args...)
	} else {
		entry.e.Panic(args...)
	}
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Debugf(format, args...)
	} else {
		entry.e.Debugf(format, args...)
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		entry.e.Infof(format, args...)
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Infof(format, args...)
	} else {
		entry.e.Infof(format, args...)
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		entry.e.Warnf(format, args...)
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warnf(format, args...)
	} else {
		entry.e.Warnf(format, args...)
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Errorf(format, args...)
	} else {
		entry.e.Errorf(format, args...)
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Fatalf(format, args...)
	} else {
		entry.e.Fatalf(format, args...)
	}
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Panicf(format, args...)
	} else {
		entry.e.Panicf(format, args...)
	}
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Debugln(args...)
	} else {
		entry.e.Debugln(args...)
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		entry.e.Infoln(args...)
	}
}

func (entry *Entry) Println(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Infoln(args...)
	} else {
		entry.e.Infoln(args...)
	}
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		entry.e.Warnln(args...)
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Warnln(args...)
	} else {
		entry.e.Warnln(args...)
	}
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Errorln(args...)
	} else {
		entry.e.Errorln(args...)
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Fatalln(args...)
	} else {
		entry.e.Fatalln(args...)
	}
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.withRunTimeFields {
		entry.e.WithFields(runtimeFields(2)).Panicln(args...)
	} else {
		entry.e.Panicln(args...)
	}
}
