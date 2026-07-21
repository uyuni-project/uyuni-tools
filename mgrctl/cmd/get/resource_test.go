// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestResourceTypesNoDuplicates(t *testing.T) {
	var seen []string

	for name, res := range resourceTypes {
		if utils.Contains(seen, name) {
			t.Errorf("Duplicate resource key found: %s", name)
		}
		seen = append(seen, name)

		for _, alias := range res.Aliases {
			if utils.Contains(seen, alias) {
				t.Errorf("Duplicate resource alias found: %s (in resource %s)", alias, name)
			}
			seen = append(seen, alias)
		}
	}
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		input     string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{"name=foo", "name", "=foo", false},
		{"extra_pkg_count>0", "extra_pkg_count", ">0", false},
		{"count>=10", "count", ">=10", false},
		{"count<=5", "count", "<=5", false},
		{"status!=active", "status", "!=active", false},
		{"id<100", "id", "<100", false},
		{"justkey", "justkey", "", false},
		{" name = foo ", "name", "= foo", false},
		{"mykey |= my value", "", "", true},
		{"bad key=foo", "", "", true},
		{"=foo", "", "", true},
	}

	for _, tt := range tests {
		key, value, err := parseFilter(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseFilter(%q): got err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && key != tt.wantKey {
			t.Errorf("parseFilter(%q): got key %q, want %q", tt.input, key, tt.wantKey)
		}
		if !tt.wantErr && value != tt.wantValue {
			t.Errorf("parseFilter(%q): got value %q, want %q", tt.input, value, tt.wantValue)
		}
	}
}
