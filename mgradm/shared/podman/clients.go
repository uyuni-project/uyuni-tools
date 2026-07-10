// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCAFingerprint returns the SHA-256 fingerprint of the CA added by 'mgradm ssl addca'.
func NewCAFingerprint() (string, error) {
	bundle, err := shared_podman.GetSecret(shared_podman.CASecret)
	if err != nil {
		return "", utils.Error(err, L("failed to read the CA bundle secret"))
	}

	fingerprints, err := ssl.SHA256Fingerprints([]byte(bundle))
	if err != nil {
		return "", utils.Error(err, L("failed to compute the CA fingerprints"))
	}

	// A single certificate means no new CA has been added on top of the install one yet.
	if len(fingerprints) < 2 {
		return "", errors.New(L("no newly added CA found in the CA bundle; run 'mgradm ssl addca' first"))
	}
	return fingerprints[len(fingerprints)-1], nil
}

type ClientCheckResult struct {
	Migrated    []string
	NotMigrated []string
	Unreachable []string
}

// AllMigrated returns true when every responsive minion trusts the new CA and none is unreachable.
func (r ClientCheckResult) AllMigrated() bool {
	return len(r.NotMigrated) == 0 && len(r.Unreachable) == 0
}

const (
	// saltPingTimeout bounds, in seconds, how long Salt waits for the minions to acknowledge the
	// presence ping when listing their status.
	saltPingTimeout = 3
	// saltGatherJobTimeout bounds, in seconds, how long Salt waits to gather the minions' answers.
	saltGatherJobTimeout = 15
)

// CheckClientsCATrust checks over Salt if the registered minions already trust the CA with
// the given SHA-256 fingerprint.
func CheckClientsCATrust(fingerprint string) (ClientCheckResult, error) {
	var result ClientCheckResult
	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")

	log.Info().Msg(L("Checking whether the clients trust the new CA over Salt; this can take a moment…"))

	// Get both the responsive (up) and unresponsive (down) minions in a single runner call.
	// Bound the Salt timeouts so the check cannot hang on slow or unresponsive minions.
	statusOut, err := cnx.Exec("salt-run", "--out=json", "manage.status",
		fmt.Sprintf("timeout=%d", saltPingTimeout),
		fmt.Sprintf("gather_job_timeout=%d", saltGatherJobTimeout),
	)
	if err != nil {
		return result, utils.Error(err, L("failed to list the minions status; is Salt available on the server?"))
	}
	var status struct {
		Up   []string `json:"up"`
		Down []string `json:"down"`
	}
	if err := json.Unmarshal(statusOut, &status); err != nil {
		return result, utils.Errorf(err, L("failed to parse the minions status: %s"), string(statusOut))
	}
	result.Unreachable = status.Down

	// Nothing to check if no minion is responding.
	if len(status.Up) == 0 {
		return result, nil
	}

	// Target the responsive minions explicitly so Salt returns as soon as they answer, instead of
	// waiting for stragglers like it does with a glob target.
	out, err := cnx.Exec(
		"salt", "--out=json", "--static", fmt.Sprintf("--timeout=%d", saltGatherJobTimeout),
		"-L", strings.Join(status.Up, ","),
		"cmd.run", clientTrustCheckCommand(fingerprint),
	)
	if err != nil {
		return result, utils.Error(err, L("failed to query the minions trust stores over Salt"))
	}

	result.Migrated, result.NotMigrated, err = parseClientTrustResults(out, nil)
	return result, err
}

// parseClientTrustResults parses the JSON output of the Salt `cmd.run` trust check into the lists of
// migrated and not-yet-migrated minions, skipping the ones already known to be unreachable.
func parseClientTrustResults(saltOutput []byte, unreachable map[string]bool) (
	migrated []string, notMigrated []string, err error,
) {
	results := map[string]string{}
	if len(bytes.TrimSpace(saltOutput)) > 0 {
		if err := json.Unmarshal(saltOutput, &results); err != nil {
			return nil, nil, utils.Errorf(err, L("failed to parse the minions response: %s"), string(saltOutput))
		}
	}

	for minion, output := range results {
		if unreachable[minion] {
			continue
		}
		if strings.TrimSpace(output) == "OK" {
			migrated = append(migrated, minion)
		} else {
			notMigrated = append(notMigrated, minion)
		}
	}
	return migrated, notMigrated, nil
}

// clientTrustCheckScript is the shell script run on each minion to check whether the
// CA file holds a certificate matching a given SHA-256 fingerprint.
// It prints "OK" on a match and "MISSING" otherwise.
const clientTrustCheckScript = `
found=MISSING

for caFile in \
	/etc/pki/trust/anchors/RHN-ORG-TRUSTED-SSL-CERT \
	/etc/pki/ca-trust/source/anchors/RHN-ORG-TRUSTED-SSL-CERT \
	/usr/local/share/ca-certificates/RHN-ORG-TRUSTED-SSL-CERT.crt; do
	# Skip the locations that do not apply to this minion's OS.
	[ -f "$caFile" ] || continue

	# Split the (possibly multi-cert) CA file into one file per certificate.
	d=$(mktemp -d)
	csplit -z -s -f "$d/cert-" "$caFile" '/-----BEGIN CERTIFICATE-----/' '{*}' 2>/dev/null

	# Report OK as soon as one of the certificates matches the expected fingerprint.
	for cert in "$d"/cert-*; do
		fp=$(openssl x509 -in "$cert" -noout -fingerprint -sha256 2>/dev/null | cut -d= -f2)
		[ "$fp" = "%[1]s" ] && found=OK
	done

	rm -rf "$d"
done

echo "$found"
`

// clientTrustCheckCommand returns the trust check shell command to run on each minion for the given
// SHA-256 fingerprint.
func clientTrustCheckCommand(fingerprint string) string {
	return fmt.Sprintf(clientTrustCheckScript, fingerprint)
}
