package clog

var logger loggerOps = newZapLogger()

// loggerActions defines the available logging actions a logger needs to support
type loggerActions interface {
	Infof(msg string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnf(msg string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorf(msg string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

// loggerOps defines the operations a logger needs to support
type loggerOps interface {
	loggerActions
	Flush()
}

// Info records an information level log
func Info(msg string) {
	logger.Infof(msg)
}

// Infof records an information level log via a formatted string
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Infow records an information level log with key-value metadata
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// Warn records a warning level log
func Warn(msg string) {
	logger.Warnf(msg)
}

// Warnf records a warning level log via a formatted string
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Warnw records a warning level log with key-value metadata
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

// Error records an error level log
func Error(msg string) {
	logger.Errorf(msg)
}

// Errorf records an error level log via a formatted string
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// Errorw records an error level log with key-value metadata
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

// Flush will ensure any in-flight logs are flushed. Should be called when the app is shutting down
func Flush() {
	logger.Flush()
}
