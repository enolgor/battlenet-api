package battlenet

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type LogLevel string

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

type battlenetLogger interface {
	SetLogger(logger *log.Logger, levels ...LogLevel)
	getLogger(level LogLevel) bncLogger
}

type battlenetLoggerImpl struct {
	loggers map[LogLevel]*bncLoggerImpl
}

func newBattleNetLogger() battlenetLogger {
	bli := &battlenetLoggerImpl{loggers: make(map[LogLevel]*bncLoggerImpl)}
	bli.loggers[DEBUG] = &bncLoggerImpl{level: DEBUG, internal: nil}
	bli.loggers[INFO] = &bncLoggerImpl{level: INFO, internal: nil}
	bli.loggers[ERROR] = &bncLoggerImpl{level: ERROR, internal: nil}
	bli.loggers[WARN] = &bncLoggerImpl{level: WARN, internal: nil}
	return bli
}

func (bli *battlenetLoggerImpl) getLogger(level LogLevel) bncLogger {
	log, ok := bli.loggers[level]
	if !ok {
		log = bli.loggers[INFO]
	}
	return log
}

func (bli *battlenetLoggerImpl) SetLogger(logger *log.Logger, levels ...LogLevel) {
	if len(levels) == 0 {
		levels = []LogLevel{INFO, ERROR, WARN, DEBUG}
	}
	for _, level := range levels {
		bli.loggers[level].internal = logger
	}
}

type bncLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Writer() io.Writer
}

type bncLoggerImpl struct {
	level    LogLevel
	internal *log.Logger
}

func (bli *bncLoggerImpl) Print(v ...interface{}) {
	if bli.internal == nil {
		return
	}
	prefix := bli.internal.Prefix()
	bli.internal.SetPrefix(fmt.Sprintf("[%s]", bli.level))
	bli.internal.Print(v...)
	bli.internal.SetPrefix(prefix)
}

func (bli *bncLoggerImpl) Printf(format string, v ...interface{}) {
	if bli.internal == nil {
		return
	}
	prefix := bli.internal.Prefix()
	bli.internal.SetPrefix(fmt.Sprintf("[%s]", bli.level))
	bli.internal.Printf(format, v...)
	bli.internal.SetPrefix(prefix)
}

func (bli *bncLoggerImpl) Println(v ...interface{}) {
	if bli.internal == nil {
		return
	}
	prefix := bli.internal.Prefix()
	bli.internal.SetPrefix(fmt.Sprintf("[%s]", bli.level))
	bli.internal.Println(v...)
	bli.internal.SetPrefix(prefix)
}

func (bli *bncLoggerImpl) Writer() io.Writer {
	if bli.internal == nil {
		return ioutil.Discard
	}
	return bli.internal.Writer()
}
