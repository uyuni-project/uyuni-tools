package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LogInit(appName string) {
	zerolog.CallerMarshalFunc = logCallerMarshalFunction

	consoleWritter := zerolog.NewConsoleWriter()
	consoleWritter.TimeFormat = time.RFC3339

	multi := zerolog.MultiLevelWriter(consoleWritter, getFileWriter())

	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
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
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

func logCallerMarshalFunction(pc uintptr, file string, line int) string {
	paths := strings.Split(file, "/")
	callerFile := file
	foundSubDir := false
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
