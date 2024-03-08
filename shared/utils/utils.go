// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

const prompt_end = ": "

func checkValueSize(value string, min int, max int) bool {
	if min == 0 && max == 0 {
		return true
	}

	if len(value) < min {
		fmt.Printf("Has to be more than %d characters long", min)
		return false
	}
	if len(value) > max {
		fmt.Printf("Has to be less than %d characters long", max)
		return false
	}
	return true
}

// AskPasswordIfMissing asks for password if missing.
// Don't perform any check if min and max are set to 0.
func AskPasswordIfMissing(value *string, prompt string, min int, max int) {
	for *value == "" {
		fmt.Print(prompt + prompt_end)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read password")
		}
		tmpValue := strings.TrimSpace(string(bytePassword))
		r := regexp.MustCompile(`^[^\t ]+$`)
		validChars := r.MatchString(tmpValue)
		if !validChars {
			fmt.Printf("Cannot contain spaces or tabs")
		}

		if validChars && checkValueSize(tmpValue, min, max) {
			*value = tmpValue
		}
		fmt.Println()
		if *value == "" {
			fmt.Println("A value is required")
		}
	}
}

// AskIfMissing asks for a value if missing.
// Don't perform any check if min and max are set to 0.
func AskIfMissing(value *string, prompt string, min int, max int, checker func(string) bool) {
	reader := bufio.NewReader(os.Stdin)
	for *value == "" {
		fmt.Print(prompt + prompt_end)
		newValue, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read input")
		}
		tmpValue := strings.TrimSpace(newValue)
		if checkValueSize(tmpValue, min, max) && (checker == nil || checker(tmpValue)) {
			*value = tmpValue
		}
		fmt.Println()
		if *value == "" {
			fmt.Println("A value is required")
		}
	}
}

// ComputeImage assembles the container image from its name and tag.
func ComputeImage(name string, tag string, appendToName ...string) (string, error) {
	imageValid := regexp.MustCompile("^((?:[^:/]+(?::[0-9]+)?/)?[^:]+)(?::([^:]+))?$")
	submatches := imageValid.FindStringSubmatch(name)
	if submatches == nil {
		return "", fmt.Errorf("invalid image name: %s", name)
	}
	if submatches[2] == `` {
		if len(tag) <= 0 {
			return name, fmt.Errorf("tag missing on %s", name)
		}
		if len(appendToName) > 0 {
			name = name + strings.Join(appendToName, ``)
		}
		// No tag provided in the URL name, append the one passed
		imageName := fmt.Sprintf("%s:%s", name, tag)
		log.Debug().Msgf("Computed image name is %s", imageName)
		return imageName, nil
	}
	imageName := submatches[1] + strings.Join(appendToName, ``) + `:` + submatches[2]
	log.Debug().Msgf("Computed image name is %s", imageName)
	return imageName, nil
}

// Get the timezone set on the machine running the tool.
func GetLocalTimezone() string {
	out, err := RunCmdOutput(zerolog.DebugLevel, "timedatectl", "show", "--value", "-p", "Timezone")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run timedatectl show --value -p Timezone")
	}
	return string(out)
}

// Check if a given path exists.
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
// This is handy for files from the kernel API.
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

// GetRandomBase64 generates random base64-encoded data.
func GetRandomBase64(size int) string {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		log.Fatal().Err(err).Msg("Failed to read random data")
	}
	return base64.StdEncoding.EncodeToString(data)
}
