// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"errors"

	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"gopkg.in/yaml.v2"
)

// MigrationData represents the files and data extracted from the migration sync phase.
type MigrationData struct {
	CaKey      string
	CaCert     string
	Data       *utils.InspectResult
	ServerCert string
	ServerKey  string
}

func extractMigrationData(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	volume types.VolumeMount,
) (*MigrationData, error) {
	// Run a pod reading the extracted data files from the volume.
	// The data are written as a YAML dictionary where the key is the file name and the value its content.
	out, err := kubernetes.RunPodLogs(namespace, "uyuni-data-extractor", image,
		pullPolicy, pullSecret, []types.VolumeMount{volume},
		"sh", "-c",
		"for f in /var/lib/uyuni-tools/*; do echo \"`basename $f`: |2\"; cat $f | sed 's/^/  /'; done",
	)
	if err != nil {
		return nil, err
	}

	// Parse the content
	files := make(map[string]string)
	if err := yaml.Unmarshal(out, &files); err != nil {
		return nil, utils.Errorf(err, L("failed to parse data extractor pod output"))
	}

	var result MigrationData
	for file, content := range files {
		switch file {
		case "RHN-ORG-PRIVATE-SSL-KEY":
			result.CaKey = content
		case "RHN-ORG-TRUSTED-SSL-CERT":
			result.CaCert = content
		case "spacewalk.crt":
			result.ServerCert = content
		case "spacewalk.key":
			result.ServerKey = content
		case "data":
			parsedData, err := utils.ReadInspectData[utils.InspectResult]([]byte(content))
			if err != nil {
				return nil, utils.Errorf(err, L("failed to parse migration data file"))
			}
			result.Data = parsedData
		}
	}

	if result.Data == nil {
		return nil, errors.New(L("found no data file after migration"))
	}

	return &result, nil
}
