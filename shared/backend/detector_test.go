// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package backend_test

import (
	"errors"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/backend"
)

// fakePath records which binaries are "in PATH".
type fakePath struct {
	installed map[string]bool
}

func (f *fakePath) LookPath(file string) bool {
	return f.installed[file]
}

// fakeContainerProber records which (bin, container) pairs are running.
type fakeContainerProber struct {
	running map[string]bool
}

func (f *fakeContainerProber) HasRunningContainer(bin, container string) bool {
	return f.running[bin+"/"+container]
}

// fakeSystemd records which service names are installed.
type fakeSystemd struct {
	installed map[string]bool
}

func (f *fakeSystemd) HasService(name string) bool {
	return f.installed[name]
}

// fakeKubernetes controls HasDeployment and HasHelmRelease results.
type fakeKubernetes struct {
	// deploymentFound and deploymentErr control the HasDeployment return values.
	deploymentFound map[string]bool
	deploymentErr   map[string]error
	helmReleases    map[string]bool
}

func (f *fakeKubernetes) HasDeployment(filter string) (bool, error) {
	if err, ok := f.deploymentErr[filter]; ok && err != nil {
		return false, err
	}
	return f.deploymentFound[filter], nil
}

func (f *fakeKubernetes) HasHelmRelease(release string) bool {
	return f.helmReleases[release]
}

// makeDetector builds a SystemDetector from fakes.
func makeDetector(
	path *fakePath,
	containers *fakeContainerProber,
	systemd *fakeSystemd,
	kubernetes *fakeKubernetes,
) *backend.SystemDetector {
	d := &backend.SystemDetector{
		Path:       path,
		Containers: containers,
		Kubernetes: kubernetes,
	}
	if systemd != nil {
		d.Systemd = systemd
	}
	return d
}

func TestDetect_ExplicitBackend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		explicit    string
		pathHas     []string
		wantCommand string
		wantErr     bool
	}{
		{
			name:        "explicit podman found in PATH",
			explicit:    "podman",
			pathHas:     []string{"podman"},
			wantCommand: "podman",
		},
		{
			name:        "explicit podman-remote found in PATH",
			explicit:    "podman-remote",
			pathHas:     []string{"podman-remote"},
			wantCommand: "podman-remote",
		},
		{
			name:        "explicit kubectl found in PATH",
			explicit:    "kubectl",
			pathHas:     []string{"kubectl"},
			wantCommand: "kubectl",
		},
		{
			name:        "host never needs a PATH lookup",
			explicit:    "host",
			pathHas:     []string{},
			wantCommand: "host",
		},
		{
			name:     "explicit podman not in PATH",
			explicit: "podman",
			pathHas:  []string{},
			wantErr:  true,
		},
		{
			name:     "explicit kubectl not in PATH",
			explicit: "kubectl",
			pathHas:  []string{},
			wantErr:  true,
		},
		{
			name:     "explicit podman-remote not in PATH",
			explicit: "podman-remote",
			pathHas:  []string{},
			wantErr:  true,
		},
		{
			name:     "unsupported backend name",
			explicit: "docker",
			pathHas:  []string{"docker"},
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			installed := make(map[string]bool)
			for _, b := range tc.pathHas {
				installed[b] = true
			}
			d := makeDetector(
				&fakePath{installed: installed},
				&fakeContainerProber{running: map[string]bool{}},
				nil,
				&fakeKubernetes{},
			)

			got, err := d.Detect(tc.explicit, "uyuni-server", "-lapp=uyuni")

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error, got command=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantCommand {
				t.Errorf("got %q, want %q", got, tc.wantCommand)
			}
		})
	}
}

func TestDetect_AutoDetect_Precedence(t *testing.T) {
	t.Parallel()

	const (
		container = "uyuni-server"
		filter    = "-lapp=uyuni"
	)

	errUnreachable := errors.New("connection refused")

	tests := []struct {
		name            string
		pathHas         []string
		runningBin      string
		systemdHas      []string
		deploymentFound bool
		deploymentErr   error
		helmReleases    []string
		wantCommand     string
		wantErr         bool
	}{
		{
			name:            "kubectl live deployment wins over everything",
			pathHas:         []string{"kubectl", "podman"},
			runningBin:      "podman",
			systemdHas:      []string{"uyuni-server"},
			deploymentFound: true,
			helmReleases:    []string{"uyuni"},
			wantCommand:     "kubectl",
		},
		{
			name:          "kubectl unreachable falls through to podman",
			pathHas:       []string{"kubectl", "podman"},
			runningBin:    "podman",
			deploymentErr: errUnreachable,
			wantCommand:   "podman",
		},
		{
			name:        "podman container running",
			pathHas:     []string{"podman", "kubectl"},
			runningBin:  "podman",
			wantCommand: "podman",
		},
		{
			name:        "podman-remote container running",
			pathHas:     []string{"podman-remote"},
			runningBin:  "podman-remote",
			wantCommand: "podman-remote",
		},
		{
			name:        "podman checked before podman-remote",
			pathHas:     []string{"podman", "podman-remote"},
			runningBin:  "podman",
			wantCommand: "podman",
		},
		{
			name:        "podman installed + uyuni-server service present",
			pathHas:     []string{"podman"},
			systemdHas:  []string{"uyuni-server"},
			wantCommand: "podman",
		},
		{
			name:        "podman installed + uyuni-proxy-pod service present",
			pathHas:     []string{"podman"},
			systemdHas:  []string{"uyuni-proxy-pod"},
			wantCommand: "podman",
		},
		{
			name:         "kubectl + uyuni helm release",
			pathHas:      []string{"kubectl"},
			helmReleases: []string{"uyuni"},
			wantCommand:  "kubectl",
		},
		{
			name:         "kubectl + uyuni-proxy helm release",
			pathHas:      []string{"kubectl"},
			helmReleases: []string{"uyuni-proxy"},
			wantCommand:  "kubectl",
		},
		{
			name:    "nothing installed",
			pathHas: []string{},
			wantErr: true,
		},
		{
			name:    "podman installed but no container running and no service",
			pathHas: []string{"podman"},
			wantErr: true,
		},
		{
			name:    "kubectl installed but no deployment and no helm release",
			pathHas: []string{"kubectl"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			installed := make(map[string]bool)
			for _, b := range tc.pathHas {
				installed[b] = true
			}

			running := make(map[string]bool)
			if tc.runningBin != "" {
				running[tc.runningBin+"/"+container] = true
			}

			systemd := &fakeSystemd{installed: make(map[string]bool)}
			for _, svc := range tc.systemdHas {
				systemd.installed[svc] = true
			}

			helmReleases := make(map[string]bool)
			for _, rel := range tc.helmReleases {
				helmReleases[rel] = true
			}

			deploymentErr := make(map[string]error)
			if tc.deploymentErr != nil {
				deploymentErr[filter] = tc.deploymentErr
			}

			d := makeDetector(
				&fakePath{installed: installed},
				&fakeContainerProber{running: running},
				systemd,
				&fakeKubernetes{
					deploymentFound: map[string]bool{filter: tc.deploymentFound},
					deploymentErr:   deploymentErr,
					helmReleases:    helmReleases,
				},
			)

			got, err := d.Detect("", container, filter)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error, got command=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantCommand {
				t.Errorf("got %q, want %q", got, tc.wantCommand)
			}
		})
	}
}

func TestDetect_ErrorMessages(t *testing.T) {
	t.Parallel()

	t.Run("error message includes the unsupported backend name", func(t *testing.T) {
		t.Parallel()
		d := makeDetector(
			&fakePath{installed: map[string]bool{}},
			&fakeContainerProber{running: map[string]bool{}},
			nil,
			&fakeKubernetes{},
		)
		_, err := d.Detect("docker", "uyuni-server", "")
		if err == nil {
			t.Fatal("expected error")
		}
		if !containsSub(err.Error(), "docker") {
			t.Errorf("error %q should mention the backend name", err.Error())
		}
	})

	t.Run("error message includes the missing binary name", func(t *testing.T) {
		t.Parallel()
		d := makeDetector(
			&fakePath{installed: map[string]bool{}},
			&fakeContainerProber{running: map[string]bool{}},
			nil,
			&fakeKubernetes{},
		)
		_, err := d.Detect("kubectl", "uyuni-server", "")
		if err == nil {
			t.Fatal("expected error")
		}
		if !containsSub(err.Error(), "kubectl") {
			t.Errorf("error %q should mention the binary name", err.Error())
		}
	})
}

func TestDetect_Idempotent(t *testing.T) {
	t.Parallel()

	d := makeDetector(
		&fakePath{installed: map[string]bool{"podman": true}},
		&fakeContainerProber{running: map[string]bool{"podman/uyuni-server": true}},
		&fakeSystemd{installed: map[string]bool{}},
		&fakeKubernetes{},
	)

	first, err := d.Detect("", "uyuni-server", "")
	if err != nil {
		t.Fatalf("first Detect: %v", err)
	}
	second, err := d.Detect("", "uyuni-server", "")
	if err != nil {
		t.Fatalf("second Detect: %v", err)
	}
	if first != second {
		t.Errorf("not idempotent: first=%q second=%q", first, second)
	}
}

func containsSub(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
