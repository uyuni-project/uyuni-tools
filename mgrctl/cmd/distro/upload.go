// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"bytes"
	"fmt"
	"mime/multipart"
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

func distroUpload(client *api.APIClient, filename string, distro []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("filename", filename); err != nil {
		return utils.Errorf(err, L("error creating distro upload request"))
	}
	part, err := writer.CreateFormFile("distro", filename)
	if err != nil {
		return utils.Errorf(err, L("error creating distro upload request"))
	}
	if _, err = part.Write(distro); err != nil {
		return utils.Errorf(err, L("error creating distro upload request"))
	}
	if err = writer.Close(); err != nil {
		return utils.Errorf(err, L("error creating distro upload request"))
	}

	response, err := api.PostRaw[float64](client, "admin/distro/uploadDistro", body, writer.FormDataContentType())
	if err != nil {
		return utils.Errorf(err, L("error uploading distro"))
	}

	if !response.Success {
		return fmt.Errorf(L("failed to upload distro: %s"), response.Message)
	}

	if int(response.Result) == 1 {
		fmt.Println(L("Distro successfully uploaded"))
	} else {
		fmt.Println(L("unable to upload distro, server returned an error"))
	}

	return nil
}

func getFilenameFromSource(source string) string {
	if parsedURL, err := url.Parse(source); err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		filename := path.Base(parsedURL.Path)
		if filename != "." && filename != "/" {
			return filename
		}
		return ""
	}
	return filepath.Base(source)
}

func readDistro(source string) ([]byte, string, error) {
	var data []byte
	var err error

	filename := strings.TrimSpace(getFilenameFromSource(source))
	if filename == "" || filename == "." || filename == "/" {
		return nil, "", fmt.Errorf(L("unable to determine distro ISO filename from %s"), source)
	}

	if _, err = os.Stat(source); err == nil {
		log.Debug().Msgf("Reading distro ISO from file %s", source)
		data, err = os.ReadFile(source)
		if err != nil {
			return nil, "", utils.Errorf(err, L("failed to read distro ISO file %s"), source)
		}
	} else {
		log.Debug().Msgf("Downloading distro ISO from %s", source)
		data, err = utils.GetURLBody(source)
		if err != nil {
			return nil, "", utils.Errorf(err, L("failed to download distro ISO from %s"), source)
		}
	}

	return data, filename, nil
}

func runDistroUpload(_ *types.GlobalFlags, flags *apiFlags, _ *cobra.Command, args []string) error {
	source := args[0]
	distro, filename, err := readDistro(source)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Uploading ISO...")
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return distroUpload(client, filename, distro)
}
