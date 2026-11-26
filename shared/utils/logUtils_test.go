// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"
)

func TestRedact(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "reportdb password in log file",
			input: `{"level":"info","time":"2024-04-29T15:23:39+02:00",
				"message":"Running /usr/bin/uyuni-setup-reportdb create ` +
				`--db reportdb --user pythia_susemanager --host localhost --address * --remote 0.0.0.0/0,::/0 --password ` +
				`/z4FffHC2HxaagBeIXFzshxtNUfbqm5Zwv/EgvxT"}`,
			expected: `{"level":"info","time":"2024-04-29T15:23:39+02:00",
				"message":"Running /usr/bin/uyuni-setup-reportdb create ` +
				`--db reportdb --user pythia_susemanager --host localhost --address * --remote 0.0.0.0/0,::/0 --password ` +
				`<REDACTED>"}`,
		},
		{
			name: "reportdb password in log output",
			input: `Running /usr/bin/uyuni-setup-reportdb create --db reportdb --user pythia_susemanager --host localhost ` +
				`--address * --remote 0.0.0.0/0,::/0 --password iVgQsuPDGxwKhFc5bfk4IjpVBbqrbyRDYKEsww+Y`,
			expected: `Running /usr/bin/uyuni-setup-reportdb create --db reportdb --user pythia_susemanager --host localhost ` +
				`--address * --remote 0.0.0.0/0,::/0 --password <REDACTED>`,
		},
		{
			name:     "password in json",
			input:    `{"adminLogin":"admin","adminPassword":"secret","email":"no@email.com"}`,
			expected: `{"adminLogin":"admin","adminPassword":"<REDACTED>","email":"no@email.com"}`,
		},
		{
			name:     "word password followed by a newline",
			input:    `password\n`,
			expected: `password\n`,
		},
		{
			name:     "escaped json password",
			input:    `\"password\": \"foo\"`,
			expected: `\"password\": \"<REDACTED>\"`,
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := redact(tt.input)

			if actual != tt.expected {
				t.Errorf("Expected %s got %s when redacting %s", tt.expected, actual, tt.input)
			}
		})
	}
}
