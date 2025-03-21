// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// InspectServer check values on a given image and deploy.
func InspectServer(
	namespace string,
	serverImage string,
	pullPolicy string,
	pullSecret string,
) (*utils.ServerInspectData, error) {
	podName := "uyuni-image-inspector"

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer cleaner()
	inspector := utils.NewServerInspector(tempDir)
	// We need the inspector to write to the pod's logs instead of a file
	inspector.DataPath = "/dev/stdout"
	script, err := inspector.GenerateScriptString()
	if err != nil {
		return nil, err
	}

	out, err := RunPodLogs(
		namespace, podName, serverImage, pullPolicy, pullSecret,
		[]types.VolumeMount{utils.EtcRhnVolumeMount, utils.VarPgsqlDataVolumeMount},
		"sh", "-c", script,
	)
	if err != nil {
		return nil, err
	}

	// Parse the data
	inspectedData, err := utils.ReadInspectDataString[utils.ServerInspectData]([]byte(out))
	if err != nil {
		return nil, utils.Errorf(err, L("failed to parse the inspected data"))
	}
	return inspectedData, nil
}
