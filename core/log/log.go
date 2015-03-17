package log

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"runtime"
)

var (
	drunsLog    *syslog.Writer
	FallbackLog *log.Logger
)

func init() {
	FallbackLog = log.New(os.Stderr, "", log.LstdFlags)
}

func Connect(name, hostAndPort string) (err error) {
	drunsLog, err = syslog.Dial("tcp", hostAndPort, syslog.LOG_INFO|syslog.LOG_LOCAL0, name)
	return
}

func ConnectLocal(name string) (err error) {
	drunsLog, err = syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, name)
	return
}

func Disconnect() error {
	if drunsLog == nil {
		return nil
	}

	err := drunsLog.Close()
	if err == nil {
		drunsLog = nil
	}
	return err
}

type Logger struct {
	identifier string
}

func NewLogger(id string) Logger {
	return Logger{"[" + id + "] "}
}

func (l Logger) Alert(m string) {
	logWithSourceInfo(drunsLog.Alert, l.identifier, m)
}

func (l Logger) Crit(m string) {
	logWithSourceInfo(drunsLog.Crit, l.identifier, m)
}

func (l Logger) Critf(m string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Crit, l.identifier, m, a...)
}

func (l Logger) Debug(m string) {
	logWithSourceInfo(drunsLog.Debug, l.identifier, m)
}

func (l Logger) Debugf(m string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Debug, l.identifier, m, a...)
}

func (l Logger) Emerg(m string) {
	logWithSourceInfo(drunsLog.Emerg, l.identifier, m)
}

func (l Logger) Error(e error) {
	msg := l.identifier + e.Error()
	if drunsLog == nil {
		FallbackLog.Println(msg)
		return
	}

	if err := drunsLog.Err(msg); err != nil {
		FallbackLog.Println("Error writing to syslog. Details:", err)
		FallbackLog.Println(msg)
	}
}

func (l Logger) Errorf(s string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Err, l.identifier, s, a...)
}

func (l Logger) Info(m string) {
	logWithSourceInfo(drunsLog.Info, l.identifier, m)
}

func (l Logger) Infof(s string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Info, l.identifier, s, a...)
}

func (l Logger) Notice(m string) {
	logWithSourceInfo(drunsLog.Notice, l.identifier, m)
}

func (l Logger) Warningf(m string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Warning, l.identifier, m, a...)
}

func (l Logger) Warning(m string) {
	logWithSourceInfo(drunsLog.Warning, l.identifier, m)
}

func logWithSourceInfo(f func(string) error, prefix, message string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	msg := fmt.Sprintf("%s%s:%d: %s", prefix, file, line, message)

	if drunsLog == nil {
		FallbackLog.Println(msg)
		return
	}

	if err := f(msg); err != nil {
		FallbackLog.Println("Error writing to syslog. Details:", err)
		FallbackLog.Println(msg)
	}
}

func logWithSourceInfof(f func(string) error, prefix, message string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	args := make([]interface{}, 2, len(a)+2)
	args[0], args[1] = file, line
	args = append(args, a...)
	msg := fmt.Sprintf(prefix+"%s:%d: "+message, args...)

	if drunsLog == nil {
		FallbackLog.Println(msg)
		return
	}

	if err := f(msg); err != nil {
		FallbackLog.Println("Error writing to syslog. Details:", err)
		FallbackLog.Println(msg)
	}
}

func Debug(m string) {
	logWithSourceInfo(drunsLog.Debug, "", m)
}

func Info(m string) {
	logWithSourceInfo(drunsLog.Info, "", m)
}

func Notice(m string) {
	logWithSourceInfo(drunsLog.Notice, "", m)
}

func Warningf(m string, a ...interface{}) {
	logWithSourceInfof(drunsLog.Warning, "", m, a...)
}

func Warning(m string) {
	logWithSourceInfo(drunsLog.Warning, "", m)
}
