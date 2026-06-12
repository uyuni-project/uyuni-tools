// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpg

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type ListGpgKeysResponse struct {
	KeyType     int64    `mapstructure:"keyType"`
	KeySize     int64    `mapstructure:"keySize"`
	Fingerprint string   `mapstructure:"fingerprint"`
	Names       []string `mapstructure:"names"`
}

func gpgKeyList(client *api.APIClient) error {
	response, err := api.GetChecked[[]ListGpgKeysResponse](
		client,
		"admin/gpg/listGpgKeys",
		"admin.gpg.list_gpg_keys",
	)

	if err != nil {
		return utils.Errorf(err, L("error listing GPG keys"))
	}

	if !response.Success {
		return fmt.Errorf(L("failed to list GPG keys: %s"), response.Message)
	}

	if len(response.Result) == 0 {
		fmt.Println("No GPG keys stored.")
	}

	for keyIdx, key := range response.Result {
		fmt.Printf("[%d]\tFingerprint:\t%s\n", keyIdx, key.Fingerprint)

		typeName := "unknown"

		switch key.KeyType {
		case 1:
			typeName = "rsa"
		case 16:
			typeName = "elgamal"
		case 17:
			typeName = "dsa"
		case 18:
			typeName = "ecdh"
		case 19:
			typeName = "ecdsa"
		case 22:
			typeName = "eddsa"
		case 25:
			typeName = "x25519"
		}

		fmt.Printf("\tKey type:\t%s%d\n", typeName, key.KeySize)

		for nameIdx, name := range key.Names {
			fmt.Printf("\t[%d]\t%s\n", nameIdx, name)
		}
		fmt.Printf("\n")
	}

	return nil
}

func runGpgKeyList(_ *types.GlobalFlags, flags *apiFlags, _ *cobra.Command, _ []string) error {
	log.Debug().Msgf("Requesting GPG keys from the server...")
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return gpgKeyList(client)
}
