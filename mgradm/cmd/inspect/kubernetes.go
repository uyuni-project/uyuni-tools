// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package inspect

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kuberneteInspect(
	_ *types.GlobalFlags,
	flags *inspectFlags,
	_ *cobra.Command,
	_ []string,
) error {
	serverImage, err := utils.ComputeImage("", utils.DefaultTag, flags.Image)
	if err != nil && len(serverImage) > 0 {
		return utils.Errorf(err, L("failed to determine image"))
	}

	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	if len(serverImage) <= 0 {
		log.Debug().Msg("Use deployed image")

		serverImage, err = kubernetes.GetRunningImage("uyuni")
		if err != nil {
			return errors.New(L("failed to find the image of the currently running server container: %s"))
		}
	}

	namespace, err := cnx.GetNamespace("")
	if err != nil {
		return utils.Errorf(err, L("failed retrieving namespace"))
	}

	// Get the SCC credentials secret if existing
	pullSecret, err := kubernetes.GetSCCSecret(namespace, &types.SCCCredentials{}, kubernetes.ServerApp)
	if err != nil {
		return err
	}

	inspectResult, err := kubernetes.InspectServer(namespace, serverImage, flags.Image.PullPolicy, pullSecret)
	if err != nil {
		return utils.Errorf(err, L("inspect command failed"))
	}

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return utils.Errorf(err, L("cannot print inspect result"))
	}

	outputString := "\n" + string(prettyInspectOutput)
	log.Info().Msgf(outputString)

	return nil
}
