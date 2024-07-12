// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/base64"
	"fmt"
	"os"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// PodmanLogin logs in the registry.suse.com registry if needed and returns an authentication file, a cleanup function and an error.
func PodmanLogin() (string, func(), error) {
	creds, err := inspectHost()
	if err != nil {
		return "", nil, err
	}

	if creds.SccPassword != "" && creds.SccUsername != "" {
		// We have SCC credentials, so we are pretty likely to need registry.suse.com
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", creds.SccUsername, creds.SccPassword)))
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
			return "", nil, utils.Errorf(err, L("failed to close the temporary auth file"))
		}

		return authFilePath, func() {
			os.Remove(authFilePath)
		}, nil
	}

	return "", func() {}, nil
}
