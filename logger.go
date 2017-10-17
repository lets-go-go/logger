package logger

import (
	"fmt"
)

var (
	fwLogger *Logger
)

func init() {
	fwLogger = NewLoggerWithConfig(DEFAULT_CONFIG)
}

// Init 初始化
func Init(config *Config) {
	fwLogger = NewLoggerWithConfig(config)
}

// InitWithConf 初始化
func InitWithConf(conf string) {

	fwLogger = NewLogger(conf)
}

// Trace logs with the TRACE severity.
// Arguments are handled in the manner of fmt.Print.
func Trace(v ...interface{}) {
	fwLogger.Output(TRACE, fmt.Sprint(v...))
}

// Traceln logs with the TRACE severity.
// Arguments are handled in the manner of fmt.Println.
func Traceln(v ...interface{}) {
	fwLogger.Output(TRACE, fmt.Sprintln(v...))
}

// Tracef logs with the TRACE severity.
// Arguments are handled in the manner of fmt.Printf.
func Tracef(format string, v ...interface{}) {
	fwLogger.Output(TRACE, fmt.Sprintf(format, v...))
}

// Debug logs with the DEBUG severity.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	fwLogger.Output(DEBUG, fmt.Sprint(v...))
}

// Debugln logs with the DEBUG severity.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	fwLogger.Output(DEBUG, fmt.Sprintln(v...))
}

// Debugf logs with the DEBUG severity.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	fwLogger.Output(DEBUG, fmt.Sprintf(format, v...))
}

// Info logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	fwLogger.Output(INFO, fmt.Sprint(v...))
}

// Infoln logs with the INFO severity.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	fwLogger.Output(INFO, fmt.Sprintln(v...))
}

// Infof logs with the INFO severity.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	fwLogger.Output(INFO, fmt.Sprintf(format, v...))
}

// Warn logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	fwLogger.Output(WARN, fmt.Sprint(v...))
}

// Warnln logs with the WARN severity.
// Arguments are handled in the manner of fmt.Println.
func Warnln(v ...interface{}) {
	fwLogger.Output(WARN, fmt.Sprintln(v...))
}

// Warnf logs with the WARN severity.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	fwLogger.Output(WARN, fmt.Sprintf(format, v...))
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	fwLogger.Output(ERROR, fmt.Sprint(v...))
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	fwLogger.Output(ERROR, fmt.Sprintln(v...))
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	fwLogger.Output(ERROR, fmt.Sprintf(format, v...))
}

// Fatal logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func Fatal(v ...interface{}) {
	fwLogger.Output(FATAL, fmt.Sprint(v...))
}

// Fatalln logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func Fatalln(v ...interface{}) {
	fwLogger.Output(FATAL, fmt.Sprintln(v...))

}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	fwLogger.Output(FATAL, fmt.Sprintf(format, v...))
}
