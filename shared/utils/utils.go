// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"golang.org/x/term"
)

const promptEnd = ": "

var prodVersionArchRegex = regexp.MustCompile(`suse\/(?:multi-linux-)?manager\/.*:`)
var imageValid = regexp.MustCompile("^((?:[^:/]+(?::[0-9]+)?/)?[^:]+)(?::([^:]+))?$")

// Taken from https://github.com/go-playground/validator/blob/2e1df48/regexes.go#L58
var fqdnValid = regexp.MustCompile(
	`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?` +
		`(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$`,
)

// InspectResult holds the results of the inspection scripts.
type InspectResult struct {
	CommonInspectData `mapstructure:",squash"`
	DBInspectData     `mapstructure:",squash"`
	Timezone          string
	HasHubXmlrpcAPI   bool `mapstructure:"has_hubxmlrpc"`
	Debug             bool `mapstructure:"debug"`
}

func checkValueSize(value string, minValue int, maxValue int) bool {
	if minValue == 0 && maxValue == 0 {
		return true
	}

	if len(value) < minValue {
		fmt.Printf(NL("Has to be more than %d character long", "Has to be more than %d characters long", minValue), minValue)
		return false
	}
	if len(value) > maxValue {
		fmt.Printf(NL("Has to be less than %d character long", "Has to be less than %d characters long", maxValue), maxValue)
		return false
	}
	return true
}

// CheckValidPassword performs check to a given password.
func CheckValidPassword(value *string, prompt string, minValue int, maxValue int) string {
	fmt.Print(prompt + promptEnd)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Error().Err(err).Msg(L("Failed to read password"))
		return ""
	}
	tmpValue := strings.TrimSpace(string(bytePassword))

	if tmpValue == "" {
		fmt.Println("A value is required")
		return ""
	}

	r := regexp.MustCompile(`[\t ]`)
	invalidChars := r.MatchString(tmpValue)

	if invalidChars {
		fmt.Println(L("Cannot contain spaces or tabs"))
		return ""
	}

	if !checkValueSize(tmpValue, minValue, maxValue) {
		fmt.Println()
		return ""
	}
	fmt.Println()
	*value = tmpValue
	return *value
}

// AskPasswordIfMissing asks for password if missing.
// Don't perform any check if min and max are set to 0.
func AskPasswordIfMissing(value *string, prompt string, minValue int, maxValue int) {
	if *value == "" && !term.IsTerminal(int(os.Stdin.Fd())) {
		log.Warn().Msg(L("not an interactive device, not asking for missing value"))
		return
	}

	for *value == "" {
		firstRound := CheckValidPassword(value, prompt, minValue, maxValue)
		if firstRound == "" {
			continue
		}
		secondRound := CheckValidPassword(value, L("Confirm the password"), minValue, maxValue)
		if secondRound != firstRound {
			fmt.Println(L("Two different passwords have been provided"))
			*value = ""
		} else {
			*value = secondRound
		}
	}
}

// AskPasswordIfMissingOnce asks for password if missing only once
// Don't perform any check if min and max are set to 0.
func AskPasswordIfMissingOnce(value *string, prompt string, minValue int, maxValue int) {
	if *value == "" && !term.IsTerminal(int(os.Stdin.Fd())) {
		log.Warn().Msg(L("not an interactive device, not asking for missing value"))
		return
	}

	for *value == "" {
		*value = CheckValidPassword(value, prompt, minValue, maxValue)
	}
}

// AskIfMissing asks for a value if missing.
// Don't perform any check if minValue and maxValue are set to 0.
func AskIfMissing(value *string, prompt string, minValue int, maxValue int, checker func(string) bool) {
	if *value == "" && !term.IsTerminal(int(os.Stdin.Fd())) {
		log.Warn().Msg(L("not an interactive device, not asking for missing value"))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for *value == "" {
		fmt.Print(prompt + promptEnd)
		newValue, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err).Msg(L("failed to read input"))
		}
		tmpValue := strings.TrimSpace(newValue)
		if checkValueSize(tmpValue, minValue, maxValue) && (checker == nil || checker(tmpValue)) {
			*value = tmpValue
		}
		fmt.Println()
		if *value == "" {
			fmt.Print(L("A value is required"))
		}
	}
}

// YesNo asks a question in CLI.
func YesNo(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]?", question)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			return true, nil
		}
		if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
			return false, nil
		}
	}
}

// SplitRegistryHostAndPath splits a registry string into domain and path.
func SplitRegistryHostAndPath(registry string) (domain string, path string) {
	separator := "://"
	index := strings.Index(registry, separator)
	if index != -1 {
		registry = registry[index+len(separator):]
	}

	idx := strings.Index(registry, "/")
	if idx == -1 {
		return registry, ""
	}
	return registry[:idx], registry[idx+1:]
}

// ComputeImage assembles the container image from its name and tag.
func ComputeImage(
	globalRegistry string,
	globalTag string,
	imageFlags types.ImageFlags,
) (string, error) {
	// Compute the registry
	registry := globalRegistry
	if imageFlags.Registry.Host != "" {
		registry = imageFlags.Registry.Host
	}

	name := imageFlags.Name

	if !StartWithFQDN(name) {
		name = path.Join(registry, imageFlags.Name)
	}

	// Compute the tag
	tag := globalTag
	if imageFlags.Tag != "" {
		tag = imageFlags.Tag
	}

	submatches := imageValid.FindStringSubmatch(name)
	if submatches == nil {
		return "", fmt.Errorf(L("invalid image name: %s"), name)
	}
	if submatches[2] == `` {
		if len(tag) <= 0 {
			return name, fmt.Errorf(L("tag missing on %s"), name)
		}
		// No tag provided in the URL name, append the one passed
		imageName := fmt.Sprintf("%s:%s", name, tag)
		imageName = strings.ToLower(imageName) // podman does not accept repo in upper case
		log.Info().Msgf(L("Computed image name is %s"), imageName)
		return imageName, nil
	}
	imageName := strings.ToLower(name) // podman does not accept repo in upper case
	log.Info().Msgf(L("Computed image name is %s"), imageName)
	return imageName, nil
}

// The fullImage must contain the pattern `suse/manager/...` or `suse/multi-linux-manager/...`
// If registry has a path, then, the fullImage must start with that path.
func ComputePTF(registry string, user string, ptfID string, fullImage string, suffix string) (string, error) {
	submatches := prodVersionArchRegex.FindStringSubmatch(fullImage)
	if submatches == nil || len(submatches) != 1 {
		return "", fmt.Errorf(L("invalid image name: %s"), fullImage)
	}
	imagePath := submatches[0]

	tag := fmt.Sprintf("latest-%s-%s", suffix, ptfID)

	registryHost, registryPath := SplitRegistryHostAndPath(registry)
	// registry.suse.de is an internal registry and ptf containers here
	// are shipped in a slightly different path
	if registryHost == "registry.suse.de" {
		sep := "containerfile/"
		idx := strings.Index(registry, sep)
		if idx > 1 {
			imagePath = registry[idx+len(sep):]
		}
		return fmt.Sprintf("%s/ptf/%s/containers/a/%s%s", registryHost, ptfID, imagePath, tag), nil
	}
	if registryPath != "" && !strings.HasPrefix(imagePath, registryPath) {
		return "", fmt.Errorf(L("image path '%[1]s' does not start with registry path '%[2]s'"), imagePath, registryPath)
	}

	return fmt.Sprintf("%s/a/%s/%s/%s%s", registryHost, strings.ToLower(user), ptfID, imagePath, tag), nil
}

// GetLocalTimezone returns the timezone set on the current machine.
func GetLocalTimezone() string {
	out, err := RunCmdOutput(zerolog.DebugLevel, "timedatectl", "show", "--value", "-p", "Timezone")
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to run %s"), "timedatectl show --value -p Timezone")
	}
	return strings.TrimSpace(string(out))
}

// GetRandomBase64 generates random base64-encoded data.
func GetRandomBase64(size int) string {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		log.Fatal().Err(err).Msg(L("Failed to read random data"))
	}
	return base64.StdEncoding.EncodeToString(data)
}

// ContainsUpperCase check if string contains an uppercase character.
func ContainsUpperCase(str string) bool {
	for _, char := range str {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

// GetURLBody provide the body content of an GET HTTP request.
func GetURLBody(URL string) ([]byte, error) {
	// Download the key from the URL
	log.Debug().Msgf("Downloading %s", URL)
	resp, err := http.Get(URL)
	if err != nil {
		return nil, Errorf(err, L("error downloading from %s"), URL)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(L("bad status: %s"), resp.Status)
	}

	var buf bytes.Buffer

	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	// Extract the byte slice from the buffer
	data := buf.Bytes()
	return data, nil
}

// DownloadFile downloads from a remote path to a local file.
func DownloadFile(filepath string, URL string) (err error) {
	data, err := GetURLBody(URL)
	if err != nil {
		return err
	}

	// Writer the body to file
	log.Debug().Msgf("Saving %s to %s", URL, filepath)
	return os.WriteFile(filepath, data, 0644)
}

// maxInts compensates the absence of max on Debian 12's go version.
func maxInts(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

// CompareVersion compare the server image version and the server deployed  version.
func CompareVersion(imageVersion string, deployedVersion string) int {
	image := versionAsSlice(imageVersion)
	deployed := versionAsSlice(deployedVersion)

	maxLen := maxInts(len(image), len(deployed))
	return getPaddedVersion(image, maxLen) - getPaddedVersion(deployed, maxLen)
}

func versionAsSlice(version string) []string {
	re := regexp.MustCompile(`[^0-9]`)
	parts := strings.Split(version, ".")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = re.ReplaceAllString(part, "")
	}
	return result
}

func getPaddedVersion(version []string, size int) int {
	padded := version
	if len(version) != size {
		padded = make([]string, size)
		copy(padded, version)
		for i, part := range padded {
			if part == "" {
				padded[i] = "0"
			}
		}
	}

	result, _ := strconv.Atoi(strings.Join(padded, ""))
	return result
}

// Errorf helps providing consistent errors.
//
// Instead of fmt.Printf(L("the message for %s: %s"), value, err) use:
//
//	Errorf(err, L("the message for %s"), value)
func Errorf(err error, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	return Error(err, formattedMessage)
}

// Error helps providing consistent errors.
//
// Instead of fmt.Printf(L("the message: %s"), err) use:
//
//	Error(err, L("the message"))
func Error(err error, message string) error {
	// l10n-ignore
	return fmt.Errorf("%s: %w", message, err)
}

// JoinErrors aggregate multiple multiple errors into one.
//
// Replacement for errors.Join which is not available in go 1.19.
func JoinErrors(errs ...error) error {
	var messages []string
	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return errors.New(strings.Join(messages, "; "))
}

// GetFqdn returns and checks the FQDN of the host system.
func GetFqdn(args []string) (string, error) {
	var fqdn string
	if len(args) == 1 {
		fqdn = args[0]
	} else {
		out, err := RunCmdOutput(zerolog.DebugLevel, "hostname", "-f")
		if err != nil {
			return "", Error(err, L("failed to compute server FQDN"))
		}
		fqdn = strings.TrimSpace(string(out))
	}
	if err := IsValidFQDN(fqdn); err != nil {
		return "", err
	}

	return fqdn, nil
}

// IsValidFQDN returns an error if the argument is not a valid FQDN.
func IsValidFQDN(fqdn string) error {
	if !IsWellFormedFQDN(fqdn) {
		return fmt.Errorf(L("%s is not a valid FQDN"), fqdn)
	}
	_, err := net.LookupHost(fqdn)
	if err != nil {
		return Errorf(err, L("cannot resolve %s"), fqdn)
	}
	return nil
}

// IsWellFormedFQDN returns an false if the argument is not a well formed FQDN.
func IsWellFormedFQDN(fqdn string) bool {
	return fqdnValid.MatchString(fqdn)
}

// StartWithFQDN returns true if the argument start with a well formed FQDN.
func StartWithFQDN(url string) bool {
	fqdn, _ := SplitRegistryHostAndPath(url)
	return IsWellFormedFQDN(fqdn)
}

// CommandExists checks if cmd exists in $PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// SaveBinaryData saves binary data to a file.
func SaveBinaryData(filename string, data []int8) error {
	// Need to convert the array of signed ints to unsigned/byte
	byteArray := make([]byte, len(data))
	for i, v := range data {
		byteArray[i] = byte(v)
	}
	file, err := os.Create(filename)
	if err != nil {
		return Errorf(err, L("error creating file %s"), filename)
	}
	defer file.Close()
	_, err = file.Write(byteArray)
	if err != nil {
		return Errorf(err, L("error writing file %s"), filename)
	}
	return nil
}

// CreateChecksum creates sha256 checksum of provided file.
// Uses system `sha256sum` binary to avoid pulling crypto dependencies.
func CreateChecksum(file string) error {
	outputFile := file + ".sha256sum"

	output, err := NewRunner("sha256sum", file).Exec()
	if err != nil {
		return Errorf(err, L("Failed to calculate checksum of the file %s"), file)
	}
	// We want only checksum, drop the filepath
	output = bytes.Split(output, []byte(" "))[0]
	if err := os.WriteFile(outputFile, output, 0622); err != nil {
		return Errorf(err, L("Failed to write checksum of the file %[1]s to the %[2]s"), file, outputFile)
	}
	return nil
}

// ValidateChecksum checks integrity of the file by checking against stored checksum
// Uses system `sha256sum` binary to avoid pulling crypt dependencies.
func ValidateChecksum(file string) error {
	checksum, err := NewRunner("sha256sum", file).Exec()
	if err != nil {
		return Errorf(err, L("Failed to calculate checksum of the file %s"), file)
	}
	// We want only checksum, drop the filepath
	checksum = bytes.Split(checksum, []byte(" "))[0]

	output, err := os.ReadFile(file + ".sha256sum")
	if err != nil {
		return Errorf(err, L("Failed to read checksum of the file %[1]s"), file)
	}
	// Split by space to work with older backups
	if !bytes.Equal(checksum, bytes.Split(output, []byte(" "))[0]) {
		return fmt.Errorf(L("Checksum of %s does not match"), file)
	}
	return nil
}
