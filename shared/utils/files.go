// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

var preserveTmpDir atomic.Bool

func SetShouldPreserveTmpDir(shouldPreserve bool) {
	preserveTmpDir.Store(shouldPreserve)
}

// IsEmptyDirectory return true if a given directory is empty.
func IsEmptyDirectory(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("cannot check content of %s"), path)
		return false
	}
	if len(files) > 0 {
		return false
	}
	return true
}

// RemoveDirectory remove a given directory.
func RemoveDirectory(path string) error {
	if err := os.Remove(path); err != nil {
		return Errorf(err, L("Cannot remove %s folder"), path)
	}
	return nil
}

// FileExists check if path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if !os.IsNotExist(err) {
		log.Fatal().Err(err).Msgf(L("Failed to get %s file informations"), path)
	}
	return false
}

// ReadFile returns the content of a file and exit if there was an error.
func ReadFile(file string) []byte {
	out, err := os.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to read file %s"), file)
	}
	return out
}

// GetFileBoolean gets the value of a file containing a boolean.
//
// This is handy for files from the kernel API.
func GetFileBoolean(file string) bool {
	return strings.TrimSpace(string(ReadFile(file))) != "0"
}

// UninstallFile uninstalls a file.
func UninstallFile(path string, dryRun bool) {
	if FileExists(path) {
		if dryRun {
			log.Info().Msgf(L("Would remove file %s"), path)
		} else {
			log.Info().Msgf(L("Removing file %s"), path)
			if err := os.Remove(path); err != nil {
				log.Info().Err(err).Msgf(L("Failed to remove file %s"), path)
			}
		}
	}
}

// TempDir creates a temporary directory.
func TempDir() (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return "", nil, Error(err, L("failed to create temporary directory"))
	}

	shouldPreserveTmpDir := preserveTmpDir.Load()
	cleaner := func() {
		if shouldPreserveTmpDir {
			log.Info().Msgf(L("Generated temporary directory will be preserved: %s"), tempDir)
		} else {
			if err := os.RemoveAll(tempDir); err != nil {
				log.Error().Err(err).Msg(L("failed to remove temporary directory"))
			}
		}
	}
	return tempDir, cleaner, nil
}

// CopyFile copies the content of the file at src path into the opened dst file.
func CopyFile(src string, dst *os.File) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return Errorf(err, L("fails to open %s file"), src)
	}

	const bufSize = 1024
	buf := make([]byte, bufSize)

	for {
		read, err := srcFile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return Errorf(err, L("failed to read %s file"), src)
		}
		_, err = dst.Write(buf[:read])
		if err != nil {
			return Errorf(err, L("failed to copy %s file"), src)
		}
	}
}
