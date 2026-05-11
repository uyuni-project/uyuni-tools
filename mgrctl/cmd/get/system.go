// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	apitypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"gopkg.in/yaml.v2"
)

func newSystemCommand(globalFlags *types.GlobalFlags, parentFlags *getFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "system [name/id]",
		Aliases: []string{"systems"},
		Short:   L("List or get details for registered systems"),
		RunE: func(_ *cobra.Command, args []string) error {
			return runGetSystem(globalFlags, parentFlags, args)
		},
	}
}

func runGetSystem(_ *types.GlobalFlags, flags *getFlags, _ []string) error {
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil && (client.Details.User != "" || client.Details.InSession) {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	res, err := api.Get[[]apitypes.System](client, "system/listSystems")
	if err != nil {
		return utils.Errorf(err, L("failed to fetch systems from API"))
	}

	systems := res.Result

	switch flags.OutputFormat {
	case "json":
		out, err := json.MarshalIndent(systems, "", "  ")
		if err != nil {
			return utils.Errorf(err, L("failed to marshal JSON"))
		}
		fmt.Println(string(out))
	case "yaml":
		out, err := yaml.Marshal(systems)
		if err != nil {
			return utils.Errorf(err, L("failed to marshal YAML"))
		}
		fmt.Println(string(out))
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tLAST_CHECKIN\tCREATED")
		for _, sys := range systems {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				sys.ID,
				sys.Name,
				sys.LastCheckin,
				sys.Created,
			)
		}
		w.Flush()
	}

	return nil
}
