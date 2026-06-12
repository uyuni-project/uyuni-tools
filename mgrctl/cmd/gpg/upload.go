// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const armorHeader = "-----BEGIN PGP PUBLIC KEY BLOCK-----"

func gpgKeyUpload(client *api.APIClient, key string) error {
	response, err := api.PostChecked[float64](
		client,
		"admin/gpg/uploadGpgKey",
		"admin.gpg.upload_gpg_key",
		map[string]interface{}{
			"gpgKey": key,
		},
	)

	if err != nil {
		return utils.Errorf(err, L("error uploading GPG key"))
	}

	if !response.Success {
		return fmt.Errorf(L("failed to upload GPG key: %s"), response.Message)
	}

	if int(response.Result) == 1 {
		fmt.Println(L("GPG key successfully uploaded"))
	} else {
		fmt.Println(L("unable to upload GPG key, server returned an error"))
	}

	return nil
}

func readKey(source string) (string, error) {
	var data []byte
	var err error

	if _, err = os.Stat(source); err == nil {
		log.Debug().Msgf("Reading GPG key from file %s", source)
		data, err = os.ReadFile(source)
		if err != nil {
			return "", utils.Errorf(err, L("failed to read key file %s"), source)
		}
	} else {
		log.Debug().Msgf("Downloading GPG key from %s", source)
		data, err = utils.GetURLBody(source)
		if err != nil {
			return "", utils.Errorf(err, L("failed to download key from %s"), source)
		}
	}

	key := string(data)
	// Armored GPG keys start with this header.
	if !strings.Contains(key, armorHeader) {
		return "", errors.New(L("the provided key is not an armored GPG key"))
	}

	return key, nil
}

func runGpgKeyUpload(_ *types.GlobalFlags, flags *apiFlags, _ *cobra.Command, args []string) error {
	source := args[0]

	key, err := readKey(source)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Uploading GPG key...")
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return gpgKeyUpload(client, key)
}
