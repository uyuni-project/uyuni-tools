package utils

import (
	"io"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var consoleFiteredWriter FilteredLevelWriter

func LogInit(appName string) {
	zerolog.CallerMarshalFunc = logCallerMarshalFunction
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	consoleWriter := zerolog.NewConsoleWriter()
	consoleFiteredWriter = FilteredLevelWriter{
		Writer: &LevelWriterAdapter{consoleWriter},
		Level:  zerolog.WarnLevel,
	}

	fileWritter := getFileWriter()
	multi := zerolog.MultiLevelWriter(&consoleFiteredWriter, fileWritter)
	log.Logger = zerolog.New(multi).With().Timestamp().Stack().Logger()
	log.Info().Msgf("welcome to %s", appName)
}

func getFileWriter() *lumberjack.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   "/var/log/uyuni-tools.log",
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     90,
		Compress:   true,
	}
	return fileLogger
}

func SetLogLevel(logLevel string) {
	globalLevel := zerolog.InfoLevel
	consoleLevel := zerolog.WarnLevel

	level, err := zerolog.ParseLevel(logLevel)
	if logLevel != "" && err == nil {
		globalLevel = level
		consoleLevel = level
	}
	if globalLevel <= zerolog.DebugLevel {
		log.Logger = log.Logger.With().Caller().Logger()
	}
	zerolog.SetGlobalLevel(globalLevel)
	consoleFiteredWriter.Level = consoleLevel
}

func logCallerMarshalFunction(pc uintptr, file string, line int) string {
	paths := strings.Split(file, "/")
	callerFile := file
	foundSubDir := false
	if strings.HasSuffix(file, "/io/io.go") {
		return "Cmd output"
	}

	for _, currentPath := range paths {
		if foundSubDir {
			if callerFile != "" {
				callerFile = callerFile + "/"
			}
			callerFile = callerFile + currentPath
		} else {
			if strings.Contains(currentPath, "uyuni-tools") {
				foundSubDir = true
				callerFile = ""
			}
		}
	}
	return callerFile + ":" + strconv.Itoa(line)
}

// Anticipating the release of https://github.com/rs/zerolog/pull/573, could be removed once out

type LevelWriterAdapter struct {
	io.Writer
}

func (lw LevelWriterAdapter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	return lw.Write(p)
}

// FilteredLevelWriter writes only logs at Level or above to Writer.
//
// It should be used only in combination with MultiLevelWriter when you
// want to write to multiple destinations at different levels. Otherwise
// you should just set the level on the logger and filter events early.
// When using MultiLevelWriter then you set the level on the logger to
// the lowest of the levels you use for writers.
type FilteredLevelWriter struct {
	Writer zerolog.LevelWriter
	Level  zerolog.Level
}

// Write writes to the underlying Writer.
func (w *FilteredLevelWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}

// WriteLevel calls WriteLevel of the underlying Writer only if the level is equal
// or above the Level.
func (w *FilteredLevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= w.Level {
		return w.Writer.WriteLevel(level, p)
	}
	return len(p), nil
}
