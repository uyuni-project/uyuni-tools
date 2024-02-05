// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func runGet(globalFlags *types.GlobalFlags, flags *apiFlags, cmd *cobra.Command, args []string) error {
	log.Debug().Msgf("Running GET command %s", args[0])
	client, err := api.Init(&flags.ConnectionDetails)

	if err != nil {
		return fmt.Errorf("unable to login to the server: %s", err)
	}
	path := args[0]
	options := args[1:]

	res, err := api.Get[interface{}](client, fmt.Sprintf("%s?%s", path, strings.Join(options, "&")))
	if err != nil {
		return fmt.Errorf("error in query %s: %s", path, err)
	}

	// TODO do this only when result is JSON or TEXT. Watchout for binary data
	// Decode JSON to the string and pretty print it
	out, err := json.MarshalIndent(res.Result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(out))

	return nil
}
