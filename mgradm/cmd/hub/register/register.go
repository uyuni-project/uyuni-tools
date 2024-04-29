// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package register

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type configFlags struct {
	Backend           string
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
}

// NewCommand command for registering peripheral server to hub.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: L("Register"),
		Long:  L("Register this peripheral server to Hub API"),
		Args:  cobra.MaximumNArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var flags configFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, register)
		},
	}
	registerCmd.SetUsageTemplate(registerCmd.UsageTemplate())

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(registerCmd)
	}

	if err := api.AddAPIFlags(registerCmd, false); err != nil {
		return nil
	}

	return registerCmd
}

func register(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	config, err := getRhnConfig(cnx)
	if err != nil {
		return err
	}
	err = registerToHub(config, &flags.ConnectionDetails)
	return err
}

func getRhnConfig(cnx *shared.Connection) (map[string]string, error) {
	out, err := cnx.Exec("/bin/cat", "/etc/rhn/rhn.conf")
	if err != nil {
		return nil, err
	}
	config := make(map[string]string)

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		log.Trace().Msgf("Config: %s", line)

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(L("invalid line format: %s"), line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	return config, nil
}

func registerToHub(config map[string]string, cnxDetails *api.ConnectionDetails) error {
	for _, key := range []string{"java.hostname", "report_db_name", "report_db_port", "report_db_user", "report_db_password"} {
		if _, ok := config[key]; !ok {
			return fmt.Errorf(L("mandatory %s entry missing in config"), key)
		}
	}
	log.Info().Msgf(L("Hub API server: %s"), cnxDetails.Server)
	client, err := api.Init(cnxDetails)
	if err != nil {
		return utils.Errorf(err, L("failed to connect to the Hub server"))
	}
	data := map[string]interface{}{
		"fqdn": config["java.hostname"],
	}

	ret, err := api.Post[int](client, "system/registerPeripheralServer", data)
	if err != nil {
		return utils.Errorf(err, L("failed to register this peripheral server"))
	}
	if !ret.Success {
		return fmt.Errorf(L("failed to register this peripheral server: %s"), ret.Message)
	}
	id := ret.Result

	data = map[string]interface{}{
		"sid":              id,
		"reportDbName":     config["report_db_name"],
		"reportDbHost":     config["java.hostname"],
		"reportDbPort":     config["report_db_port"],
		"reportDbUser":     config["report_db_user"],
		"reportDbPassword": config["report_db_password"],
	}
	ret, err = api.Post[int](client, "system/updatePeripheralServerInfo", data)
	if err != nil {
		return utils.Errorf(err, L("failed to update peripheral server info"))
	}

	if !ret.Success {
		return fmt.Errorf(L("failed to update peripheral server info: %s"), ret.Message)
	}
	log.Info().Msgf(L("Registered peripheral server: %[1]s, ID: %[2]d"), config["java.hostname"], id)
	return nil
}
