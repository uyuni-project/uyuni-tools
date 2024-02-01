// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package inspect

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"

	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

var kubernetesBuilt = true

func inspectKubernetes(serverImage string, pullPolicy string) (map[string]string, error) {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			log.Fatal().Err(err).Msgf("install %s before running this command", binary)
		}
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create temporary directory")
	}

	inspect_shared.GenerateInspectScript(scriptDir)

	command := inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename

	podName := "inspector"

	nodeName := shared_kubernetes.GetNode("uyuni")

	overridesArgs := []string{"--override-type=strategic", "--overrides", `{"apiVersion":"v1","spec":{"restartPolicy":"Never","containers":[{"name":` + strconv.Quote(podName) + `,"image":` + strconv.Quote(serverImage) + `,"volumeMounts":[{"mountPath":"` + inspect_shared.InspectOutputFile.Directory + `","name":"var-lib-uyuni-tools"}]}],"volumes":[{"name":"var-lib-uyuni-tools","hostPath":{"path":` + strconv.Quote(scriptDir) + `,"type":"Directory"}}]}}`}

	//TODO quick and dirty, but it works, just add nodeName, mountPath and volume var-pgsql as PVC
	if len(nodeName) > 0 {
		overridesArgs = []string{"--override-type=strategic", "--overrides", `{"apiVersion":"v1","spec":{"nodeName":"` + nodeName + `","restartPolicy":"Never","containers":[{"name":` + strconv.Quote(podName) + `,"image":` + strconv.Quote(serverImage) + `,"volumeMounts":[{"mountPath":"/var/lib/pgsql","name":"var-pgsql"},{"mountPath":"` + inspect_shared.InspectOutputFile.Directory + `","name":"var-lib-uyuni-tools"}]}],"volumes":[{"name":"var-pgsql","persistentVolumeClaim":{"claimName":"var-pgsql"}},{"name":"var-lib-uyuni-tools","hostPath":{"path":` + strconv.Quote(scriptDir) + `,"type":"Directory"}}]}}`}
	}

	shared_kubernetes.RunPod(podName, serverImage, pullPolicy, command, overridesArgs...)

	shared_kubernetes.WaitForPod(podName, "Succeeded")

	inspectResult := inspect_shared.ReadInspectData(scriptDir)

	shared_kubernetes.DeletePod(podName)

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot print inspect result")
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return inspectResult, err
}
