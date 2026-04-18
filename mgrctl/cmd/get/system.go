// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
    "encoding/json"
    "fmt"
    "strings"
    "text/tabwriter"
    "os"

    "github.com/spf13/cobra"
    "github.com/uyuni-project/uyuni-tools/shared/api"
    . "github.com/uyuni-project/uyuni-tools/shared/l10n"
    "github.com/uyuni-project/uyuni-tools/shared/types"
    "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// System represents a system from the Uyuni API.
type System struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    LastCheckin string `json:"last_checkin"`
}

func newSystemCommand(globalFlags *types.GlobalFlags) *cobra.Command {
    var flags getFlags

    cmd := &cobra.Command{
        Use:   "system [name]",
        Short: L("List or search systems"),
        Long: L(`List all registered systems or search by name.

Examples:
  # List all systems
  mgrctl get system

  # Search systems by name
  mgrctl get system webserver

  # Output as JSON
  mgrctl get system -o json`),
        RunE: func(cmd *cobra.Command, args []string) error {
            return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runSystem)
        },
    }

    cmd.Flags().StringP("output", "o", "table", L("Output format: table, json, yaml"))

    return cmd
}

func runSystem(_ *types.GlobalFlags, flags *getFlags, _ *cobra.Command, args []string) error {
    client, err := api.Init(&flags.ConnectionDetails)
    if err != nil {
        return utils.Errorf(err, L("failed to initialize API client"))
    }

    if client.Details.User != "" || client.Details.InSession {
        if err = client.Login(); err != nil {
            return utils.Errorf(err, L("unable to login to the server"))
        }
    }

    if len(args) > 0 {
        return searchSystem(client, args[0], flags.Output)
    }
    return listSystems(client, flags.Output)
}

func listSystems(client *api.APIClient, format string) error {
    res, err := api.Get[[]System](client, "system/listSystems")
    if err != nil {
        return utils.Errorf(err, L("failed to list systems"))
    }

    return printSystems(res.Result, format)
}

func searchSystem(client *api.APIClient, name string, format string) error {
    res, err := api.Get[[]System](client, fmt.Sprintf("system/searchByName?searchTerm=%s", name))
    if err != nil {
        return utils.Errorf(err, L("failed to search systems"))
    }

    return printSystems(res.Result, format)
}

func printSystems(systems []System, format string) error {
    switch strings.ToLower(format) {
    case "json":
        out, err := json.MarshalIndent(systems, "", "  ")
        if err != nil {
            return utils.Errorf(err, L("failed to format output as JSON"))
        }
        fmt.Println(string(out))
    case "table", "":
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        fmt.Fprintln(w, "ID\tNAME\tLAST CHECKIN")
        for _, s := range systems {
            fmt.Fprintf(w, "%d\t%s\t%s\n", s.ID, s.Name, s.LastCheckin)
        }
        w.Flush()
    default:
        return fmt.Errorf(L("unsupported output format: %s"), format)
    }

    return nil
}