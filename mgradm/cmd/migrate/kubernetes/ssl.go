// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"os"
	"path"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installExistingCertificate(namespace string, extractedData *MigrationData) error {
	// Store the certificates and key to file to load them
	tmpDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	caCrtPath := path.Join(tmpDir, "ca.crt")
	if err := os.WriteFile(caCrtPath, []byte(extractedData.CaCert), 0700); err != nil {
		return utils.Errorf(err, L("failed to create temporary ca.crt file"))
	}

	srvCrtPath := path.Join(tmpDir, "srv.crt")
	if err := os.WriteFile(srvCrtPath, []byte(extractedData.ServerCert), 0700); err != nil {
		return utils.Errorf(err, L("failed to create temporary srv.crt file"))
	}

	srvKeyPath := path.Join(tmpDir, "srv.key")
	if err := os.WriteFile(srvKeyPath, []byte(extractedData.ServerKey), 0700); err != nil {
		return utils.Errorf(err, L("failed to create temporary srv.key file"))
	}

	sslFlags := adm_utils.InstallSSLFlags{
		Ca: types.CaChain{Root: caCrtPath},
		Server: types.SSLPair{
			Key:  srvKeyPath,
			Cert: srvCrtPath,
		},
	}
	return kubernetes.DeployExistingCertificate(namespace, &sslFlags)
}
