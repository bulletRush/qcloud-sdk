package qcloud

import (
	"io"
	"bytes"
	"fmt"
	"time"
)

type loggerLevel int

const (
	LvlError loggerLevel = iota
	LvlWarn
	LvlInfo
	LvlDebug
)

type Logger interface {
	Error(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
}

type Formatter interface {
	Format(level loggerLevel, msg string, ctx ...interface{}) []byte
}

type logger struct {
	writer io.Writer
	formatter Formatter
}

func (log logger) Error(msg string, ctx ...interface{}) {
	log.writer.Write(log.formatter.Format(LvlError, msg, ctx...))
}

func (log logger) Warn(msg string, ctx ...interface{}) {
	log.writer.Write(log.formatter.Format(LvlWarn, msg, ctx...))
}

func (log logger) Info(msg string, ctx ...interface{}) {
	log.writer.Write(log.formatter.Format(LvlInfo, msg, ctx...))
}

func (log logger) Debug(msg string, ctx ...interface{}) {
	log.writer.Write(log.formatter.Format(LvlDebug, msg, ctx...))
}

func NewLogger(writer io.Writer, formatter Formatter) Logger {
	return logger{
		writer: writer,
		formatter: formatter,
	}
}

type stdoutFormatter struct {}

func (stdoutFormatter) Format(lvl loggerLevel, msg string, ctx ...interface{}) []byte {
	b := &bytes.Buffer{}
	defaultColor := "\x1b[0m"
	lvlStr, lvlColor := func(lvl loggerLevel) (string, string) {
		switch lvl {
		case LvlError:
			return "ERROR", "\x1b[31m"
		case LvlWarn:
			return "WARN", "\x1b[33m"
		case LvlInfo:
			return "INFO", "\x1b[32m"
		case LvlDebug:
			return "DEBUG", "\x1b[36m"
		}
		return "UNKNOWN", defaultColor
	} (lvl)
	msgColor := "\x1b[34m"
	fmt.Fprintf(b, "%s%s%s[%s] | %s%s%s |", lvlColor, lvlStr, defaultColor,
		time.Now().Format("2006-01-02 15:04:05"), msgColor, msg, defaultColor)
	if len(ctx) % 2 == 1 {
		ctx = append(ctx, "FIX")
	}
	for i:=0; i<len(ctx); i+=2 {
		if i != 0 {
			fmt.Fprint(b, " ")
		}
		fmt.Fprint(b, lvlColor, ctx[i], defaultColor, ":", ctx[i+1])
	}
	b.WriteByte('\n')
	return b.Bytes()
}