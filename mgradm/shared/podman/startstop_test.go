// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func TestStartStopWarnWhenServicesDisabled(t *testing.T) {
	driver := testutils.FakeSystemdDriver{
		Installed: []string{podman.ServerService, podman.DBService},
	}
	systemd = podman.NewSystemdWithDriver(&driver)
	var output bytes.Buffer
	oldLogger := log.Logger
	log.Logger = zerolog.New(&output)
	t.Cleanup(func() {
		log.Logger = oldLogger
	})

	testutils.AssertNoError(t, "start failed: ", StartServices())
	testutils.AssertNoError(t, "stop failed: ", StopServices())
	if count := strings.Count(output.String(),
		"Server services are disabled and will not start automatically at boot"); count != 2 {
		t.Errorf("expected start and stop warnings, got %d", count)
	}
	testutils.AssertEquals(t, "start or stop changed boot enablement", 0, len(driver.Enabled))
}

func TestStartServices(t *testing.T) {
	cases := []struct {
		name               string
		installed          []string
		enabled            []string
		expectedStarted    []string
		expectedNotStarted []string
		startErrors        map[string]error
		err                string
	}{
		{
			name:            "only server and DB containers",
			installed:       allServices,
			enabled:         []string{podman.ServerService, podman.DBService},
			expectedStarted: []string{podman.ServerService, podman.DBService},
			expectedNotStarted: []string{
				podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@",
				podman.SalineService + "@", podman.TFTPService,
			},
		},
		{
			name:      "an instance of all services",
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
		{
			name: "disabled services",
			installed: []string{
				podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@",
			},
			enabled:         []string{podman.ServerService},
			expectedStarted: []string{podman.ServerService},
			expectedNotStarted: []string{
				podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@", podman.DBService, podman.TFTPService,
			},
		},
		{
			name:      "error when both server and DB service fail to start",
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
			err: "failed to start DB\nfailed to start server",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			driver := testutils.FakeSystemdDriver{
				Installed:          testCase.installed,
				Enabled:            testCase.enabled,
				StartServiceErrors: testCase.startErrors,
			}

			systemd = podman.NewSystemdWithDriver(&driver)

			err := StartServices()

			for _, service := range testCase.expectedStarted {
				testutils.AssertContains(t, fmt.Sprintf("%s not started", service), driver.Running, service)
			}
			for _, service := range testCase.expectedNotStarted {
				testutils.AssertNotContains(t, fmt.Sprintf("%s has been started", service), driver.Running, service)
			}
			testutils.AssertError(t, testCase.err, err)
		})
	}
}

func TestStopServices(t *testing.T) {
	cases := []struct {
		name               string
		installed          []string
		enabled            []string
		started            []string
		expectedStopped    []string
		expectedNotStopped []string
		stopErrors         map[string]error
		err                string
	}{
		{
			name:      "only server and DB containers",
			installed: allServices,
			enabled:   []string{podman.ServerService, podman.DBService},
			started:   []string{podman.ServerService, podman.DBService},
			expectedNotStopped: []string{
				podman.ServerService, podman.DBService, podman.HubXmlrpcService + "@",
				podman.ServerAttestationService + "@", podman.SalineService + "@", podman.TFTPService,
			},
		},
		{
			name:      "an instance of all services",
			installed: allServices,
			enabled: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			started: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
			expectedNotStopped: []string{
				podman.ServerService, podman.DBService, podman.TFTPService,
				podman.HubXmlrpcService + "@0", podman.ServerAttestationService + "@0", podman.SalineService + "@0",
			},
		},
		{
			name:      "disabled services",
			installed: []string{podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@"},
			enabled:   []string{podman.ServerService},
			started:   []string{podman.ServerService},
			expectedNotStopped: []string{
				podman.ServerService, podman.HubXmlrpcService + "@", podman.ServerAttestationService + "@", podman.DBService,
				podman.TFTPService,
			},
		},
		{
			name:            "error when both server and DB service fail to stop",
			installed:       allServices,
			enabled:         []string{podman.ServerService, podman.DBService},
			started:         []string{podman.ServerService, podman.DBService},
			expectedStopped: []string{podman.ServerService, podman.DBService},
			expectedNotStopped: []string{
				podman.HubXmlrpcService + "@",
				podman.ServerAttestationService + "@", podman.SalineService + "@", podman.TFTPService,
			},
			stopErrors: map[string]error{
				podman.ServerService: errors.New("failed to stop server"),
				podman.DBService:     errors.New("failed to stop DB"),
			},
			err: "failed to stop server\nfailed to stop DB",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			driver := testutils.FakeSystemdDriver{
				Installed:         testCase.installed,
				Enabled:           testCase.enabled,
				Running:           testCase.started,
				StopServiceErrors: testCase.stopErrors,
			}

			systemd = podman.NewSystemdWithDriver(&driver)

			err := StopServices()

			for _, service := range testCase.expectedStopped {
				testutils.AssertContains(t, fmt.Sprintf("%s has been stopped", service), driver.Running, service)
			}
			for _, service := range testCase.expectedNotStopped {
				testutils.AssertNotContains(t, fmt.Sprintf("%s has not been stopped", service), driver.Running, service)
			}
			testutils.AssertError(t, testCase.err, err)
		})
	}
}
