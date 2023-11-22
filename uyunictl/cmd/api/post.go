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

func runPost(globalFlags *types.GlobalFlags, flags *api.ConnectionDetails, cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Running POST command %s", args[0])
	client, err := api.Init(flags)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to login to the server")
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

	res, err := client.Post(path, data)
	if err != nil {
		log.Error().Err(err).Msgf("Error in query %s", path)
		return
	}

	if !res["success"].(bool) {
		log.Error().Msg(res["message"].(string))
	}
	out, err := json.MarshalIndent(res["result"], "", "  ")
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Print(out)
}
