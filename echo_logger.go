package main

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

func NewEchoLogger(logger *logrus.Entry) echo.Logger {
	return &echoLogger{logger}
}

type echoLogger struct {
	*logrus.Entry
}

func (l *echoLogger) Debugj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Debug()
}

func (l *echoLogger) Errorj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Error()
}

func (l *echoLogger) Fatalj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Fatal()
}

func (l *echoLogger) Infoj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Info()
}

func (l *echoLogger) Panicj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Panic()
}

func (l *echoLogger) Printj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Print()
}

func (l *echoLogger) Warnj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Warn()
}

func (l *echoLogger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.FatalLevel:
		return log.ERROR
	case logrus.PanicLevel:
		return log.ERROR
	}

	return log.OFF
}

func (l *echoLogger) Output() io.Writer {
	return l.Logger.Out
}

func (l *echoLogger) Prefix() string {
	return ""
}

func (l *echoLogger) SetLevel(log.Lvl) {

}

func (l *echoLogger) SetOutput(io.Writer) {

}

func (l *echoLogger) SetPrefix(string) {

}

func (l *echoLogger) SetHeader(string) {

}
