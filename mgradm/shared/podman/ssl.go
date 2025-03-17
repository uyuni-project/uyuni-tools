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
	"github.com/uyuni-project/uyuni-tools/shared"
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

// generateSSLCertificates creates the self-signed certificates if needed.
func GenerateSSLCertificates(image string, ssl *adm_utils.InstallSSLFlags, tz string, fqdn string) error {
	// Write the ordered cert and Root CA to temp files
	tempDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return err
	}

	var ca *types.CaChain
	var pair *types.SSLPair

	var dbCa *types.CaChain
	var dbPair *types.SSLPair

	if ssl.UseProvided() {
		ca = &ssl.Ca
		pair = &ssl.Server
	}

	if ca != nil && pair != nil {
		serverDir := path.Join(tempDir, "server")
		if err := prepareThirdPartyCertificate(
			ca, pair, serverDir,
		); err != nil {
			return err
		}

		// Create secrets for the server key and certificate
		if err := shared_podman.CreateTLSSecrets(
			shared_podman.CASecret, path.Join(serverDir, "ca.crt"),
			shared_podman.SSLCertSecret, path.Join(serverDir, "server.crt"),
			shared_podman.SSLKeySecret, pair.Key,
		); err != nil {
			return err
		}
	}

	if ssl.UseProvidedDB() {
		dbCa = &ssl.DB.CA
		dbPair = &ssl.DB.SSLPair
	} else if ssl.UseProvided() {
		//if DB SSL info are not provided, use the server one.
		dbCa = &ssl.Ca
		dbPair = &ssl.Server
	} else if ssl.DB.IsDefined() && !ssl.DB.CA.IsThirdParty() {
		//let's use the CA already present in the server
		dbPair = &ssl.DB.SSLPair
	}

	if dbPair != nil {
		dbDir := path.Join(tempDir, "db")
		if dbCa == nil {
			log.Debug().Msg("A 3rd party SSL was set, but not a CA. Use the CA from the server.")
			cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
			const containerCertPath = "server:/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT"
			if err := cnx.Copy(containerCertPath, dbDir, "root", "root"); err != nil {
				return err
			}
			dbCa = &types.CaChain{Root: path.Join(dbDir)}
		}

		if err := prepareThirdPartyCertificate(
			dbCa,
			dbPair,
			dbDir,
		); err != nil {
			return err
		}

		// Create secrets for the database key and certificate
		if err := shared_podman.CreateTLSSecrets(
			shared_podman.DBCASecret, path.Join(dbDir, "ca.crt"),
			shared_podman.DBSSLCertSecret, path.Join(dbDir, "server.crt"),
			shared_podman.DBSSLKeySecret, dbPair.Key,
		); err != nil {
			return err
		}
		return nil
	}

	env := map[string]string{
		"CERT_O":       ssl.Org,
		"CERT_OU":      ssl.OU,
		"CERT_CITY":    ssl.City,
		"CERT_STATE":   ssl.State,
		"CERT_COUNTRY": ssl.Country,
		"CERT_EMAIL":   ssl.Email,
		"CERT_CNAMES":  strings.Join(append([]string{fqdn}, ssl.Cnames...), " "),
		"CERT_PASS":    ssl.Password,
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
		"-e", "TZ=" + tz,
		"-v", utils.RootVolumeMount.Name + ":" + utils.RootVolumeMount.MountPath,
		"-v", tempDir + ":/ssl:z", // Bind mount for the generated certificates
	}
	command = append(command, envNames...)
	command = append(command, image)

	// Fail fast with `-e`.
	command = append(command, "/usr/bin/sh", "-e", "-c", sslSetupScript)

	if _, err := newRunner("podman", command...).Env(envValues).StdMapping().Exec(); err != nil {
		return utils.Error(err, L("SSL certificates generation failed"))
	}

	log.Info().Msg(L("SSL certificates generated"))

	// Create secret for the database key and certificate
	if err := shared_podman.CreateTLSSecrets(
		shared_podman.CASecret, path.Join(tempDir, "ca.crt"),
		shared_podman.SSLCertSecret, path.Join(tempDir, "server.crt"),
		shared_podman.SSLKeySecret, path.Join(tempDir, "server.key"),
	); err != nil {
		return err
	}

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

const sslSetupScript = `
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

	echo "Generating DB certificate..."
	rhn-ssl-tool --gen-server --no-rpm --cert-expiration 3650 \
		--dir /root/ssl-build --password "$CERT_PASS" \
		--set-country "$CERT_COUNTRY" --set-state "$CERT_STATE" --set-city "$CERT_CITY" \
	    --set-org "$CERT_O" --set-org-unit "$CERT_OU" \
	    --set-hostname reportdb.mgr.internal --cert-expiration 3650 --set-email "$CERT_EMAIL" \
		--set-cname reportdb --set-cname db $cert_args

	cp /root/ssl-build/reportdb/server.crt /ssl/reportdb.crt
	cp /root/ssl-build/reportdb/server.key /ssl/reportdb.key
`
