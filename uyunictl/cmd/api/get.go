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

func runGet(globalFlags *types.GlobalFlags, flags *api.ConnectionDetails, cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Running GET command %s", args[0])
	client := api.Init(flags)

	if flags.User != "" && flags.Password != "" {
		if err := client.Login(flags.User, flags.Password); err != nil {
			log.Fatal().Err(err).Msg("Unable to login to the server")
		}
	}
	path := args[0]
	options := args[1:]

	res, err := client.Get(fmt.Sprintf("%s?%s", path, strings.Join(options, "&")))
	if err != nil {
		log.Error().Err(err).Msgf("Error in query %s", path)
	}

	// TODO do this only when result is JSON or TEXT. Watchout for binary data
	// Decode JSON to the string and pretty print it
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Print(string(out))
}
