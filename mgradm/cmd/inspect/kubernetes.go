// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package inspect

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"

	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

var kubernetesBuilt = true

func InspectKubernetes(serverImage string, pullPolicy string) (map[string]string, error) {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return map[string]string{}, fmt.Errorf("install %s before running this command. %s", binary, err)
		}
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to create temporary directory. %s", err)
	}

	inspect_shared.GenerateInspectScript(scriptDir)

	command := inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename

	podName := "inspector"

	nodeName, err := shared_kubernetes.GetNode("uyuni")

	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot find node for app uyuni %s", err)
	}

	overridesArgs := []string{"--override-type=strategic", "--overrides", `{"apiVersion":"v1","spec":{"nodeName":"` + nodeName + `","restartPolicy":"Never","containers":[{"name":` + strconv.Quote(podName) + `,"image":` + strconv.Quote(serverImage) + `,"volumeMounts":[{"mountPath":"/var/lib/pgsql","name":"var-pgsql"},{"mountPath":"` + inspect_shared.InspectOutputFile.Directory + `","name":"var-lib-uyuni-tools"}]}],"volumes":[{"name":"var-pgsql","persistentVolumeClaim":{"claimName":"var-pgsql"}},{"name":"var-lib-uyuni-tools","hostPath":{"path":` + strconv.Quote(scriptDir) + `,"type":"Directory"}}]}}`}

	shared_kubernetes.RunPod(podName, serverImage, pullPolicy, command, overridesArgs...)

	inspectResult, err := inspect_shared.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot inspect data. %s", err)
	}

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot print inspect result. %s", err)
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return inspectResult, err
}
