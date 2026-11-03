// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ContainerStatus is the part of the container state relevant to startup checks.
type ContainerStatus struct {
	State              string
	Running            bool
	ExitCode           int
	RestartCount       int
	Health             string // starting, healthy, unhealthy, or none
	FailingStreak      int
	HealthLog          string
	HealthcheckTimeout time.Duration
}

type containerInspect struct {
	RestartCount int `json:"RestartCount"`
	Config       struct {
		Healthcheck *struct {
			StartPeriod time.Duration `json:"StartPeriod"`
			Interval    time.Duration `json:"Interval"`
			Timeout     time.Duration `json:"Timeout"`
			Retries     int           `json:"Retries"`
		} `json:"Healthcheck"`
	} `json:"Config"`
	State struct {
		Status   string `json:"Status"`
		Running  bool   `json:"Running"`
		ExitCode int    `json:"ExitCode"`
		Health   *struct {
			Status        string `json:"Status"`
			FailingStreak int    `json:"FailingStreak"`
			Log           []struct {
				ExitCode int    `json:"ExitCode"`
				Output   string `json:"Output"`
			} `json:"Log"`
		} `json:"Health"`
	} `json:"State"`
}

// GetContainerStatus returns the current state of the container.
func GetContainerStatus(command string, container string) (ContainerStatus, error) {
	out, err := runner(command, "inspect", "--format", "{{json .}}", container).
		Log(zerolog.DebugLevel).Spinner("").Exec()
	if err != nil {
		return ContainerStatus{}, utils.Errorf(err, L("cannot inspect container %s"), container)
	}

	var inspect containerInspect
	if err := json.Unmarshal(out, &inspect); err != nil {
		return ContainerStatus{}, utils.Errorf(err, L("failed to parse the state of container %s"), container)
	}

	status := ContainerStatus{
		State:        inspect.State.Status,
		Running:      inspect.State.Running,
		ExitCode:     inspect.State.ExitCode,
		RestartCount: inspect.RestartCount,
		Health:       "none",
	}
	if hc := inspect.Config.Healthcheck; hc != nil {
		if hc.StartPeriod != 0 || hc.Interval != 0 || hc.Timeout != 0 {
			retries := hc.Retries
			if retries < 1 {
				retries = 1
			}
			status.HealthcheckTimeout = hc.StartPeriod + time.Duration(retries+1)*(hc.Interval+hc.Timeout)
		}
	}
	if health := inspect.State.Health; health != nil {
		status.Health = health.Status
		status.FailingStreak = health.FailingStreak
		if len(health.Log) > 0 {
			status.HealthLog = strings.TrimSpace(health.Log[len(health.Log)-1].Output)
		}
	}

	return status, nil
}

// ContainerEvent is a podman container event used to track the startup.
type ContainerEvent struct {
	ID                string `json:"ID"`
	Status            string `json:"Status"`
	TimeNano          int64  `json:"timeNano"`
	ContainerExitCode int    `json:"ContainerExitCode"`
}

// GetContainerEvents returns the container lifecycle events since the given time.
func GetContainerEvents(command string, container string, since time.Time) ([]ContainerEvent, error) {
	// --since only has a second granularity, filter precisely on TimeNano below.
	out, err := runner(command, "events",
		"--since", since.Add(-time.Second).UTC().Format(time.RFC3339),
		"--stream=false",
		"--filter", "container="+container,
		"--filter", "event=create",
		"--filter", "event=start",
		"--filter", "event=died",
		"--filter", "event=remove",
		"--format", "{{json .}}",
	).Log(zerolog.DebugLevel).Spinner("").Exec()
	if err != nil {
		return nil, err
	}

	var events []ContainerEvent
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var event ContainerEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, utils.Errorf(err, L("failed to parse a podman event"))
		}
		if event.TimeNano < since.UnixNano() {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

// ContainerStartState summarizes the lifecycle events of a single start attempt.
type ContainerStartState struct {
	Created    bool
	Started    bool
	FirstStart time.Time
}

// AnalyzeContainerEvents checks the lifecycle events of a single start attempt.
func AnalyzeContainerEvents(events []ContainerEvent) (ContainerStartState, error) {
	state := ContainerStartState{}
	creates := map[string]struct{}{}
	known := map[string]struct{}{}
	starts := 0

	for _, event := range events {
		switch event.Status {
		case "create":
			creates[event.ID] = struct{}{}
			known[event.ID] = struct{}{}
		case "start":
			starts++
			known[event.ID] = struct{}{}
			if state.FirstStart.IsZero() || event.TimeNano < state.FirstStart.UnixNano() {
				state.FirstStart = time.Unix(0, event.TimeNano)
			}
		case "died":
			// Deaths of containers started before the event window belong to the
			// previous instance being replaced during a restart, ignore them.
			if _, ok := known[event.ID]; ok {
				return state, fmt.Errorf(L("container exited during startup with exit code %d"), event.ContainerExitCode)
			}
		case "remove":
			// A remove event without a start is the signature of a failed container creation.
			if _, created := creates[event.ID]; created && starts == 0 {
				return state, errors.New(L("container was removed before it started"))
			}
		}
	}

	state.Created = len(creates) > 0 || starts > 0
	state.Started = starts > 0

	if len(creates) > 1 || starts > 1 {
		return state, errors.New(
			L("container was started more than once during startup, this indicates a crash or an unexpected restart"),
		)
	}
	return state, nil
}
