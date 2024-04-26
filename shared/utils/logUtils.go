// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
)

var redactRegex = regexp.MustCompile(`(password\s+)[^\s"]+`)

// UyuniLogger is an io.WriteCloser that writes to the specified filename.
type UyuniLogger struct {
	logger *lumberjack.Logger
}

// UyuniConsoleWriter parses the JSON input and writes it in an (optionally) colorized, human-friendly format to Out.
type UyuniConsoleWriter struct {
	consoleWriter zerolog.ConsoleWriter
}

func (l *UyuniLogger) Write(p []byte) (n int, err error) {
	_, err = l.logger.Write([]byte(redact(string(p)) + "\n"))
	if err != nil {
		return 0, err
	}
	//using len(p) prevents "zerolog: could not write event: short write" error
	return len(p), nil
}

// Close implements io.Closer, and closes the current logfile.
func (l *UyuniLogger) Close() error {
	return l.logger.Close()
}

// Rotate causes Logger to close the existing log file and immediately create a
// new one.  This is a helper function for applications that want to initiate
// rotations outside of the normal rotation rules, such as in response to
// SIGHUP.  After rotating, this initiates compression and removal of old log
// files according to the configuration.
func (l *UyuniLogger) Rotate() error {
	return l.logger.Rotate()
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (c UyuniConsoleWriter) Write(p []byte) (n int, err error) {
	_, err = c.consoleWriter.Write([]byte(redact(string(p))))
	if err != nil {
		return 0, err
	}
	//using len(p) prevents "zerolog: could not write event: short write" error
	return len(p), nil
}

func redact(line string) string {
	return redactRegex.ReplaceAllString(line, "${1}<REDACTED>")
}

// LogInit initialize logs.
func LogInit(logToConsole bool) {
	zerolog.CallerMarshalFunc = logCallerMarshalFunction
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	fileWriter := getFileWriter()
	writers := []io.Writer{fileWriter}
	if logToConsole {
		consoleWriter := zerolog.NewConsoleWriter()
		uyuniConsoleWriter := UyuniConsoleWriter{
			consoleWriter: consoleWriter,
		}
		consoleWriter.NoColor = !term.IsTerminal(int(os.Stdout.Fd()))
		writers = append(writers, uyuniConsoleWriter)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Stack().Logger()
}

func getFileWriter() *UyuniLogger {
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
	uyuniLogger := &UyuniLogger{
		logger: fileLogger,
	}
	return uyuniLogger
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
