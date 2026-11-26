// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// distributionsPath is where the server extracts the distributions.
const distributionsPath = "/srv/www/distributions"

type distroUploadResult struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Registered string `json:"registered"`
	detection  bool
}

// parseUploadResult handles both the legacy integer result and the structured one.
func parseUploadResult(raw json.RawMessage) (*distroUploadResult, error) {
	var number float64
	if err := json.Unmarshal(raw, &number); err == nil {
		if int(number) != 1 {
			return nil, errors.New(L("unable to upload distro, server returned an error"))
		}
		return &distroUploadResult{}, nil
	}

	var result distroUploadResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, utils.Errorf(err, L("unable to parse the distro upload result"))
	}
	result.detection = true
	return &result, nil
}

// defaultDistroName returns the name if set, or the filename without .iso.
func defaultDistroName(filename string, name string) string {
	if name != "" {
		return strings.TrimSpace(name)
	}
	distroName := strings.TrimSpace(filename)
	if strings.HasSuffix(strings.ToLower(distroName), ".iso") {
		distroName = strings.TrimSpace(distroName[:len(distroName)-len(".iso")])
	}
	return distroName
}

// validateDistroName checks the name with the same rules as the server.
func validateDistroName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New(L("the distribution name must not be empty"))
	}
	if strings.ContainsAny(trimmed, "/\\") {
		return fmt.Errorf(L("invalid distribution name %s: it must not contain path separators"), name)
	}
	if strings.HasPrefix(trimmed, ".") {
		return fmt.Errorf(L("invalid distribution name %s: it must not start with a dot"), name)
	}
	return nil
}

// newDistroUploadBody creates a multipart body streaming the distro content.
// `size` is the size of the distro in bytes.
//
// Returns the body, its total length and its content type.
func newDistroUploadBody(filename string, name string, distro io.Reader, size int64) (
	io.Reader, int64, string, error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	if err := writer.WriteField("filename", filename); err != nil {
		return nil, 0, "", err
	}
	if name != "" {
		if err := writer.WriteField("name", name); err != nil {
			return nil, 0, "", err
		}
	}
	if _, err := writer.CreateFormFile("distro", filename); err != nil {
		return nil, 0, "", err
	}
	head := bytes.Clone(buf.Bytes())

	buf.Reset()
	if err := writer.Close(); err != nil {
		return nil, 0, "", err
	}
	tail := buf.Bytes()

	length := int64(len(head)) + size + int64(len(tail))
	body := io.MultiReader(bytes.NewReader(head), distro, bytes.NewReader(tail))
	return body, length, writer.FormDataContentType(), nil
}

// progressReader logs the progress of reading a stream of known size.
type progressReader struct {
	reader      io.Reader
	size        int64
	read        int64
	lastPercent int64
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.read += int64(n)
	// Guard against empty files to avoid dividing by zero
	if r.size > 0 {
		percent := r.read * 100 / r.size
		if percent/10 > r.lastPercent/10 {
			r.lastPercent = percent
			log.Info().Msgf(L("Uploaded %d%% (%d / %d MiB)"), percent, r.read>>20, r.size>>20)
		}
	}
	if errors.Is(err, io.EOF) {
		log.Info().Msg(L("Upload finished, waiting for the server to extract the distro..."))
	}
	return n, err
}

func distroUpload(client *api.APIClient, filename string, name string, distro io.Reader, size int64) error {
	progress := &progressReader{reader: distro, size: size}
	body, length, contentType, err := newDistroUploadBody(filename, name, progress, size)
	if err != nil {
		return utils.Errorf(err, L("error creating distro upload request"))
	}

	response, err := api.PostRawChecked[json.RawMessage](client, "admin/distro/uploadDistro", body, length, contentType)
	if err != nil {
		return utils.Errorf(err, L("error uploading distro"))
	}

	if !response.Success {
		return fmt.Errorf(L("failed to upload distro: %s"), response.Message)
	}

	result, err := parseUploadResult(response.Result)
	if err != nil {
		return err
	}

	treePath := result.Path
	if treePath == "" {
		treePath = path.Join(distributionsPath, defaultDistroName(filename, name))
	}
	fmt.Printf(L("Distro successfully uploaded to %s")+"\n", treePath)

	switch {
	case result.Registered != "":
		fmt.Printf(L("Distribution %s successfully registered")+"\n", result.Registered)
	case result.detection:
		fmt.Println(L("Distribution not registered automatically, declare it in the Web UI"))
	default:
		fmt.Println(L("Declare the distribution in the Web UI if needed"))
	}
	return nil
}

// parseURL returns the parsed source if it is an HTTP or HTTPS URL, nil otherwise.
func parseURL(source string) *url.URL {
	parsedURL, err := url.Parse(source)
	if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil
	}
	return parsedURL
}

// getFilenameFromSource returns the ISO filename from a local path or URL, or an empty string if there is none.
func getFilenameFromSource(source string) string {
	var filename string
	if parsedURL := parseURL(source); parsedURL != nil {
		filename = path.Base(parsedURL.Path)
	} else {
		filename = filepath.Base(source)
	}
	filename = strings.TrimSpace(filename)
	if filename == "." || filename == "/" || filename == string(filepath.Separator) {
		return ""
	}
	return filename
}

// openDistro opens the distro ISO from a local file or a URL for streaming.
//
// Returns the reader to close after use, the size of the ISO and its filename.
func openDistro(source string) (io.ReadCloser, int64, string, error) {
	filename := getFilenameFromSource(source)
	if filename == "" {
		return nil, 0, "", fmt.Errorf(L("unable to determine distro ISO filename from %s"), source)
	}

	if parseURL(source) == nil {
		log.Debug().Msgf("Reading distro ISO from file %s", source)
		info, err := os.Stat(source)
		if err != nil {
			return nil, 0, "", utils.Errorf(err, L("failed to read distro ISO file %s"), source)
		}
		if !info.Mode().IsRegular() {
			return nil, 0, "", fmt.Errorf(L("%s is not a regular file"), source)
		}
		file, err := os.Open(source)
		if err != nil {
			return nil, 0, "", utils.Errorf(err, L("failed to read distro ISO file %s"), source)
		}
		return file, info.Size(), filename, nil
	}

	log.Debug().Msgf("Downloading distro ISO from %s", source)
	//nolint:gosec // Downloading from the user-provided URL is the purpose of this function
	resp, err := http.Get(source)
	if err != nil {
		return nil, 0, "", utils.Errorf(err, L("failed to download distro ISO from %s"), source)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, "", fmt.Errorf(L("failed to download distro ISO from %[1]s: bad status %[2]s"),
			source, resp.Status)
	}
	if resp.ContentLength < 0 {
		resp.Body.Close()
		return nil, 0, "", fmt.Errorf(
			L("the server at %s did not report the ISO size, download it first and upload the local file"), source)
	}
	return resp.Body, resp.ContentLength, filename, nil
}

func runDistroUpload(_ *types.GlobalFlags, flags *distroUploadFlags, _ *cobra.Command, args []string) error {
	distro, size, filename, err := openDistro(args[0])
	if err != nil {
		return err
	}
	defer distro.Close()

	if err := validateDistroName(defaultDistroName(filename, flags.Name)); err != nil {
		return err
	}

	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	log.Debug().Msgf("Uploading ISO...")
	return distroUpload(client, filename, flags.Name, distro, size)
}
