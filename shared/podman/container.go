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
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HealthState is the health state of a container as reported by podman.
type HealthState string

const (
	// HealthStarting means that the health is not final yet: the startup healthcheck
	// has not passed or the regular healthcheck did not report a result yet.
	HealthStarting HealthState = "starting"
	// HealthHealthy means that the healthcheck passed.
	HealthHealthy HealthState = "healthy"
	// HealthUnhealthy means that the healthcheck failed more than its retries.
	HealthUnhealthy HealthState = "unhealthy"
	// HealthNone means that the container has no healthcheck configured.
	HealthNone HealthState = "none"
)

// HealthStatus is the current health of a container.
type HealthStatus struct {
	// State is the health state. Podman reports an empty state before the first
	// healthcheck run: this is mapped to HealthStarting. Any value that is not
	// one of the constants above is kept as is.
	State         HealthState
	FailingStreak int
	// LastOutput is the trimmed output of the last healthcheck run.
	LastOutput string
}

// HealthcheckTimeouts describes how long podman can legitimately take to report a final
// health, based on the healthcheck configuration of the container.
type HealthcheckTimeouts struct {
	// Regular bounds the regular healthcheck: the container is expected to be healthy
	// or unhealthy after it. 0 when there is no bound or no healthcheck.
	Regular time.Duration
	// HasStartup is true when a startup healthcheck (podman --health-startup-cmd) is
	// configured. It has to pass before the regular healthcheck even starts.
	HasStartup bool
	// Startup bounds the startup healthcheck, based on its
	// --health-startup-retries. It is only meaningful when HasStartup is true.
	// 0 means that the startup phase is unbounded: with the podman default
	// --health-startup-retries=0 the container is neither marked unhealthy nor
	// restarted while the startup healthcheck keeps failing, so it can
	// legitimately take arbitrarily long.
	Startup time.Duration
}

// WaitTimeout returns how long a container with the given health can be waited for
// before concluding that it will not become healthy. 0 means that there is no bound.
//
// While the startup healthcheck is deciding on readiness the regular healthcheck does
// not run yet, so its timeout does not apply.
func (c HealthcheckTimeouts) WaitTimeout(state HealthState) time.Duration {
	if c.HasStartup && state == HealthStarting {
		return c.Startup
	}
	return c.Regular
}

// ContainerStatus is the part of the container state relevant to startup checks.
type ContainerStatus struct {
	State        string
	Running      bool
	ExitCode     int
	RestartCount int
	Health       HealthStatus
	Timeouts     HealthcheckTimeouts
}

// inspectHealthcheck is the healthcheck configuration in the podman inspect output.
// Durations are in nanoseconds, as time.Duration.
type inspectHealthcheck struct {
	StartPeriod time.Duration `json:"StartPeriod"`
	Interval    time.Duration `json:"Interval"`
	Timeout     time.Duration `json:"Timeout"`
	Retries     int           `json:"Retries"`
	// StartInterval is the interval of the checks during the start period.
	StartInterval time.Duration `json:"StartInterval"`
}

// regularTimeout computes the time after which podman reports the container unhealthy.
func (h *inspectHealthcheck) regularTimeout() time.Duration {
	if h == nil || (h.StartPeriod == 0 && h.Interval == 0 && h.Timeout == 0) {
		return 0
	}
	retries := h.Retries
	if retries < 1 {
		retries = 1
	}
	return h.StartPeriod + time.Duration(retries+1)*(h.Interval+h.Timeout)
}

// startupTimeout computes the time after which podman reports the container unhealthy
// because of its startup healthcheck, or 0 if the startup phase is unbounded.
// regular is the regular healthcheck configuration and can be nil.
func (h *inspectHealthcheck) startupTimeout(regular *inspectHealthcheck) time.Duration {
	if h.Retries <= 0 {
		return 0
	}
	// The startup healthcheck runs every StartInterval, falling back to Interval
	// on older podman versions.
	interval := h.StartInterval
	if interval == 0 {
		interval = h.Interval
	}
	// Podman runs the startup healthcheck exec with the regular healthcheck timeout,
	// so take the larger of the two.
	timeout := h.Timeout
	if regular != nil && regular.Timeout > timeout {
		timeout = regular.Timeout
	}
	return time.Duration(h.Retries) * (interval + timeout)
}

type containerInspect struct {
	RestartCount int `json:"RestartCount"`
	Config       struct {
		Healthcheck        *inspectHealthcheck `json:"Healthcheck"`
		StartupHealthCheck *inspectHealthcheck `json:"StartupHealthCheck"`
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
		Health:       HealthStatus{State: HealthNone},
		Timeouts:     HealthcheckTimeouts{Regular: inspect.Config.Healthcheck.regularTimeout()},
	}
	if shc := inspect.Config.StartupHealthCheck; shc != nil {
		status.Timeouts.HasStartup = true
		status.Timeouts.Startup = shc.startupTimeout(inspect.Config.Healthcheck)
	}
	if health := inspect.State.Health; health != nil {
		status.Health.State = HealthState(health.Status)
		if status.Health.State == "" {
			status.Health.State = HealthStarting
		}
		status.Health.FailingStreak = health.FailingStreak
		if len(health.Log) > 0 {
			status.Health.LastOutput = strings.TrimSpace(health.Log[len(health.Log)-1].Output)
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
		default:
			// GetContainerEvents only asks podman for the events handled above:
			// anything else is not relevant to the startup.
			log.Debug().Msgf("ignoring unexpected podman event %q of container %s", event.Status, event.ID)
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
