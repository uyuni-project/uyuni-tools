// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogInit initialize logs.
func LogInit(logToConsole bool) {
	zerolog.CallerMarshalFunc = logCallerMarshalFunction
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	fileWriter := getFileWriter()
	writers := []io.Writer{fileWriter}
	if logToConsole {
		consoleWriter := zerolog.NewConsoleWriter()
		consoleWriter.NoColor = !term.IsTerminal(int(os.Stdout.Fd()))
		writers = append(writers, consoleWriter)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Stack().Logger()
}

func getFileWriter() *lumberjack.Logger {
	const globalLogPath = "/var/log/"
	logPath := globalLogPath

	if file, err := os.OpenFile(globalLogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600); err != nil {
		logPath, err = os.UserHomeDir()
		if err != nil {
			logPath = "./"
		}
	} else {
		file.Close()
	}

	fileLogger := &lumberjack.Logger{
		Filename:   path.Join(logPath, "uyuni-tools.log"),
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     90,
		Compress:   true,
	}
	return fileLogger
}

// SetLogLevel sets the loglevel.
func SetLogLevel(logLevel string) {
	globalLevel := zerolog.InfoLevel

	level, err := zerolog.ParseLevel(logLevel)
	if logLevel != "" && err == nil {
		globalLevel = level
	}
	if globalLevel <= zerolog.DebugLevel {
		log.Logger = log.Logger.With().Caller().Logger()
	}
	zerolog.SetGlobalLevel(globalLevel)
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
