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

func validateCA(image string, sslFlags *adm_utils.InstallSSLFlags, tz string) error {
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}
	env := map[string]string{
		"CERT_PASS": sslFlags.Password,
	}

	if err := runSSLContainer(sslValidateCA, tempDir, image, tz, env); err != nil {
		return utils.Error(err, L("CA validation failed!"))
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

	// Check if this is an upgrade scenario and there is existing CA
	if reused, err := reuseExistingCertificates(image, tempDir, false); reused && err == nil {
		// We succeffuly loaded existing certificates
		return nil
	} else if reused && err != nil {
		// We found certificates, but there was trouble loading it
		return err
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

	// Check if this is an upgrade scenario and there is existing CA
	if reused, err := reuseExistingCertificates(image, tempDir, true); reused && err == nil {
		// We succeffuly loaded existing certificates
		return nil
	} else if reused && err != nil {
		// We found certificates, but there was trouble loading it
		return err
	}

	// Not provided and not an upgrade, generate new
	return generateDatabaseCertificate(image, sslFlags, tz, fqdn)
}

func reuseExistingCertificates(image string, tempDir string, isDatabaseCheck bool) (bool, error) {
	// Upgrading from 5.1+ with all cerst in secrets
	if reuseExistingCertificatesFromSecrets(isDatabaseCheck) {
		return true, nil
	}

	// Upgrading from 5.0- with all certs in files
	return reuseExistingCertificatesFromMounts(image, tempDir, isDatabaseCheck)
}

func reuseExistingCertificatesFromSecrets(isDatabaseCheck bool) bool {
	if isDatabaseCheck {
		return shared_podman.HasSecret(shared_podman.DBCASecret) &&
			shared_podman.HasSecret(shared_podman.DBSSLCertSecret) &&
			shared_podman.HasSecret(shared_podman.DBSSLKeySecret)
	}
	return shared_podman.HasSecret(shared_podman.CASecret) &&
		shared_podman.HasSecret(shared_podman.SSLCertSecret) &&
		shared_podman.HasSecret(shared_podman.SSLKeySecret)
}

func reuseExistingCertificatesFromMounts(image string, tempDir string, isDatabaseCheck bool) (bool, error) {
	// Basic init
	caPath := path.Join(tempDir, "existing-ca.crt")
	serverCert := path.Join(tempDir, "existing-server.crt")
	serverKey := path.Join(tempDir, "existing-key.crt")

	// No longer used by 5.1+, but contain existing certs in migration scenarios
	etcTLSVolume := types.VolumeMount{Name: "etc-tls", MountPath: "/etc/pki/tls"}

	// Paths for server side checking
	volumes := append(utils.ServerVolumeMounts, etcTLSVolume)
	caCheckPath := ssl.CAContainerPath
	crtCheckPath := ssl.ServerCertPath
	keyCheckPath := ssl.ServerCertKeyPath

	if isDatabaseCheck {
		// Path for database side checking.
		// It is necessary to include etc-tls and ca-cert volume mounts to simulate non-split installation
		volumes = append(utils.PgsqlRequiredVolumeMounts, etcTLSVolume, utils.CaCertVolumeMount)
		caCheckPath = ssl.DBCAContainerPath
		crtCheckPath = ssl.DBCertPath
		keyCheckPath = ssl.DBCertKeyPath
	}

	const containerName = "uyuni-read-certs"

	// Check if we have existing CA
	rootCA, err := shared_podman.ReadFromContainer(containerName, image, volumes, nil,
		caCheckPath)
	if err != nil {
		log.Info().Msgf(L("CA file %s not found. New CA and certificates will be created."), caCheckPath)
		return false, nil
	}

	if err = os.WriteFile(caPath, rootCA, 0444); err != nil {
		return true, utils.Error(err, L("cannot write existing CA certificate"))
	}

	// Check for server certificate
	cert, err := shared_podman.ReadFromContainer(containerName, image, volumes, nil,
		crtCheckPath)
	if err != nil {
		log.Info().Msgf(L("Cert file %s not found. A new certificate will be created."), crtCheckPath)
		return false, nil
	}
	if err = os.WriteFile(serverCert, cert, 0444); err != nil {
		return true, utils.Error(err, L("cannot write existing server certificate"))
	}

	// Check for server certificate key
	keyData, err := shared_podman.ReadFromContainer(containerName, image, volumes, nil,
		keyCheckPath)
	if err != nil {
		log.Info().Msgf(L("Cert key file %s not found. A new certificate will be created."), keyCheckPath)
		return false, nil
	}
	if err = os.WriteFile(serverKey, keyData, 0400); err != nil {
		return true, utils.Error(err, L("cannot write existing server key"))
	}

	log.Info().Msg(L("Reusing existing certificates"))
	if isDatabaseCheck {
		return true, shared_podman.CreateTLSSecrets(
			shared_podman.DBCASecret, caPath,
			shared_podman.DBSSLCertSecret, serverCert,
			shared_podman.DBSSLKeySecret, serverKey,
		)
	}
	return true, shared_podman.CreateTLSSecrets(
		shared_podman.CASecret, caPath,
		shared_podman.SSLCertSecret, serverCert,
		shared_podman.SSLKeySecret, serverKey,
	)
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
		return utils.Error(err, L("Failed to generate server SSL certificates. Please check the input parameters."))
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

	if err := validateCA(image, sslFlags, tz); err != nil {
		return utils.Error(err, L("Cannot generate database certificate"))
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
		return utils.Error(err, L("Failed to generate server Database SSL certificates. Please check the input parameters."))
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
	rhn-ssl-tool --gen-ca --force --dir /root/ssl-build \
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

	rhn-ssl-tool --gen-server --cert-expiration 3650 \
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
	rhn-ssl-tool --gen-server --cert-expiration 3650 \
		--dir /root/ssl-build --password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --set-hostname reportdb.mgr.internal --cert-expiration 3650 --set-email "$CERT_EMAIL" \
		--set-cname reportdb --set-cname db $cert_args

	cp /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT /ssl/ca.crt
	cp /root/ssl-build/reportdb/server.crt /ssl/reportdb.crt
	cp /root/ssl-build/reportdb/server.key /ssl/reportdb.key
`
const sslValidateCA = `
	CA_KEY=/root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY
	CA_PASS_FILE=/ssl/ca_pass
	trap "test -f \"$CA_PASS_FILE\" && /bin/rm -f -- \"$CA_PASS_FILE\" " 0 1 2 3 13 15

	echo "Validating CA..."
	echo "$CERT_PASS" > "$CA_PASS_FILE"

	test -f $CA_KEY || (echo "CA key is not available" && exit 1)
	test -r "$CA_KEY" || (echo "CA key is not readable" && exit 2)

	openssl rsa -noout -in "/root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY" -passin "file:$CA_PASS_FILE" || \
	    (echo "Wrong CA key password" && exit 3)
`
