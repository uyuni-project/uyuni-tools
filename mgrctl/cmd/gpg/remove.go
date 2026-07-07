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

func gpgKeyRemove(client *api.APIClient, fingerprint string) error {
	response, err := api.PostChecked[float64](
		client,
		"admin/gpg/removeGpgKey",
		"admin.gpg.remove_gpg_key",
		map[string]interface{}{
			"fingerprint": fingerprint,
		},
	)

	if err != nil {
		return utils.Errorf(err, L("error removing GPG key"))
	}

	if !response.Success {
		return fmt.Errorf(L("failed to remove GPG key: %s"), response.Message)
	}

	if int(response.Result) == 1 {
		fmt.Println(L("GPG key successfully removed"))
	} else {
		fmt.Println(L("unable to remove GPG key, server returned an error"))
	}

	return nil
}

func runGpgKeyRemove(_ *types.GlobalFlags, flags *apiFlags, _ *cobra.Command, args []string) error {
	fingerprint := args[0]

	log.Debug().Msgf("Requesting GPG key deletion from the server...")

	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return gpgKeyRemove(client, fingerprint)
}
