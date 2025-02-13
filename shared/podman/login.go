// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// PodmanLogin logs in the registry.suse.com registry if needed.
//
// It returns an authentication file, a cleanup function and an error.
func PodmanLogin(hostData *HostInspectData, scc types.SCCCredentials) (string, func(), error) {
	sccUser := hostData.SCCUsername
	sccPassword := hostData.SCCPassword
	if scc.User != "" && scc.Password != "" {
		log.Info().Msg(L("SCC credentials parameters will be used. SCC credentials from host will be ignored."))
		sccUser = scc.User
		sccPassword = scc.Password
	}
	if sccUser != "" && sccPassword != "" {
		// We have SCC credentials, so we are pretty likely to need registry.suse.com
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", sccUser, sccPassword)))
		authFileContent := fmt.Sprintf(`{
	"auths": {
		"registry.suse.com" : {
			"auth": "%s"
		}
	}
}`, token)
		authFile, err := os.CreateTemp("", "mgradm-")
		if err != nil {
			return "", nil, err
		}
		authFilePath := authFile.Name()

		if _, err := authFile.Write([]byte(authFileContent)); err != nil {
			os.Remove(authFilePath)
			return "", nil, err
		}

		if err := authFile.Close(); err != nil {
			os.Remove(authFilePath)
			return "", nil, utils.Error(err, L("failed to close the temporary auth file"))
		}

		return authFilePath, func() {
			os.Remove(authFilePath)
		}, nil
	}

	noopCleaner := func() {
		// Nothing to clean
	}

	return "", noopCleaner, nil
}
