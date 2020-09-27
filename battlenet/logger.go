package battlenet

import (
	"io"
	"io/ioutil"
	"log"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type battlenetLogger interface {
	SetLogOutput(writer io.Writer, levels ...LogLevel)
	getLogger(level LogLevel) *log.Logger
}

type battlenetLoggerImpl struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

func newBattleNetLogger() battlenetLogger {
	bli := &battlenetLoggerImpl{}
	bli.debugLogger = log.New(ioutil.Discard, "DEBUG: ", log.LstdFlags)
	bli.infoLogger = log.New(ioutil.Discard, "INFO: ", log.LstdFlags)
	bli.errorLogger = log.New(ioutil.Discard, "ERROR: ", log.LstdFlags)
	bli.warnLogger = log.New(ioutil.Discard, "WARN: ", log.LstdFlags)
	return bli
}

func (bli *battlenetLoggerImpl) getLogger(level LogLevel) *log.Logger {
	switch level {
	case INFO:
		return bli.infoLogger
	case DEBUG:
		return bli.debugLogger
	case ERROR:
		return bli.errorLogger
	case WARN:
		return bli.warnLogger
	}
	panic("Missing logger in switch")
}

func (bli *battlenetLoggerImpl) SetLogOutput(writer io.Writer, levels ...LogLevel) {
	if len(levels) == 0 {
		levels = []LogLevel{INFO, ERROR, WARN, DEBUG}
	}
	for _, level := range levels {
		bli.getLogger(level).SetOutput(writer)
	}
}
