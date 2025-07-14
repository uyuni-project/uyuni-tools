// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"
)

type TestCase struct {
	Input    ImageFlags
	Expected string
}

func TestRegistryFQDN(t *testing.T) {
	data := []TestCase{
		{
			Input:    ImageFlags{Registry: `registry.opensuse.org/uyuni`},
			Expected: `registry.opensuse.org`,
		},
		{
			Input:    ImageFlags{Registry: `registry.suse.com`},
			Expected: `registry.suse.com`,
		},
		{
			Input: ImageFlags{
				Registry: `updates.suse.com/SUSE/Updates/SLE-Product-SUSE-Manager-Proxy/4.3-LTS/x86_64/update`},
			Expected: `updates.suse.com`,
		},
		{
			Input: ImageFlags{
				Registry: `https://updates.suse.com/SUSE/Updates/SLE-Product-SUSE-Manager-Proxy/4.3-LTS/x86_64/update`},
			Expected: `https://updates.suse.com`,
		},
	}
	for i, testCase := range data {
		actual := testCase.Input.ComputeRegistryFQDN()

		if actual != testCase.Expected {
			t.Errorf("Testcase %d: Expected %s got %s when registry %s", i, testCase.Expected, actual, testCase.Input.Registry)
		}
	}
}
