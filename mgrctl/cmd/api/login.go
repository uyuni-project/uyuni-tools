// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func runLogin(globalFlags *types.GlobalFlags, flags *apiFlags, cmd *cobra.Command, args []string) error {
	log.Debug().Msg("Running login command")

	utils.AskIfMissing(&flags.User, cmd.Flag("api-user").Usage, 0, 0, nil)
	utils.AskPasswordIfMissing(&flags.Password, cmd.Flag("api-password").Usage, 0, 0)
	// ToDO add FQDN checker from rebase
	utils.AskIfMissing(&flags.Server, cmd.Flag("api-server").Usage, 0, 0, nil)

	return api.StoreLoginCreds(cmd.Context(), &flags.ConnectionDetails)
}

func runLogout(globalFlags *types.GlobalFlags, flags *apiFlags, cmd *cobra.Command, args []string) error {
	log.Debug().Msg("Running logout command")

	return api.RemoveLoginCreds(cmd.Context())
}
