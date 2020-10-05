package logging

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"sync"
	"sync/atomic"
)

const (
	DebugLevel = logrus.DebugLevel
	InfoLevel = logrus.InfoLevel
	WarnLevel = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
)

type Fields logrus.Fields
type Logger struct {
	*logrus.Logger
}

var defaultLevel = logrus.InfoLevel
var loggerMap sync.Map

func GetLogger(prefix string) *Logger {
	logger := logrus.New()
	logger.SetFormatter(&prefixed.TextFormatter{

	})
	logger.AddHook(&PrefixHook{prefix:prefix})
	logger.SetLevel(logrus.Level(atomic.LoadUint32((*uint32)(&defaultLevel))))

	loggerWrapper := &Logger{
		Logger: logger,
	}

	result, _ := loggerMap.LoadOrStore(prefix, loggerWrapper)
	return (result).(*Logger)
}

func GetLevel(levelString string) logrus.Level {
	switch levelString {
	case "debug": return DebugLevel
	case "info": return InfoLevel
	case "warn": return WarnLevel
	case "error": return ErrorLevel
	default:
		return InfoLevel
	}
}

func SetGlobalLevel(level logrus.Level) {
	atomic.StoreUint32((*uint32)(&defaultLevel), uint32(level))
	loggerMap.Range(func(_, value interface{}) bool {
		value.(*Logger).SetLevel(defaultLevel)
		return true
	})
}

func SetGlobalOutput(writer io.Writer) {
	loggerMap.Range(func(_, value interface{}) bool {
		value.(*Logger).SetOutput(writer)
		return true
	})
}

type PrefixHook struct {
	prefix string
}

func (h *PrefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *PrefixHook) Fire(e *logrus.Entry) error {
	e.Data["prefix"] = h.prefix
	return nil
}

func (logger Logger) WithFields(fields Fields) *logrus.Entry {
	return logger.Logger.WithFields(logrus.Fields(fields))
}