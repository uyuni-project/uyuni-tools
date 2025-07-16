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

// PodmanLogin logs in the SCC registry registry if needed.
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
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", sccUser, sccPassword)))
		authFileContent := fmt.Sprintf(`{
	"auths": {
		"%s" : {
			"auth": "%s"
		}
	}
}`, scc.Registry, token)
		authFile, err := os.CreateTemp("", "mgradm-")
		if err != nil {
			return "", nil, utils.Errorf(err, L("failed to login to %s"), scc.Registry)
		}
		authFilePath := authFile.Name()

		if _, err := authFile.Write([]byte(authFileContent)); err != nil {
			os.Remove(authFilePath)
			return "", nil, utils.Errorf(err, L("failed to login to %s"), scc.Registry)
		}

		if err := authFile.Close(); err != nil {
			os.Remove(authFilePath)
			return "", nil, utils.Errorf(err,
				L("failed to close the temporary auth file. Failed to login to %s"), scc.Registry)
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
