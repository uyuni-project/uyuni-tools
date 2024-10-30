// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
)

var redactRegex = regexp.MustCompile(`([pP]assword[\t :"\\]+)[^\t "\\]+`)

// The default directory where log files are written.
const logDir = "/var/log/"
const logFileName = "uyuni-tools.log"
const GlobalLogPath = logDir + logFileName

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
	// using len(p) prevents "zerolog: could not write event: short write" error
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
	// using len(p) prevents "zerolog: could not write event: short write" error
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
		consoleWriter.NoColor = !term.IsTerminal(int(os.Stdout.Fd()))
		uyuniConsoleWriter := UyuniConsoleWriter{
			consoleWriter: consoleWriter,
		}
		writers = append(writers, uyuniConsoleWriter)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Stack().Logger()

	if fileWriter.logger.Filename != GlobalLogPath {
		log.Warn().Msgf(
			L("Couldn't open %[1]s file for writing, writing log to %[2]s"),
			GlobalLogPath, fileWriter.logger.Filename,
		)
	}
}

func getFileWriter() *UyuniLogger {
	logPath := GlobalLogPath

	if file, err := os.OpenFile(GlobalLogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600); err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logPath = path.Join(".", logFileName)
		} else {
			logPath = path.Join(homeDir, logFileName)
		}
	} else {
		file.Close()
	}

	fileLogger := &lumberjack.Logger{
		Filename:   logPath,
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

func logCallerMarshalFunction(_ uintptr, file string, line int) string {
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
