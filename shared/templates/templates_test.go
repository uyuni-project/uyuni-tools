// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// templateRenderer is a local interface to avoid importing shared/utils,
// which would cause an import cycle.
type templateRenderer interface {
	Render(wr io.Writer) error
}

func TestTemplatesRender(t *testing.T) {
	tests := []struct {
		name     string
		template templateRenderer
	}{
		{
			name: "InspectTemplateData",
			template: InspectTemplateData{
				Values: []types.InspectData{
					types.NewInspectData("TEST_VAR", "echo 'hello'"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.template.Render(io.Discard); err != nil {
				t.Errorf("%s render failed: %v", tt.name, err)
			}
		})
	}
}
