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
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func runPost(globalFlags *types.GlobalFlags, flags *apiFlags, cmd *cobra.Command, args []string) error {
	log.Debug().Msgf("Running POST command %s", args[0])
	client, err := api.Init(&flags.ConnectionDetails)

	if err != nil {
		return fmt.Errorf(L("unable to login to the server: %s"), err)
	}

	path := args[0]
	options := args[1:]

	var data map[string]interface{}

	if len(options) > 1 {
		log.Debug().Msg("Multiple options specified, assuming non JSON data")
		data = map[string]interface{}{}
		for _, o := range options {
			s := strings.SplitN(o, "=", 2)
			data[s[0]] = s[1]
		}
	} else {
		if err := json.NewDecoder(strings.NewReader(args[1])).Decode(&data); err != nil {
			log.Debug().Msg("Failed to decode parameters as JSON, assuming key=value pairs")
		}
	}

	res, err := api.Post[interface{}](client, path, data)
	if err != nil {
		return fmt.Errorf(L("error in query %s: %s"), path, err)
	}

	if !res.Success {
		log.Error().Msg(res.Message)
	}
	out, err := json.MarshalIndent(res.Result, "", "  ")
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Print(string(out))

	return nil
}
