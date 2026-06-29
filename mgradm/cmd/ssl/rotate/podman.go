// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package rotate

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_podman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func rotateForPodman(_ *types.GlobalFlags, flags *rotateFlags, cmd *cobra.Command, args []string) error {
	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}

	if flags.Emergency {
		return emergencyRotate(flags, cmd, fqdn)
	}
	return plannedRotate(flags, cmd, fqdn)
}

func plannedRotate(flags *rotateFlags, cmd *cobra.Command, fqdn string) error {
	checkOnly, _ := cmd.Flags().GetBool("check-only")

	fingerprint, err := adm_podman.NewCAFingerprint()
	if err != nil {
		return err
	}
	if err := checkClientReadiness(fingerprint, flags.Force); err != nil {
		return err
	}
	if checkOnly {
		return nil
	}

	if err := switchServerCertificate(flags, cmd, fqdn, false); err != nil {
		return err
	}
	return adm_podman.ApplyNewCertificates(fqdn)
}

func emergencyRotate(flags *rotateFlags, cmd *cobra.Command, fqdn string) error {
	if checkOnly, _ := cmd.Flags().GetBool("check-only"); checkOnly {
		log.Warn().Msg(L("--check-only is ignored with --emergency! Proceeding with the emergency rotation"))
	}

	if err := switchServerCertificate(flags, cmd, fqdn, true); err != nil {
		return err
	}
	return adm_podman.ApplyNewCertificates(fqdn)
}

func switchServerCertificate(flags *rotateFlags, cmd *cobra.Command, fqdn string, regenerateCA bool) error {
	if flags.SSL.Ca.IsThirdParty() {
		if !flags.SSL.UseProvided() {
			return errors.New(L("the server certificate, key and root CA need to be all provided"))
		}
		log.Info().Msg(L("Installing the provided 3rd party server certificate"))
		return adm_podman.SetThirdPartyCertificates(&flags.SSL, fqdn)
	}

	utils.AskPasswordIfMissing(&flags.SSL.Password, cmd.Flag("ssl-password").Usage, 0, 0)
	if flags.SSL.Password == "" {
		return errors.New(L("the CA key password is required"))
	}

	image := podman.GetServiceImage(podman.ServerService)
	tz := adm_podman.GetContainerTimezone()
	if regenerateCA {
		log.Info().Msg(L("Generating a new CA and server certificate"))
		return adm_podman.RegenerateCAAndCertificate(image, &flags.SSL, tz, fqdn)
	}
	log.Info().Msg(L("Generating a new server certificate signed by the new CA"))
	return adm_podman.RotateServerCertificate(image, &flags.SSL, tz, fqdn)
}

func checkClientReadiness(fingerprint string, force bool) error {
	if force {
		log.Warn().Msg(L("Skipping the client CA trust check as --force was set"))
		return nil
	}

	result, err := adm_podman.CheckClientsCATrust(fingerprint)
	if err != nil {
		return utils.Error(err, L("failed to verify client readiness, rerun with --force to rotate without checking"))
	}

	reportClientReadiness(result)
	if !result.AllMigrated() {
		return errors.New(
			L("some clients do not trust the new CA yet! Distribute it and retry, or use --force to skip this check"),
		)
	}
	return nil
}

func reportClientReadiness(result adm_podman.ClientCheckResult) {
	log.Info().Msgf(L("Clients trusting the new CA: %d"), len(result.Migrated))
	if len(result.NotMigrated) > 0 {
		log.Warn().Msgf(L("Clients not yet trusting the new CA (%[1]d): %[2]s"),
			len(result.NotMigrated), strings.Join(result.NotMigrated, ", "))
	}
	if len(result.Unreachable) > 0 {
		log.Warn().Msgf(L("Unreachable minions (%[1]d): %[2]s"),
			len(result.Unreachable), strings.Join(result.Unreachable, ", "))
	}
	log.Warn().Msg(L("Only Salt-managed clients are checked! Non-managed systems and proxy " +
		"configurations are not verified and must be handled manually."))
}
