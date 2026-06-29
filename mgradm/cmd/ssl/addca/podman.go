// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package addca

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_podman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func addCAForPodman(_ *types.GlobalFlags, flags *addCAFlags, cmd *cobra.Command, args []string) error {
	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}

	if !flags.SSL.Ca.IsThirdParty() {
		utils.AskPasswordIfMissing(&flags.SSL.Password, cmd.Flag("ssl-password").Usage, 0, 0)
		if flags.SSL.Password == "" {
			return errors.New(L("a password for the new CA key is required"))
		}
	}

	image := podman.GetServiceImage(podman.ServerService)
	tz := adm_podman.GetContainerTimezone()

	if err := adm_podman.AddCA(image, &flags.SSL, tz, fqdn); err != nil {
		return err
	}

	if err := adm_podman.ApplyNewCertificates(fqdn); err != nil {
		return err
	}

	if fingerprint, err := adm_podman.NewCAFingerprint(); err != nil {
		log.Warn().Err(err).Msg(L("Could not read the new CA fingerprint"))
	} else {
		log.Info().Msgf(L("The new root CA (SHA-256 fingerprint %s) has been added to the trusted bundle."),
			fingerprint)
	}

	return nil
}
