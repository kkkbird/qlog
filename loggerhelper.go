package qlog

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	// system param
	gPid      = os.Getpid()
	gProgram  = filepath.Base(os.Args[0])
	gHost     = "unknownhost"
	gUserName = "unknownuser"
)

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

func initSysParams() error {
	h, err := os.Hostname()
	if err == nil {
		gHost = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		gUserName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	gUserName = strings.Replace(gUserName, `\`, "_", -1)

	return nil
}

func filterLoggerFlags(args []string, keep bool) []string {
	rlt := make([]string, 0)

	for _, arg := range args {
		if strings.HasPrefix(arg, "--logger.") {
			if keep {
				rlt = append(rlt, arg)
			}
		} else {
			if !keep {
				rlt = append(rlt, arg)
			}
		}
	}

	return rlt
}
