package qlog

// StandardLogger returns default logger, direct access to Logger object is not recommended in qlog
// func StandardLogger() *logrus.Logger {
// 	return logrus.StandardLogger()
// }

// FilterFlags filter --logger.xxx flags
func FilterFlags(args []string) []string {
	return filterLoggerFlags(args, false)
}
