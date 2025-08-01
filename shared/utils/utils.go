// SPDX-FileCopyrightText: 2025 SUSE LLC
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

var prodVersionArchRegex = regexp.MustCompile(`suse\/manager\/.*:`)
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

// RemoveRegistryFromImage removes registry fqdn from image path.
func RemoveRegistryFromImage(imagePath string) string {
	separator := "://"
	index := strings.Index(imagePath, separator)
	if index != -1 {
		imagePath = imagePath[index+len(separator):]
	}

	parts := strings.Split(imagePath, "/")
	if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") || index != -1 {
		// first part is a registry fqdn
		parts = parts[1:]
	}
	return strings.Join(parts, "/")
}

// ComputeImage assembles the container image from its name and tag.
func ComputeImage(
	registry string,
	globalTag string,
	imageFlags types.ImageFlags,
	appendToName ...string,
) (string, error) {
	if !strings.Contains(DefaultRegistry, registry) {
		log.Info().Msgf(L("Registry %[1]s would be used instead of namespace %[2]s"), registry, DefaultRegistry)
	}
	name := imageFlags.Name
	if !strings.Contains(imageFlags.Name, registry) {
		name = path.Join(registry, RemoveRegistryFromImage(imageFlags.Name))
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
		if len(appendToName) > 0 {
			name = name + strings.Join(appendToName, ``)
		}
		// No tag provided in the URL name, append the one passed
		imageName := fmt.Sprintf("%s:%s", name, tag)
		imageName = strings.ToLower(imageName) // podman does not accept repo in upper case
		log.Info().Msgf(L("Computed image name is %s"), imageName)
		return imageName, nil
	}
	imageName := submatches[1] + strings.Join(appendToName, ``) + `:` + submatches[2]
	imageName = strings.ToLower(imageName) // podman does not accept repo in upper case
	log.Info().Msgf(L("Computed image name is %s"), imageName)
	return imageName, nil
}

// ComputePTF returns a PTF or Test image from registry.suse.com.
func ComputePTF(user string, ptfID string, fullImage string, suffix string) (string, error) {
	prefix := fmt.Sprintf("registry.suse.com/a/%s/%s/", user, ptfID)
	submatches := prodVersionArchRegex.FindStringSubmatch(fullImage)
	if submatches == nil || len(submatches) > 1 {
		return "", fmt.Errorf(L("invalid image name: %s"), fullImage)
	}
	tag := fmt.Sprintf("latest-%s-%s", suffix, ptfID)
	return prefix + submatches[0] + tag, nil
}

// GetLocalTimezone returns the timezone set on the current machine.
func GetLocalTimezone() string {
	out, err := RunCmdOutput(zerolog.DebugLevel, "timedatectl", "show", "--value", "-p", "Timezone")
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to run %s"), "timedatectl show --value -p Timezone")
	}
	return strings.TrimSpace(string(out))
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
	cleaner := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Error().Err(err).Msg(L("failed to remove temporary directory"))
		}
	}
	return tempDir, cleaner, nil
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

// CompareVersion compare the server image version and the server deployed  version.
func CompareVersion(imageVersion string, deployedVersion string) int {
	re := regexp.MustCompile(`\((.*?)\)`)
	imageVersionCleaned := strings.ReplaceAll(imageVersion, ".", "")
	imageVersionCleaned = strings.TrimSpace(imageVersionCleaned)
	imageVersionCleaned = re.ReplaceAllString(imageVersionCleaned, "")
	imageVersionInt, _ := strconv.Atoi(imageVersionCleaned)

	deployedVersionCleaned := strings.ReplaceAll(deployedVersion, ".", "")
	deployedVersionCleaned = strings.TrimSpace(deployedVersionCleaned)
	deployedVersionCleaned = re.ReplaceAllString(deployedVersionCleaned, "")
	deployedVersionInt, _ := strconv.Atoi(deployedVersionCleaned)
	return imageVersionInt - deployedVersionInt
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

	output, err := RunCmdOutput(zerolog.DebugLevel, "sha256sum", file)
	if err != nil {
		return Errorf(err, L("Failed to calculate checksum of the file %s"), file)
	}

	if err := os.WriteFile(outputFile, output, 0622); err != nil {
		return Errorf(err, L("Failed to write checksum of the file %s to the %s"), file, outputFile)
	}
	return nil
}

// ValidateChecksum checks integrity of the file by checking against stored checksum
// Uses system `sha256sum` binary to avoid pulling crypt dependencies.
func ValidateChecksum(file string) error {
	err := RunCmd("sha256sum", "--check", "--status", file+".sha256sum")
	if err != nil {
		return fmt.Errorf(L("Checksum of %s does not match"), file)
	}
	return nil
}
