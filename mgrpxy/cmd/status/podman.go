// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func podmanStatus(
	globalFlags *types.GlobalFlags,
	flags *statusFlags,
	cmd *cobra.Command,
	args []string,
) error {
	var returnErr error
	services := []string{"httpd", "salt-broker", "squid", "ssh", "tftpd", "pod"}
	for _, service := range services {
		serviceName := fmt.Sprintf("uyuni-proxy-%s", service)
		if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", serviceName); err != nil {
			log.Error().Err(err).Msgf("Failed to get status of the %s service", serviceName)
			returnErr = errors.New("failed to get the status of at least one service")
		}
	}
	return returnErr
}
