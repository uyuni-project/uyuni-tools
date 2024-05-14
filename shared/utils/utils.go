// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"golang.org/x/term"
)

const prompt_end = ": "

var prodVersionArchRegex = regexp.MustCompile(`suse\/manager\/.*:`)
var imageValid = regexp.MustCompile("^((?:[^:/]+(?::[0-9]+)?/)?[^:]+)(?::([^:]+))?$")

// InspectScriptFilename is the inspect script basename.
var InspectScriptFilename = "inspect.sh"

var inspectValues = []types.InspectData{
	types.NewInspectData("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3 || true"),
	types.NewInspectData("suse_manager_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4 || true"),
	types.NewInspectData("architecture", "lscpu | grep Architecture | awk '{print $2}' || true"),
	types.NewInspectData("fqdn", "cat /etc/rhn/rhn.conf 2>/dev/null | grep 'java.hostname' | cut -d' ' -f3 || true"),
	types.NewInspectData("image_pg_version", "rpm -qa --qf '%{VERSION}\\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1 || true"),
	types.NewInspectData("current_pg_version", "(test -e /var/lib/pgsql/data/PG_VERSION && cat /var/lib/pgsql/data/PG_VERSION) || true"),
	types.NewInspectData("registration_info", "env LC_ALL=C LC_MESSAGES=C LANG=C transactional-update --quiet register --status 2>/dev/null || true"),
	types.NewInspectData("scc_username", "cat /etc/zypp/credentials.d/SCCcredentials 2>&1 /dev/null | grep username | cut -d= -f2 || true"),
	types.NewInspectData("scc_password", "cat /etc/zypp/credentials.d/SCCcredentials 2>&1 /dev/null | grep password | cut -d= -f2 || true"),
}

// InspectOutputFile represents the directory and the basename where the inspect values are stored.
var InspectOutputFile = types.InspectFile{
	Directory: "/var/lib/uyuni-tools",
	Basename:  "data",
}

func checkValueSize(value string, min int, max int) bool {
	if min == 0 && max == 0 {
		return true
	}

	if len(value) < min {
		fmt.Printf(NL("Has to be more than %d character long", "Has to be more than %d characters long", min), min)
		return false
	}
	if len(value) > max {
		fmt.Printf(NL("Has to be less than %d character long", "Has to be less than %d characters long", max), max)
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
			log.Fatal().Err(err).Msgf(L("Failed to read password"))
		}
		tmpValue := strings.TrimSpace(string(bytePassword))
		r := regexp.MustCompile(`^[^\t ]+$`)
		validChars := r.MatchString(tmpValue)
		if !validChars {
			fmt.Printf(L("Cannot contain spaces or tabs"))
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
			log.Fatal().Err(err).Msgf(L("Failed to read input"))
		}
		tmpValue := strings.TrimSpace(newValue)
		if checkValueSize(tmpValue, min, max) && (checker == nil || checker(tmpValue)) {
			*value = tmpValue
		}
		fmt.Println()
		if *value == "" {
			fmt.Println(L("A value is required"))
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
		return false, nil
	}
}

// ComputeImage assembles the container image from its name and tag.
func ComputeImage(name string, tag string, appendToName ...string) (string, error) {
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
		log.Debug().Msgf("Computed image name is %s", imageName)
		return imageName, nil
	}
	imageName := submatches[1] + strings.Join(appendToName, ``) + `:` + submatches[2]
	imageName = strings.ToLower(imageName) // podman does not accept repo in upper case
	log.Debug().Msgf("Computed image name is %s", imageName)
	return imageName, nil
}

// ComputePTF returns a PTF or Test image from registry.suse.com.
func ComputePTF(user string, ptfId string, fullImage string, suffix string) (string, error) {
	prefix := fmt.Sprintf("registry.suse.com/a/%s/%s/", user, ptfId)
	submatches := prodVersionArchRegex.FindStringSubmatch(fullImage)
	if submatches == nil || len(submatches) > 1 {
		return "", fmt.Errorf(L("invalid image name: %s"), fullImage)
	}
	tag := fmt.Sprintf("latest-%s-%s", suffix, ptfId)
	return prefix + submatches[0] + tag, nil
}

// Get the timezone set on the machine running the tool.
func GetLocalTimezone() string {
	out, err := RunCmdOutput(zerolog.DebugLevel, "timedatectl", "show", "--value", "-p", "Timezone")
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to run %s"), "timedatectl show --value -p Timezone")
	}
	return string(out)
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

// Check if a given path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if !os.IsNotExist(err) {
		log.Fatal().Err(err).Msgf(L("Failed to get %s file informations"), path)
	}
	return false
}

// Returns the content of a file and exit if there was an error.
func ReadFile(file string) []byte {
	out, err := os.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to read file %s"), file)
	}
	return out
}

// Get the value of a file containing a boolean.
// This is handy for files from the kernel API.
func GetFileBoolean(file string) bool {
	return strings.TrimSpace(string(ReadFile(file))) != "0"
}

// Uninstalls a file.
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
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return err
	}

	return nil
}

// ReadInspectData returns a map with the values inspected by an image and deploy.
func ReadInspectData(scriptDir string, prefix ...string) (map[string]string, error) {
	path := filepath.Join(scriptDir, "data")
	log.Debug().Msgf("Trying to read %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, Errorf(err, L("cannot parse file %s"), path)
	}

	inspectResult := make(map[string]string)

	viper.SetConfigType("env")
	if err := viper.MergeConfig(bytes.NewBuffer(data)); err != nil {
		return map[string]string{}, Errorf(err, L("cannot read config"))
	}

	for _, v := range inspectValues {
		if len(viper.GetString(v.Variable)) > 0 {
			index := v.Variable
			/* Just the first value of prefix is used.
			 * This slice is just to allow an empty argument
			 */
			if len(prefix) >= 1 {
				index = prefix[0] + v.Variable
			}
			inspectResult[index] = viper.GetString(v.Variable)
		}
	}
	return inspectResult, nil
}

// InspectHost check values on a host machine.
func InspectHost() (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, Errorf(err, L("failed to create temporary directory"))
	}

	if err := GenerateInspectHostScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	if err := RunCmdStdMapping(zerolog.DebugLevel, scriptDir+"/inspect.sh"); err != nil {
		return map[string]string{}, Errorf(err, L("failed to run inspect script in host system"))
	}

	inspectResult, err := ReadInspectData(scriptDir, "host_")
	if err != nil {
		return map[string]string{}, Errorf(err, L("cannot inspect host data"))
	}

	return inspectResult, err
}

// GenerateInspectContainerScript create the host inspect script.
func GenerateInspectHostScript(scriptDir string) error {
	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: scriptDir + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return Errorf(err, L("failed to generate inspect script"))
	}
	return nil
}

// GenerateInspectContainerScript create the container inspect script.
func GenerateInspectContainerScript(scriptDir string) error {
	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: InspectOutputFile.Directory + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return Errorf(err, L("failed to generate inspect script"))
	}
	return nil
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
	appended := fmt.Sprintf(message, args...) + ": " + err.Error()
	return errors.New(appended)
}
