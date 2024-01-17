// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

const PROMPT_END = ": "

func AskPasswordIfMissing(value *string, prompt string) {
	if *value == "" {
		fmt.Print(prompt + PROMPT_END)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read password")
		}
		*value = string(bytePassword)
		fmt.Println()
	}
}

func AskIfMissing(value *string, prompt string) {
	if *value == "" {
		fmt.Print(prompt + PROMPT_END)
		reader := bufio.NewReader(os.Stdin)
		newValue, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read input")
		}
		*value = newValue
		fmt.Println()
	}
}

// ComputeImage assembles the container image from its name and tag.
func ComputeImage(name string, tag string) (string, error) {
	imageValid := regexp.MustCompile("^((?:[^:/]+(?::[0-9]+)?/)?[^:]+)(?::([^:]+))?$")
	submatches := imageValid.FindStringSubmatch(name)
	if submatches == nil {
		return "", fmt.Errorf("invalid image name: %s", name)
	}
	if submatches[2] == "" {
		// No tag provided in the URL name, append the one passed
		return fmt.Sprintf("%s:%s", name, tag), nil
	}
	return name, nil
}

// Get the timezone set on the machine running the tool
func GetLocalTimezone() string {

	out, err := RunCmdOutput(zerolog.DebugLevel, "timedatectl", "show", "--value", "-p", "Timezone")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run timedatectl show --value -p Timezone")
	}
	return string(out)
}

// Check if a given path exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if !os.IsNotExist(err) {
		log.Fatal().Err(err).Msgf("Failed to stat %s file", path)
	}
	return false
}

// Returns the content of a file and exit if there was an error.
func ReadFile(file string) []byte {
	out, err := os.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to read file %s", file)
	}
	return out
}

// Get the value of a file containing a boolean.
// This is handy for files from the kernel API
func GetFileBoolean(file string) bool {
	return string(ReadFile(file)) != "0"
}

// Uninstalls a file.
func UninstallFile(path string, dryRun bool) {
	if FileExists(path) {
		if dryRun {
			log.Info().Msgf("Would remove file %s", path)
		} else {
			log.Info().Msgf("Removing file %s", path)
			if err := os.Remove(path); err != nil {
				log.Info().Err(err).Msgf("Failed to remove file %s", path)
			}
		}
	}
}

// Extracts a tar.gz file.
func ExtractTarGz(tarballPath string, dstPath string) error {
	reader, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path, err := filepath.Abs(filepath.Join(dstPath, header.Name))
		if err != nil {
			return err
		}
		if !strings.HasPrefix(path, dstPath) {
			log.Warn().Msgf("Skipping extraction of %s in %s file as is resolves outside the target path",
				header.Name, tarballPath)
			continue
		}

		info := header.FileInfo()
		if info.IsDir() {
			log.Debug().Msgf("Creating folder %s", path)
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		log.Debug().Msgf("Extracting file %s", path)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err = io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}
