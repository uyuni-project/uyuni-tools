// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
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

func prepareThirdPartyCertificate(
	caChain *types.CaChain,
	pair *types.SSLPair,
	caSecretName string,
	certSecretName string,
	keySecretName string,
	fqdns ...string,
) error {
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	// OrderCas checks the chain of certificates to report problems early
	// We also sort the certificates of the chain in a single blob for Apache and PostgreSQL
	var orderedCert, rootCA []byte
	if orderedCert, rootCA, err = ssl.OrderCas(caChain, pair); err != nil {
		return err
	}

	// Check that the private key is not encrypted
	if err := ssl.CheckKey(pair.Key); err != nil {
		return err
	}

	// Write the ordered cert and Root CA to temp files
	caPath := path.Join(tempDir, "ca.crt")
	if err = os.WriteFile(caPath, rootCA, 0600); err != nil {
		return err
	}

	serverCertPath := path.Join(tempDir, "server.crt")
	if err = os.WriteFile(serverCertPath, orderedCert, 0600); err != nil {
		return err
	}

	errors := []error{}
	for _, fqdn := range fqdns {
		errors = append(errors, ssl.VerifyHostname(caPath, serverCertPath, fqdn))
	}

	err = utils.JoinErrors(errors...)
	if err != nil {
		return err
	}

	// Create secrets for CA
	return shared_podman.CreateTLSSecrets(
		caSecretName, path.Join(tempDir, "ca.crt"),
		certSecretName, path.Join(tempDir, "server.crt"),
		keySecretName, pair.Key,
	)
}

var newRunner = utils.NewRunner

func prepareThirdPartyCertificates(
	sslFlags *adm_utils.InstallSSLFlags, fqdn string,
) (serverCertReady bool, dbCertReady bool, err error) {
	var errs []error

	// Check if we have been provided certificates as parameters
	if sslFlags.UseProvided() {
		log.Info().Msg(L("Using provided 3rd party server certificates"))
		if err := prepareThirdPartyCertificate(&sslFlags.Ca, &sslFlags.Server,
			shared_podman.CASecret, shared_podman.SSLCertSecret, shared_podman.SSLKeySecret, fqdn,
		); err != nil {
			errs = append(errs, err)
		}
		serverCertReady = true
	}

	if sslFlags.UseProvidedDB() {
		log.Info().Msg(L("Using provided 3rd party database certificates"))
		if err := prepareThirdPartyCertificate(
			&sslFlags.DB.CA, &sslFlags.DB.SSLPair, shared_podman.DBCASecret,
			shared_podman.DBSSLCertSecret, shared_podman.DBSSLKeySecret, fqdn, "db", "reportdb",
		); err != nil {
			errs = append(errs, err)
		}
		dbCertReady = true
	}

	err = utils.JoinErrors(errs...)
	return
}

// PrepareSSLCertificates prepares SSL environment for the server and database.
// If 3rd party certificates are provided, it uses them, else new certificates are generated.
// This function is called in both new installation and upgrade scenarios.
func PrepareSSLCertificates(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	var serverCertReady bool
	var dbCertReady bool

	serverCertReady, dbCertReady, err := prepareThirdPartyCertificates(sslFlags, fqdn)
	if err != nil {
		return utils.Errorf(err, L("Failed to create secrets from the provided SSL arguments"))
	}

	// Do we have secrets or certificates from volumes to reuse?
	if !serverCertReady {
		// Check if this is an upgrade scenario and there is existing CA and cert/key pair
		serverCertReady, err = reuseExistingCertificates(image, fqdn, false)
	}
	if serverCertReady && err != nil {
		// we found certificates, but there was trouble loading it
		return err
	}

	if !dbCertReady {
		// Check if this is an upgrade scenario and there is existing CA and cert/key pair
		dbCertReady, err = reuseExistingCertificates(image, fqdn, true)
	}
	if dbCertReady && err != nil {
		// we found certificates, but there was trouble loading it
		return err
	}

	if dbCertReady && serverCertReady {
		log.Info().Msg(L("Reusing the existing server and database certificates"))
		return nil
	}

	log.Info().Msg(L("Generating both the server and database certificates since one is missing"))
	// No need to check the CA if there is no server cert to reuse: it's likely not there too as in install.
	if serverCertReady {
		// Do we have generated certificates files?
		_, err = shared_podman.ReadFromContainer(
			"ca-key-reader", image, []types.VolumeMount{utils.RootVolumeMount}, []string{},
			"/root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY",
		)
		if err != nil {
			return utils.Error(err,
				L("Cannot generate certificates as the SSL CA key cannot be found. Please set up third-party certificates."))
		}

		if err := validateCA(image, sslFlags, tz); err != nil {
			return err
		}
	}

	// Generate them all in order to have the same expiration date on both.
	return utils.JoinErrors(
		generateServerCertificate(image, sslFlags, tz, fqdn),
		generateDatabaseCertificate(image, sslFlags, tz, fqdn),
	)
}

func validateCA(image string, sslFlags *adm_utils.InstallSSLFlags, tz string) error {
	log.Info().Msg(L("Verifying the CA key password…"))
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}
	env := map[string]string{
		"CERT_PASS": sslFlags.Password,
	}

	if err := runSSLContainer(sslValidateCA, tempDir, image, tz, env); err != nil {
		return errors.New(L("Failed to verify the CA key password!"))
	}
	return nil
}

func reuseExistingCertificates(image string, fqdn string, isDatabaseCheck bool) (reused bool, err error) {
	// Write the ordered cert and Root CA to temp files
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return
	}

	// Upgrading from 5.1+ with all certificates as secrets
	if reuseExistingCertificatesFromSecrets(isDatabaseCheck) {
		secretName := shared_podman.SSLCertSecret
		caSecretName := shared_podman.CASecret
		fqdns := []string{fqdn}
		msg := L("Reusing the existing server certificate secrets")
		if isDatabaseCheck {
			secretName = shared_podman.DBSSLCertSecret
			caSecretName = shared_podman.DBCASecret
			fqdns = append(fqdns, "db", "reportdb")
			msg = L("Reusing the existing database certificate secrets")
		}
		reused = isFQDNMatchingCertificateSecret(secretName, caSecretName, fqdns...)
		if reused {
			log.Info().Msg(msg)
		}
		return
	}

	// Upgrading from 5.0- with all certs in files
	return reuseExistingCertificatesFromMounts(image, tempDir, fqdn, isDatabaseCheck)
}

func isFQDNMatchingCertificateSecret(secretName string, caSecretName string, fqdns ...string) bool {
	cert, err := shared_podman.GetSecret(secretName)
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}

	tmpDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}
	caPath := path.Join(tmpDir, "ca.crt")
	certPath := path.Join(tmpDir, "toverify.crt")

	caCert, err := shared_podman.GetSecret(caSecretName)
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}

	if err := os.WriteFile(caPath, []byte(caCert), 0700); err != nil {
		log.Error().Err(err).Send()
		return false
	}

	if err := os.WriteFile(certPath, []byte(cert), 0700); err != nil {
		log.Error().Err(err).Send()
		return false
	}

	missingFQDNs := []string{}
	for _, fqdn := range fqdns {
		if !isFQDNMatchingCertificate(fqdn, certPath, caPath) {
			missingFQDNs = append(missingFQDNs, fqdn)
		}
	}
	if len(missingFQDNs) > 0 {
		log.Error().Msgf(L("Missing SANs or subjects in certificate of secret %[1]: %[2]s"),
			secretName, strings.Join(fqdns, ", "))
		return false
	}
	return true
}

func isFQDNMatchingCertificate(fqdn string, certPath string, caPath string) bool {
	err := ssl.VerifyHostname(caPath, certPath, fqdn)
	if err != nil {
		log.Debug().Msgf("SSL verification error: %s", err)
	}
	return err == nil
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

func reuseExistingCertificatesFromMounts(
	image string,
	tempDir string,
	fqdn string,
	isDatabaseCheck bool,
) (bool, error) {
	caPath := path.Join(tempDir, "existing-ca.crt")
	serverCert := path.Join(tempDir, "existing-server.crt")
	serverKey := path.Join(tempDir, "existing-key.crt")

	msg := L("Reusing the existing server certificate files")

	// Paths for server side checking
	caCheckPath := ssl.CAContainerPath
	crtCheckPath := ssl.ServerCertPath
	keyCheckPath := ssl.ServerCertKeyPath

	if isDatabaseCheck {
		// Path for database side checking.
		caCheckPath = ssl.DBCAContainerPath
		crtCheckPath = ssl.DBCertPath
		keyCheckPath = ssl.DBCertKeyPath
		msg = L("Reusing the existing database certificate files")
	}

	const containerName = "uyuni-read-certs"

	// Check if we have existing CA
	rootCA, err := shared_podman.ReadFromContainer(containerName, image, utils.SSLMigrationVolumeMounts, nil,
		caCheckPath)
	if err != nil {
		log.Info().Msgf(L("CA file %s not found."), caCheckPath)
		return false, nil
	}

	if err = os.WriteFile(caPath, rootCA, 0444); err != nil {
		return true, utils.Error(err, L("cannot write existing CA certificate"))
	}

	// Check for server certificate
	cert, err := shared_podman.ReadFromContainer(containerName, image, utils.SSLMigrationVolumeMounts, nil,
		crtCheckPath)
	if err != nil {
		log.Info().Msgf(L("Certificate file %s not found."), crtCheckPath)
		return false, nil
	}
	if err = os.WriteFile(serverCert, cert, 0444); err != nil {
		return true, utils.Error(err, L("cannot write existing server certificate"))
	}

	// We cannot reuse certificates not matching the requested FQDN
	if !isFQDNMatchingCertificate(fqdn, serverCert, caPath) {
		log.Warn().Msgf(L("Certificate file %[1]s doesn't match %[2]s FQDN."), crtCheckPath, fqdn)
		return false, nil
	}

	// Check for server certificate key
	keyData, err := shared_podman.ReadFromContainer(containerName, image, utils.SSLMigrationVolumeMounts, nil,
		keyCheckPath)
	if err != nil {
		log.Warn().Msgf(L("Certificate key file %s not found."), keyCheckPath)
		return false, nil
	}
	if err = os.WriteFile(serverKey, keyData, 0400); err != nil {
		return true, utils.Error(err, L("cannot write existing server key"))
	}

	log.Info().Msg(msg)
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
	// This generally should not happen, otherwise we would ask for CA password in parameters check.
	// However there are some paths, e.g. upgrade reusing existing certs and db provided 3rd party certs where we
	// do not check for the password, but on existing certs check we can fail and drop here.
	if sslFlags.Password == "" {
		return errors.New(L("Cannot generate new certificates without a CA password. Please check input options"))
	}
	log.Info().Msg(L("Generating the server certificate…"))

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
		return utils.Error(err, L("Failed to generate server SSL certificate. Please check the input parameters."))
	}

	log.Info().Msg(L("Server SSL certificate generated"))

	// Create secret for the database key and certificate
	return shared_podman.CreateTLSSecrets(
		shared_podman.CASecret, path.Join(tempDir, "ca.crt"),
		shared_podman.SSLCertSecret, path.Join(tempDir, "server.crt"),
		shared_podman.SSLKeySecret, path.Join(tempDir, "server.key"),
	)
}

func generateDatabaseCertificate(image string, sslFlags *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	log.Info().Msg(L("Generating the database certificate…"))
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
	}
	if err := runSSLContainer(sslSetupDatabaseScript, tempDir, image, tz, env); err != nil {
		return utils.Error(err, L("Failed to generate database SSL certificate"))
	}

	log.Info().Msg(L("Database SSL certificate generated"))

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

	# Only generate a CA is we don't have it yet, like for install
	if ! test -e /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT; then
		echo "Generating the self-signed SSL CA..."
		mkdir -p /root/ssl-build
		rhn-ssl-tool --gen-ca --force --dir /root/ssl-build \
			--password "$CERT_PASS" \
			--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
			--set-org "$CERT_O" --set-org-unit "$CERT_OU" \
			--set-common-name "$HOSTNAME" --cert-expiration 3650
	fi

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
	cert_args=""
	for CERT_CNAME in $CERT_CNAMES; do
		cert_args="$cert_args --set-cname $CERT_CNAME"
	done

	rhn-ssl-tool --gen-server --cert-expiration 3650 \
		--dir /root/ssl-build --password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --cert-expiration 3650 --set-email "$CERT_EMAIL" \
		--set-hostname reportdb --set-cname db $cert_args

	cp /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT /ssl/ca.crt
	cp /root/ssl-build/reportdb/server.crt /ssl/reportdb.crt
	cp /root/ssl-build/reportdb/server.key /ssl/reportdb.key
`
const sslValidateCA = `
	CA_KEY=/root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY
	CA_PASS_FILE=/ssl/ca_pass
	trap "test -f \"$CA_PASS_FILE\" && /bin/rm -f -- \"$CA_PASS_FILE\" " 0 1 2 3 13 15

	echo "$CERT_PASS" > "$CA_PASS_FILE"

	test -f $CA_KEY || (echo "CA key is not available" && exit 1)
	test -r "$CA_KEY" || (echo "CA key is not readable" && exit 2)

	openssl rsa -noout -in "/root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY" -passin "file:$CA_PASS_FILE" || \
	    (echo "Wrong CA key password" 1>&2 && exit 3)
`
