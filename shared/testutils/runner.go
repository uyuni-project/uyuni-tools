// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"bytes"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"golang.org/x/net/context"
)

type fakeRunner struct {
	out []byte
	err error
}

func (r fakeRunner) Log(_ zerolog.Level) types.Runner {
	return r
}

func (r fakeRunner) Spinner(_ string) types.Runner {
	return r
}

func (r fakeRunner) StdMapping() types.Runner {
	return r
}

func (r fakeRunner) Std(buf *bytes.Buffer) types.Runner {
	buf.Write(r.out)
	return r
}

func (r fakeRunner) InputString(_ string) types.Runner {
	return r
}

func (r fakeRunner) Env(_ []string) types.Runner {
	return r
}

func (r fakeRunner) Exec() ([]byte, error) {
	return r.out, r.err
}

func (r fakeRunner) Start() error {
	return r.err
}

func (r fakeRunner) Wait() error {
	return r.err
}

// FakeRunnerGenerator creates NewRunner function generating a FakeRunner.
// out and err are the returns of the mocked Exec().
func FakeRunnerGenerator(out string, err error) func(string, ...string) types.Runner {
	return func(_ string, _ ...string) types.Runner {
		runner := fakeRunner{
			out: []byte(out),
			err: err,
		}
		return runner
	}
}

// FakeContextRunnerGenerator creates NewContextRunner function generating a FakeRunner.
// out is returned via Std mocking to the buffer. To be used with Start() and Wait().
func FakeContextRunnerGenerator(out string, err error) func(context.Context, string, ...string) types.Runner {
	return func(_ context.Context, _ string, _ ...string) types.Runner {
		runner := fakeRunner{
			out: []byte(out),
			err: err,
		}
		return runner
	}
}
