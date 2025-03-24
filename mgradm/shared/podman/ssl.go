// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func prepareThirdPartyCertificate(caChain *types.CaChain, pair *types.SSLPair, outDir string) error {
	// OrderCas checks the chain of certificates to report problems early
	// We also sort the certificates of the chain in a single blob for Apache and PostgreSQL
	var orderedCert, rootCA []byte
	var err error
	if orderedCert, rootCA, err = ssl.OrderCas(caChain, pair); err != nil {
		return err
	}

	// Check that the private key is not encrypted
	if err := ssl.CheckKey(pair.Key); err != nil {
		return err
	}

	if err := os.Mkdir(outDir, 0600); err != nil {
		return err
	}

	// Write the ordered cert and Root CA to temp files
	caPath := path.Join(outDir, "ca.crt")
	if err = os.WriteFile(caPath, rootCA, 0600); err != nil {
		return err
	}

	serverCertPath := path.Join(outDir, "server.crt")
	if err = os.WriteFile(serverCertPath, orderedCert, 0600); err != nil {
		return err
	}

	return nil
}

var newRunner = utils.NewRunner

// PrepareSSLCertificates prepares SSL environment for the server and database.
// If 3rd party certificates are provided, it uses them, else new certificates are generated.
// This function is called in both new installation and upgrade scenarios.
func PrepareSSLCertificates(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	// Prepare Server certificates
	if err := prepareServerSSLcertificates(image, sslFlags, tz, fqdn); err != nil {
		return err
	}
	// Prepare database certificates
	if err := prepareDatabaseSSLcertificates(image, sslFlags, tz, fqdn); err != nil {
		return err
	}
	return nil
}

func prepareServerSSLcertificates(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	// Check for provided certificates
	if sslFlags.UseProvided() {
		log.Info().Msg(L("Using provided 3rd party server certificates"))
		ca := &sslFlags.Ca
		pair := &sslFlags.Server

		serverDir := path.Join(tempDir, "server")
		if err := prepareThirdPartyCertificate(ca, pair, serverDir); err != nil {
			return err
		}
		// Create secrets for CA
		return shared_podman.CreateTLSSecrets(
			shared_podman.CASecret, path.Join(serverDir, "ca.crt"),
			shared_podman.SSLCertSecret, path.Join(serverDir, "server.crt"),
			shared_podman.SSLKeySecret, pair.Key,
		)
	}
	// Not provided but check for upgrade scenario
	// Check if this is not an upgrade scenario and there is existing CA
	rootCA, err := shared_podman.ReadFromContainer("uyuni-read-ca", image, utils.ServerVolumeMounts, nil,
		ssl.CAContainerPath)
	if err == nil {
		log.Info().Msg(L("Reusing existing server CA certificate"))

		caPath := path.Join(tempDir, "existing-ca.crt")
		if err = os.WriteFile(caPath, rootCA, 0444); err != nil {
			return utils.Error(err, L("cannot write existing CA certificate"))
		}

		return shared_podman.CreateCASecrets(
			shared_podman.CASecret, caPath,
		)
	}

	// Not provided and not an upgrade, generate new
	return generateServerCertificate(image, sslFlags, tz, fqdn)
}

func prepareDatabaseSSLcertificates(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	// Write the ordered cert and Root CA to temp files
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	if sslFlags.UseProvidedDB() {
		log.Info().Msg(L("Using provided 3rd party database certificates"))
		dbCa := &sslFlags.DB.CA
		dbPair := &sslFlags.DB.SSLPair

		dbDir := path.Join(tempDir, "db")
		if err := prepareThirdPartyCertificate(dbCa, dbPair, dbDir); err != nil {
			return err
		}
		// Create secrets for the database key and certificate
		return shared_podman.CreateTLSSecrets(
			shared_podman.DBCASecret, path.Join(dbDir, "ca.crt"),
			shared_podman.DBSSLCertSecret, path.Join(dbDir, "server.crt"),
			shared_podman.DBSSLKeySecret, dbPair.Key,
		)
	}

	// Not provided and not an upgrade, generate new
	return generateDatabaseCertificate(image, sslFlags, tz, fqdn)
}

func runSSLContainer(script string, workdir string, image string, tz string, env map[string]string) error {
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
		"-e", "TZ=" + tz,
		"-v", utils.RootVolumeMount.Name + ":" + utils.RootVolumeMount.MountPath,
		"-v", workdir + ":/ssl:z", // Bind mount for the generated certificates
	}
	command = append(command, envNames...)
	command = append(command, image)

	// Fail fast with `-e`.
	command = append(command, "/usr/bin/sh", "-e", "-c", script)

	_, err := newRunner("podman", command...).Env(envValues).StdMapping().Exec()
	return err
}

func generateServerCertificate(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	env := map[string]string{
		"CERT_O":       sslFlags.Org,
		"CERT_OU":      sslFlags.OU,
		"CERT_CITY":    sslFlags.City,
		"CERT_STATE":   sslFlags.State,
		"CERT_COUNTRY": sslFlags.Country,
		"CERT_EMAIL":   sslFlags.Email,
		"CERT_CNAMES":  strings.Join(append([]string{fqdn}, sslFlags.Cnames...), " "),
		"CERT_PASS":    sslFlags.Password,
		"HOSTNAME":     fqdn,
	}
	if err := runSSLContainer(sslSetupServerScript, tempDir, image, tz, env); err != nil {
		return utils.Error(err, L("Server SSL certificates generation failed"))
	}

	log.Info().Msg(L("Server SSL certificates generated"))

	// Create secret for the database key and certificate
	return shared_podman.CreateTLSSecrets(
		shared_podman.CASecret, path.Join(tempDir, "ca.crt"),
		shared_podman.SSLCertSecret, path.Join(tempDir, "server.crt"),
		shared_podman.SSLKeySecret, path.Join(tempDir, "server.key"),
	)
}

func generateDatabaseCertificate(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	// Write the ordered cert and Root CA to temp files
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	env := map[string]string{
		"CERT_O":       sslFlags.Org,
		"CERT_OU":      sslFlags.OU,
		"CERT_CITY":    sslFlags.City,
		"CERT_STATE":   sslFlags.State,
		"CERT_COUNTRY": sslFlags.Country,
		"CERT_EMAIL":   sslFlags.Email,
		"CERT_CNAMES":  strings.Join(append([]string{fqdn}, sslFlags.Cnames...), " "),
		"CERT_PASS":    sslFlags.Password,
		"HOSTNAME":     fqdn,
	}
	if err := runSSLContainer(sslSetupDatabaseScript, tempDir, image, tz, env); err != nil {
		return utils.Error(err, L("Database SSL certificates generation failed"))
	}

	log.Info().Msg(L("Database SSL certificates generated"))

	// Create secret for the database key and certificate
	if err := shared_podman.CreateTLSSecrets(
		shared_podman.DBCASecret, path.Join(tempDir, "ca.crt"),
		shared_podman.DBSSLCertSecret, path.Join(tempDir, "reportdb.crt"),
		shared_podman.DBSSLKeySecret, path.Join(tempDir, "reportdb.key"),
	); err != nil {
		return err
	}

	return nil
}

const sslSetupServerScript = `
	getMachineName() {
	  hostname="$1"

	  hostname=$(echo "$hostname" | sed 's/\*/_star_/g')

	  field_count=$(echo "$hostname" | awk -F. '{print NF}')

	  if [ "$field_count" -lt 3 ]; then
		echo "$hostname"
		return 0
	  fi

	  end_field=$(expr "$field_count" - 2)

	  result=$(echo "$hostname" | cut -d. -f1-"$end_field")

	  echo "$result"
	}

	echo "Generating the self-signed SSL CA..."
	mkdir -p /root/ssl-build
	rhn-ssl-tool --gen-ca --no-rpm --force --dir /root/ssl-build \
		--password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --set-common-name "$HOSTNAME" --cert-expiration 3650
	cp /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT /ssl/ca.crt

	echo "Generate apache certificate..."
	cert_args=""
	for CERT_CNAME in $CERT_CNAMES; do
		cert_args="$cert_args --set-cname $CERT_CNAME"
	done

	rhn-ssl-tool --gen-server --no-rpm --cert-expiration 3650 \
		--dir /root/ssl-build --password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --set-hostname "$HOSTNAME" --cert-expiration 3650 --set-email "$CERT_EMAIL" \
		$cert_args

	MACHINE_NAME=$(getMachineName "$HOSTNAME")
	cp "/root/ssl-build/$MACHINE_NAME/server.crt" /ssl/server.crt
	cp "/root/ssl-build/$MACHINE_NAME/server.key" /ssl/server.key
`

// This is assuming CA cert is generated by server script.
// If we in any point in the future allow mix of 3rd party server and self signed ca for database
// this will need to be updated to include check for ca cert and build if needed.
const sslSetupDatabaseScript = `
	echo "Generating DB certificate..."
	rhn-ssl-tool --gen-server --no-rpm --cert-expiration 3650 \
		--dir /root/ssl-build --password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --set-hostname reportdb.mgr.internal --cert-expiration 3650 --set-email "$CERT_EMAIL" \
		--set-cname reportdb --set-cname db $cert_args

	cp /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT /ssl/ca.crt
	cp /root/ssl-build/reportdb/server.crt /ssl/reportdb.crt
	cp /root/ssl-build/reportdb/server.key /ssl/reportdb.key
`
