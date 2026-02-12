// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

var allServices = []string{
	podman.ServerService,
	podman.DBService,
	podman.SalineService + "@",
	podman.ServerAttestationService + "@",
	podman.HubXmlrpcService + "@",
}

func TestStartServices(t *testing.T) {
	cases := []struct {
		installed          []string
		enabled            []string
		expectedStarted    []string
		expectedNotStarted []string
		startErrors        map[string]error
		err                error
	}{
		// Regular case with only server and DB containers
		{
			installed:       allServices,
			enabled:         []string{podman.ServerService, podman.DBService},
			expectedStarted: []string{podman.ServerService, podman.DBService},
			expectedNotStarted: []string{
				podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@",
				podman.SalineService + "@", podman.TFTPService,
			},
		},
		// Regular case with an instance of all services.
		{
			installed: allServices,
			enabled: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			expectedStarted: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			expectedNotStarted: []string{},
		},
		// In a migration from non-split DB to split DB we have no DB container yet
		{
			installed: []string{
				podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@",
			},
			enabled:         []string{podman.ServerService},
			expectedStarted: []string{podman.ServerService},
			expectedNotStarted: []string{
				podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@", podman.DBService, podman.TFTPService,
			},
		},
		// Error case where both the server and the DB service fail to start
		{
			installed: allServices,
			enabled:   []string{podman.ServerService, podman.DBService},
			expectedNotStarted: []string{
				podman.ServerService, podman.DBService, podman.HubXmlrpcService + "@",
				podman.ServerAttestationService + "@", podman.SalineService + "@", podman.TFTPService,
			},
			startErrors: map[string]error{
				podman.ServerService: errors.New("failed to start server"),
				podman.DBService:     errors.New("failed to start DB"),
			},
			err: errors.New("failed to start DB; failed to start server"),
		},
	}

	for i, testCase := range cases {
		driver := testutils.FakeSystemdDriver{
			Installed:          testCase.installed,
			Enabled:            testCase.enabled,
			StartServiceErrors: testCase.startErrors,
		}

		systemd = podman.NewSystemdWithDriver(&driver)

		err := StartServices()

		prefix := fmt.Sprintf("case %d - ", i+1)
		for _, service := range testCase.expectedStarted {
			testutils.AssertContains(t, fmt.Sprintf("%s%s not started", prefix, service), driver.Running, service)
		}
		for _, service := range testCase.expectedNotStarted {
			testutils.AssertNotContains(t, fmt.Sprintf("%s%s has been started", prefix, service), driver.Running, service)
		}
		testutils.AssertEquals(t, prefix+"unexpected error returned", testCase.err, err)
	}
}

func TestStopServices(t *testing.T) {
	cases := []struct {
		installed          []string
		enabled            []string
		started            []string
		expectedStarted    []string
		expectedNotStarted []string
		stopErrors         map[string]error
		err                error
	}{
		// Regular case with only server and DB containers
		{
			installed: allServices,
			enabled:   []string{podman.ServerService, podman.DBService},
			started:   []string{podman.ServerService, podman.DBService},
			expectedNotStarted: []string{
				podman.ServerService, podman.DBService, podman.HubXmlrpcService + "@",
				podman.ServerAttestationService + "@", podman.SalineService + "@", podman.TFTPService,
			},
		},
		// Regular case with an instance of all services.
		{
			installed: allServices,
			enabled: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			started: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			expectedNotStarted: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
		},
		// In a migration from non-split DB to split DB we have no DB container yet
		{
			installed: []string{podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@"},
			enabled:   []string{podman.ServerService},
			started:   []string{podman.ServerService},
			expectedNotStarted: []string{
				podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@", podman.DBService,
				podman.TFTPService,
			},
		},
		// Error case where both the server and the DB service fail to start
		{
			installed:       allServices,
			enabled:         []string{podman.ServerService, podman.DBService},
			started:         []string{podman.ServerService, podman.DBService},
			expectedStarted: []string{podman.ServerService, podman.DBService},
			expectedNotStarted: []string{
				podman.HubXmlrpcService + "@",
				podman.ServerAttestationService + "@", podman.SalineService + "@", podman.TFTPService,
			},
			stopErrors: map[string]error{
				podman.ServerService: errors.New("failed to stop server"),
				podman.DBService:     errors.New("failed to stop DB"),
			},
			err: errors.New("failed to stop server; failed to stop DB"),
		},
	}

	for i, testCase := range cases {
		driver := testutils.FakeSystemdDriver{
			Installed:         testCase.installed,
			Enabled:           testCase.enabled,
			Running:           testCase.started,
			StopServiceErrors: testCase.stopErrors,
		}

		systemd = podman.NewSystemdWithDriver(&driver)

		err := StopServices()

		prefix := fmt.Sprintf("case %d - ", i+1)
		for _, service := range testCase.expectedStarted {
			testutils.AssertContains(t, fmt.Sprintf("%s%s has been stopped", prefix, service), driver.Running, service)
		}
		for _, service := range testCase.expectedNotStarted {
			testutils.AssertNotContains(t, fmt.Sprintf("%s%s has not been stopped", prefix, service), driver.Running, service)
		}
		testutils.AssertEquals(t, prefix+"unexpected error returned", testCase.err, err)
	}
}
