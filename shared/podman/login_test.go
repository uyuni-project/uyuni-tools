// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestGeneratePodmanLoginContent(t *testing.T) {
	tests := []struct {
		hostData *HostInspectData
		registry types.Registry
		scc      types.SCCCredentials
		expected string
	}{
		// host registered, registry == registry.suse.com, both scc and registry credentials.
		// in this case, since registry.suse.com, scc credentials flags are used.
		{
			hostData: &HostInspectData{
				SCCUsername: "sccuserhost",
				SCCPassword: "sccpasswordhost",
			},
			registry: types.Registry{
				Host:     "registry.suse.com/suse/some/paths/",
				User:     "registryuserflag",
				Password: "registrypasswordflag",
			},
			scc: types.SCCCredentials{
				User:     "sccuserflag",
				Password: "sccpasswordflag",
			},
			expected: fmt.Sprintf(`{
	"auths": {
		"registry.suse.com/suse/some/paths/" : {
			"auth": "%s"
		}
	}
}`, base64.StdEncoding.EncodeToString([]byte("sccuserflag:sccpasswordflag"))),
		},
		// host registered, registry != registry.suse.com, both scc and registry credentials.
		// in this case, since != registry.suse.com, registry credentials flags are used.
		{
			hostData: &HostInspectData{
				SCCUsername: "sccuserhost",
				SCCPassword: "sccpasswordhost",
			},
			registry: types.Registry{
				Host:     "myregistry.com/suse/some/paths/",
				User:     "registryuserflag",
				Password: "registrypasswordflag",
			},
			scc: types.SCCCredentials{
				User:     "sccuserflag",
				Password: "sccpasswordflag",
			},
			expected: fmt.Sprintf(`{
	"auths": {
		"myregistry.com/suse/some/paths/" : {
			"auth": "%s"
		}
	}
}`, base64.StdEncoding.EncodeToString([]byte("registryuserflag:registrypasswordflag"))),
		},
		// host registered, registry == registry.suse.com, no flag credentials.
		// in this case, since == registry.suse.com, host credentials are used.
		{
			hostData: &HostInspectData{
				SCCUsername: "sccuserhost",
				SCCPassword: "sccpasswordhost",
			},
			registry: types.Registry{
				Host:     "registry.suse.com",
				User:     "registryuserflag",
				Password: "registrypasswordflag",
			},
			scc: types.SCCCredentials{
				User:     "",
				Password: "",
			},
			expected: fmt.Sprintf(`{
	"auths": {
		"registry.suse.com" : {
			"auth": "%s"
		}
	}
}`, base64.StdEncoding.EncodeToString([]byte("sccuserhost:sccpasswordhost"))),
		},
		// just registry != registry.suse.com, no credentials: no auth file.
		//
		{
			hostData: &HostInspectData{
				SCCUsername: "",
				SCCPassword: "",
			},
			registry: types.Registry{
				Host:     "myregistry.com",
				User:     "",
				Password: "",
			},
			scc: types.SCCCredentials{
				User:     "",
				Password: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		actual := GeneratePodmanLoginContent(tt.hostData, tt.registry, tt.scc)
		if actual != tt.expected {
			t.Errorf("PodmanLogin error:\n actual:\n %s \n expected:\n %s", actual, tt.expected)
		}
	}
}
