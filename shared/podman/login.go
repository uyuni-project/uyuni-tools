// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// GeneratePodmanLoginContent
//
// it return a string with the content of the authfile.
func GeneratePodmanLoginContent(
	hostData *HostInspectData, registry types.Registry, scc types.SCCCredentials) string {
	authFileContent := ""
	User := hostData.SCCUsername
	Password := hostData.SCCPassword
	if strings.Contains(registry.Host, "registry.suse.com") {
		log.Info().Msg(L("Registry is registry.suse.com. Using SCC credentials."))
		if scc.User != "" && scc.Password != "" {
			log.Info().Msg(L("SCC credentials has been provided, SCChost credentials will be ignored."))
			User = scc.User
			Password = scc.Password
		}
	} else if registry.User != "" && registry.Password != "" {
		log.Info().Msg(L("Registry parameters will be used. SCC credentials from host will be ignored."))
		User = registry.User
		Password = registry.Password
	}
	if User != "" && Password != "" {
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", User, Password)))
		authFileContent = fmt.Sprintf(`{
	"auths": {
		"%s" : {
			"auth": "%s"
		}
	}
}`, registry.Host, token)
	}
	return authFileContent
}

// PodmanLogin logs in the registry if needed.
//
// It returns an authentication file, a cleanup function and an error.
func PodmanLogin(hostData *HostInspectData, registry types.Registry, scc types.SCCCredentials) (string, func(), error) {
	authFileContent := GeneratePodmanLoginContent(hostData, registry, scc)
	if authFileContent != "" {
		authFile, err := os.CreateTemp("", "mgradm-")
		if err != nil {
			return "", nil, utils.Errorf(err, L("failed to set authentication for %s"), registry.Host)
		}
		authFilePath := authFile.Name()

		if _, err := authFile.Write([]byte(authFileContent)); err != nil {
			os.Remove(authFilePath)
			return "", nil, utils.Errorf(err, L("failed to set authentication for %s"), registry.Host)
		}

		if err := authFile.Close(); err != nil {
			os.Remove(authFilePath)
			return "", nil, utils.Errorf(err,
				L("failed to close the temporary auth file. Cannot set authentication for %s"), registry.Host)
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
