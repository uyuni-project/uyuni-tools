// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"
)

func TestRedact(t *testing.T) {
	data := [][]string{
		{
			`{"level":"info","time":"2024-04-29T15:23:39+02:00","message":"Running /usr/bin/uyuni-setup-reportdb create ` +
				`--db reportdb --user pythia_susemanager --host localhost --address * --remote 0.0.0.0/0,::/0 --password ` +
				`/z4FffHC2HxaagBeIXFzshxtNUfbqm5Zwv/EgvxT"}`,
			`{"level":"info","time":"2024-04-29T15:23:39+02:00","message":"Running /usr/bin/uyuni-setup-reportdb create ` +
				`--db reportdb --user pythia_susemanager --host localhost --address * --remote 0.0.0.0/0,::/0 --password ` +
				`<REDACTED>"}`,
		},
		{
			`Running /usr/bin/uyuni-setup-reportdb create --db reportdb --user pythia_susemanager --host localhost --address *` +
				` --remote 0.0.0.0/0,::/0 --password iVgQsuPDGxwKhFc5bfk4IjpVBbqrbyRDYKEsww+Y`,
			`Running /usr/bin/uyuni-setup-reportdb create --db reportdb --user pythia_susemanager --host localhost --address *` +
				` --remote 0.0.0.0/0,::/0 --password <REDACTED>`,
		},
		{
			`{"adminLogin":"admin","adminPassword":"secret","email":"no@email.com"}`,
			`{"adminLogin":"admin","adminPassword":"<REDACTED>","email":"no@email.com"}`,
		},
	}

	for i, testCase := range data {
		input := testCase[0]
		expected := testCase[1]

		actual := redact(input)

		if actual != expected {
			t.Errorf("Testcase %d: Expected %s got %s when redacting  %s", i, expected, actual, input)
		}
	}
}
