package log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	log        *log.Logger
	DebugLevel DebugLevel
}

var (
	logger *Logger

	ErrWriterNil = errors.New("logger: writer is nil")
)

func New(out io.Writer, debugLevel DebugLevel) error {
	if out == nil {
		return ErrWriterNil
	}

	defaultFlag := log.LstdFlags

	if debugLevel&DebugLevelError == DebugLevelError {
		defaultFlag |= log.Llongfile
	} else if debugLevel&DebugLevelWarn == DebugLevelWarn {
		defaultFlag |= log.Lshortfile
	}

	logger = &Logger{log.New(out, "", defaultFlag), debugLevel}

	return nil
}

func Error(err any) {
	if logger.DebugLevel&DebugLevelError != DebugLevelError {
		return
	}

	out := fmt.Append(nil, err)
	logger.log.Output(2, ColorRed+"[ERROR] - "+ColorReset+string(out)+"\n")
}

func Errorf(format string, v ...any) {
	if logger.DebugLevel&DebugLevelError != DebugLevelError {
		return
	}

	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, ColorRed+"[ERROR] - "+ColorReset+string(out)+"\n")
}

func Warn(v ...any) {
	if logger.DebugLevel&DebugLevelWarn != DebugLevelWarn {
		return
	}

	logger.log.Printf(ColorYellow+"[WARN] - "+ColorReset+"%s\n", v)
}

func Warnf(format string, v ...any) {
	if logger.DebugLevel&DebugLevelWarn != DebugLevelWarn {
		return
	}

	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, ColorYellow+"[WARN] - "+ColorReset+string(out)+"\n")
}

func Info(v ...any) {
	if logger.DebugLevel&DebugLevelInfo != DebugLevelInfo {
		return
	}

	logger.log.Printf(ColorGreen+"[INFO] - "+ColorReset+"%s\n", v)
}

func Infof(format string, v ...any) {
	if logger.DebugLevel&DebugLevelInfo != DebugLevelInfo {
		return
	}

	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, ColorGreen+"[INFO] - "+ColorReset+string(out)+"\n")
}

func Debug(format string, v ...any) {
	if logger.DebugLevel&DebugLevelNone == DebugLevelNone {
		return
	}

	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, ColorBlue+"[DEBUG] - "+ColorReset+string(out)+"\n")
}

func Print(v ...any) {
	out := fmt.Append(nil, v...)
	logger.log.Output(2, string(out)+"\n")
}

func Printf(format string, v ...any) {
	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, string(out)+"\n")
}

func Fatal(v ...any) {
	out := fmt.Append(nil, v...)
	logger.log.Output(2, ColorRed+"[FATAL] - "+string(out)+ColorReset+"\n")
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	out := fmt.Appendf(nil, format, v...)
	logger.log.Output(2, ColorRed+"[FATAL] - "+string(out)+ColorReset+"\n")
	os.Exit(1)
}
