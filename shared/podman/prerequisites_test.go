// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestCheckPodmanRunningContainers(t *testing.T) {
	type testCase struct {
		name    string
		output  string
		err     error
		wantErr bool
	}

	tests := []testCase{
		{
			name:    "no containers running",
			output:  "",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "containers running",
			output:  "abc123\n",
			err:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldRunner := runner
			defer func() { runner = oldRunner }()

			runner = testutils.FakeRunnerGenerator(tt.output, tt.err)

			err := CheckPodmanRunningContainers()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPodmanRunningContainers() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
