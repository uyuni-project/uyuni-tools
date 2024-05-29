// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// GetSupportConfigPath returns the support config tarball path.
func GetSupportConfigPath(out string) string {
	re := regexp.MustCompile(`/var/log/scc_(.*?)\.txz`)
	return re.FindString(out)
}

// GetSupportConfigFileSaveName returns the support config file name.
func GetSupportConfigFileSaveName() string {
	hostname_b, err := RunCmdOutput(zerolog.DebugLevel, "hostname")
	hostname := "localhost"
	if err != nil {
		log.Warn().Err(err).Msg(L("Unable to detect hostname, using localhost"))
		hostname = strings.TrimSpace(string(hostname_b))
	}
	now := time.Now()
	return fmt.Sprintf("scc_%s_%s", hostname, now.Format("20060102_1504"))
}
