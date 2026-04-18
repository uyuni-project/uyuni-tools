// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package knownhost

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func removeKnownHost(client *api.APIClient, hostname string, port string) error {
	params := url.Values{
		"hostname": {hostname},
		"port":     {port},
	}
	path := "admin/ssh/removeKnownHost?" + params.Encode()

	var data map[string]interface{}
	res, err := api.Post[interface{}](client, path, data)
	if err != nil {
		return utils.Errorf(err, L("error in query '%s'"), path)
	}

	if !res.Success {
		return fmt.Errorf(L("failed to remove the host: %s"), res.Message)
	}

	success, ok := res.Result.(float64)
	if !ok {
		return fmt.Errorf(L("unexpected server response '%v'"), res.Result)
	}

	if int(success) == 1 {
		fmt.Println(L("successfully removed host"))
	} else {
		fmt.Println(L("unable to remove host, server returned an error"))
	}

	return nil
}

func runRemoveKnownHost(_ *types.GlobalFlags, flags *apiFlags, _ *cobra.Command, args []string) error {
	hostname := args[0]
	// Default to port 22 if none is provided
	port := "22"
	if len(args) == 2 {
		port = args[1]
	}

	log.Debug().Msgf("Removing host %s:%s", hostname, port)
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return removeKnownHost(client, hostname, port)
}
