// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

// Package backend handles detection of the container runtime used to run
// a given Uyuni workload (podman, podman-remote, kubectl, or host).
//
// Detection follows a fixed precedence order:
//
//  1. kubectl with a live deployment matching kubernetesFilter
//  2. podman / podman-remote with the target container already running
//  3. podman / podman-remote with uyuni-server or uyuni-proxy-pod service installed
//  4. kubectl with a Helm release named "uyuni" or "uyuni-proxy"
//
// Each I/O concern (PATH lookup, container probing, systemd, kubernetes) is
// behind its own interface so unit tests can swap in fakes without spawning
// real processes.
package backend

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// PathLookup abstracts exec.LookPath so tests can control which binaries
// appear to be installed without modifying PATH.
type PathLookup interface {
	LookPath(file string) bool
}

// ContainerProber checks whether a named container is currently running.
type ContainerProber interface {
	// HasRunningContainer returns true when the named container is present and
	// running under the given binary (podman or podman-remote).
	HasRunningContainer(bin, container string) bool
}

// SystemdInspector reports whether a systemd service unit file is installed.
type SystemdInspector interface {
	HasService(name string) bool
}

// KubernetesProber detects live Kubernetes workloads and Helm releases.
type KubernetesProber interface {
	// HasDeployment returns (true, nil) when at least one deployment matching
	// filter is found in any namespace. It returns (false, err) when the
	// cluster is unreachable so the caller can log a meaningful message.
	HasDeployment(filter string) (bool, error)

	// HasHelmRelease returns true when a Helm release with the given name
	// exists in the cluster.
	HasHelmRelease(release string) bool
}

// Runner is the process-runner factory used throughout the project.
type Runner func(command string, args ...string) types.Runner

// BackendDetector resolves the backend command for a given Uyuni deployment.
// An empty explicit value triggers auto-detection.
type BackendDetector interface {
	// Detect returns the resolved command name ("podman", "podman-remote",
	// "kubectl", or "host"), or an error when no suitable backend is found.
	//
	// explicit is the operator-supplied --backend value, or "" for auto-detect.
	// container is the podman container name (e.g. "uyuni-server").
	// kubernetesFilter is the -l selector used to find the k8s deployment.
	Detect(explicit, container, kubernetesFilter string) (string, error)
}

// SystemDetector is the production BackendDetector. Its fields are exported
// so NewConnection can wire Systemd after construction without an extra
// constructor argument.
type SystemDetector struct {
	Path       PathLookup
	Containers ContainerProber
	Systemd    SystemdInspector
	Kubernetes KubernetesProber
}

// NewSystemDetector returns a SystemDetector wired to real system calls.
// The Systemd field is left nil; callers set it from the Connection's systemd instance.
func NewSystemDetector(r Runner) *SystemDetector {
	return &SystemDetector{
		Path:       &execPathLookup{},
		Containers: &podmanProber{runner: r},
		Systemd:    nil,
		Kubernetes: &kubectlProber{runner: r},
	}
}

// Detect implements BackendDetector.
func (d *SystemDetector) Detect(explicit, container, kubernetesFilter string) (string, error) {
	if explicit != "" {
		return d.validateExplicit(explicit)
	}
	return d.autoDetect(container, kubernetesFilter)
}

// validateExplicit confirms the operator-supplied backend is usable.
func (d *SystemDetector) validateExplicit(bk string) (string, error) {
	switch bk {
	case "podman", "podman-remote", "kubectl":
		if !d.Path.LookPath(bk) {
			return "", fmt.Errorf(L("backend command not found in PATH: %s"), bk)
		}
		return bk, nil
	case "host":
		return "host", nil
	default:
		return "", fmt.Errorf(L("unsupported backend %s"), bk)
	}
}

// autoDetect probes the system in documented precedence order.
func (d *SystemDetector) autoDetect(container, kubernetesFilter string) (string, error) {
	hasKubectl := d.Path.LookPath("kubectl")
	hasPodman := false

	// Precedence 1: kubectl with a live deployment.
	if hasKubectl {
		found, err := d.Kubernetes.HasDeployment(kubernetesFilter)
		if err != nil {
			log.Info().Msg(L("kubectl not configured to connect to a cluster, ignoring"))
		} else if found {
			return "kubectl", nil
		}
	}

	// Precedence 2: podman / podman-remote with the container already running.
	for _, bin := range []string{"podman", "podman-remote"} {
		if d.Path.LookPath(bin) {
			hasPodman = true
			if d.Containers.HasRunningContainer(bin, container) {
				return bin, nil
			}
		}
	}

	// Precedence 3: podman present and a uyuni systemd service is installed.
	if hasPodman && d.Systemd != nil {
		if d.Systemd.HasService("uyuni-server") || d.Systemd.HasService("uyuni-proxy-pod") {
			return "podman", nil
		}
	}

	// Precedence 4: kubectl with a known Helm release.
	if hasKubectl {
		if d.Kubernetes.HasHelmRelease("uyuni") || d.Kubernetes.HasHelmRelease("uyuni-proxy") {
			return "kubectl", nil
		}
	}

	return "", fmt.Errorf(L("uyuni container is not accessible with one of podman, podman-remote or kubectl"))
}

// execPathLookup wraps exec.LookPath.
type execPathLookup struct{}

func (e *execPathLookup) LookPath(file string) bool {
	_, err := exec.LookPath(file)
	return err == nil
}

// podmanProber inspects running containers via podman or podman-remote.
type podmanProber struct {
	runner Runner
}

func (p *podmanProber) HasRunningContainer(bin, container string) bool {
	_, err := p.runner(bin, "inspect", container, "--format", "{{.Name}}").
		Spinner("").Exec()
	return err == nil
}

// kubectlProber detects live Kubernetes workloads and Helm releases.
type kubectlProber struct {
	runner Runner
}

// HasDeployment returns whether at least one deployment matching filter exists.
// An error is returned only when kubectl cannot reach the cluster at all.
func (k *kubectlProber) HasDeployment(filter string) (bool, error) {
	out, err := k.runner("kubectl",
		"--request-timeout=30s", "get", "deploy", filter,
		"-A", "-o=jsonpath={.items[*].metadata.name}",
	).Log(zerolog.DebugLevel).Spinner("").Exec()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(out)) != 0, nil
}

// HasHelmRelease returns whether a Helm release with the given name exists.
// Delegates to kubernetes.HasHelmRelease so k3s kubeconfig handling is consistent
// with the rest of the project.
func (k *kubectlProber) HasHelmRelease(release string) bool {
	if !utils.IsInstalled("helm") {
		return false
	}
	kubeconfig := ""
	if infos, err := kubernetes.CheckCluster(); err == nil {
		kubeconfig = infos.GetKubeconfig()
	}
	return kubernetes.HasHelmRelease(release, kubeconfig)
}