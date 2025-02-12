// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var noopCleaner = func() {
	// Nothing to clean
}

// generateSSLCertificates creates the self-signed certificates if needed.
// It returns the podman arguments to mount them in the setup container, a cleaner function and possibly an error.
func generateSSLCertificates(image string, flags *adm_utils.ServerFlags, fqdn string) ([]string, func(), error) {
	if flags.Installation.SSL.UseExisting() {
		// OrderCas checks the chain of certificates to report problems early
		if _, _, err := ssl.OrderCas(&flags.Installation.SSL.Ca, &flags.Installation.SSL.Server); err != nil {
			return []string{}, noopCleaner, err
		}

		// Check that the private key is not encrypted
		if err := ssl.CheckKey(flags.Installation.SSL.Server.Key); err != nil {
			return []string{}, noopCleaner, err
		}

		// Add mount options for the existing files into the ssl directory
		// The mgr-setup script expects the certificates in the /ssl folder with fixed names.
		opts := []string{
			"-v", flags.Installation.SSL.Ca.Root + ":/ssl/ca.crt",
			"-v", flags.Installation.SSL.Server.Cert + ":/ssl/server.crt",
			"-v", flags.Installation.SSL.Server.Key + ":/ssl/server.key",
		}
		for i, intermediate := range flags.Installation.SSL.Ca.Intermediate {
			opts = append(opts, "-v", fmt.Sprintf("%s:/ssl/intermediate-%d.crt", intermediate, i))
		}

		return opts, noopCleaner, nil
	}

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return []string{}, cleaner, err
	}

	env := map[string]string{
		"CERT_O":       flags.Installation.SSL.Org,
		"CERT_OU":      flags.Installation.SSL.OU,
		"CERT_CITY":    flags.Installation.SSL.City,
		"CERT_STATE":   flags.Installation.SSL.State,
		"CERT_COUNTRY": flags.Installation.SSL.Country,
		"CERT_EMAIL":   flags.Installation.SSL.Email,
		"CERT_CNAMES":  strings.Join(append([]string{fqdn}, flags.Installation.SSL.Cnames...), " "),
		"CERT_PASS":    flags.Installation.SSL.Password,
		"HOSTNAME":     fqdn,
	}
	envNames := []string{}
	envValues := []string{}
	for key, value := range env {
		envNames = append(envNames, "-e", key)
		envValues = append(envValues, fmt.Sprintf("%s=%s", key, value))
	}

	command := []string{
		"run",
		"--rm",
		"--name", "uyuni-ssl-generator",
		"--network", shared_podman.UyuniNetwork,
		"-e", "TZ=" + flags.Installation.TZ,
		"-v", utils.RootVolumeMount.Name + ":" + utils.RootVolumeMount.MountPath,
		"-v", tempDir + ":/ssl:z", // Bind mount for the generated certificates
	}
	command = append(command, envNames...)
	command = append(command, image)

	command = append(command, "/usr/bin/sh", "-x", "-c", sslSetupScript)

	if _, err := newRunner("podman", command...).Env(envValues).StdMapping().Exec(); err != nil {
		return []string{}, cleaner, utils.Error(err, L("SSL certificates generation failed"))
	}

	log.Info().Msg(L("SSL certificates generated"))

	return []string{"-v", tempDir + ":/ssl"}, cleaner, nil
}

const sslSetupScript = `
	echo "Generating the self-signed SSL CA..."
	mkdir -p /root/ssl-build
	rhn-ssl-tool --gen-ca --no-rpm --force --dir /root/ssl-build \
		--password $CERT_PASS \
		--set-country $CERT_COUNTRY --set-state $CERT_STATE --set-city $CERT_CITY \
	    --set-org $CERT_O --set-org-unit $CERT_OU \
	    --set-common-name $HOSTNAME --cert-expiration 3650
	cp /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT /ssl/ca.crt

	echo "Generate apache certificate..."
	cert_args=""
	for CERT_CNAME in $CERT_CNAMES; do
		cert_args=$cert_args --set-cname $CERT_CNAME
	done

	rhn-ssl-tool --gen-server --no-rpm --cert-expiration 3650 \
		--dir /root/ssl-build --password $CERT_PASS \
		--set-country $CERT_COUNTRY --set-state $CERT_STATE --set-city $CERT_CITY \
	    --set-org $CERT_O --set-org-unit $CERT_OU \
	    --set-hostname $HOSTNAME --cert-expiration 3650 --set-email $CERT_EMAIL \
		$cert_args

	NAME=${HOSTNAME%%.*}
	cp /root/ssl-build/${NAME}/server.crt /ssl/server.crt
	cp /root/ssl-build/${NAME}/server.key /ssl/server.key

	echo "Generating DB certificate..."
	rhn-ssl-tool --gen-server --no-rpm --cert-expiration 3650 \
		--dir /root/ssl-build --password $CERT_PASS \
		--set-country $CERT_COUNTRY --set-state $CERT_STATE --set-city $CERT_CITY \
	    --set-org $CERT_O --set-org-unit $CERT_OU \
	    --set-hostname reportdb.mgr.internal --cert-expiration 3650 --set-email $CERT_EMAIL \
		$cert_args

	cp /root/ssl-build/reportdb/server.crt /ssl/reportdb.crt
	cp /root/ssl-build/reportdb/server.key /ssl/reportdb.key
`
