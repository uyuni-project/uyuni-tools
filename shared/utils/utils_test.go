// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"testing"

	expect "github.com/Netflix/go-expect"
)

type askTestData struct {
	value           string
	expectedMessage string
	min             int
	max             int
	checker         func(string) bool
}

func TestAskIfMissing(t *testing.T) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		t.Errorf("Failed to create fake console")
	}
	defer c.Close()

	origStdin := os.Stdin
	origStdout := os.Stdout

	os.Stdin = c.Tty()
	os.Stdout = c.Tty()
	defer func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
	}()

	fChecker := func(v string) bool {
		if !strings.Contains(v, "f") {
			fmt.Println("Has to contain an 'f'")
			return false
		}
		return true
	}

	data := []askTestData{
		{value: "\n", expectedMessage: "A value is required", min: 1, max: 5, checker: nil},
		{value: "superlong\n", expectedMessage: "Has to be less than 5 characters long", min: 1, max: 5, checker: nil},
		{value: "a\n", expectedMessage: "Has to be more than 2 characters long", min: 2, max: 5, checker: nil},
		{value: "booh\n", expectedMessage: "Has to contain an 'f'", min: 0, max: 0, checker: fChecker},
	}

	for i, testCase := range data {
		go func() {
			if _, err := c.ExpectString("Prompted value: "); err != nil {
				t.Errorf("Testcase %d: Expected prompt error: %s", i, err)
			}
			if _, err := c.Send(testCase.value); err != nil {
				t.Errorf("Testcase %d: Failed to send value to fake console: %s", i, err)
			}
			if _, err := c.Expect(expect.Regexp(regexp.MustCompile(testCase.expectedMessage))); err != nil {
				t.Errorf("Testcase %d: Expected '%s' message: %s", i, testCase.expectedMessage, err)
			}
			if _, err := c.ExpectString("Prompted value: "); err != nil {
				t.Errorf("Testcase %d: Expected prompt error: %s", i, err)
			}
			if _, err := c.Send("foo\n"); err != nil {
				t.Errorf("Testcase %d: Failed to send value to fake console: %s", i, err)
			}
		}()

		var value string
		AskIfMissing(&value, "Prompted value", testCase.min, testCase.max, testCase.checker)
		if value != "foo" {
			t.Errorf("Testcase %d: Expected 'foo', got '%s' value", i, value)
		}
	}
}

func TestAskPasswordIfMissing(t *testing.T) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		t.Errorf("Failed to create fake console")
	}
	defer c.Close()

	origStdin := syscall.Stdin
	origStdout := os.Stdout

	syscall.Stdin = int(c.Tty().Fd())
	os.Stdout = c.Tty()
	defer func() {
		syscall.Stdin = origStdin
		os.Stdout = origStdout
	}()

	data := []askTestData{
		{value: "\n", expectedMessage: "A value is required", min: 1, max: 5, checker: nil},
		{value: "superlong\n", expectedMessage: "Has to be less than 5 characters long", min: 1, max: 5, checker: nil},
		{value: "a\n", expectedMessage: "Has to be more than 2 characters long", min: 2, max: 5, checker: nil},
	}

	for i, testCase := range data {
		go func() {
			if _, err := c.ExpectString("Prompted password: "); err != nil {
				t.Errorf("Testcase %d: Expected prompt error: %s", i, err)
			}
			if _, err := c.Send(testCase.value); err != nil {
				t.Errorf("Testcase %d: Failed to send value to fake console: %s", i, err)
			}
			if _, err := c.Expect(expect.Regexp(regexp.MustCompile(testCase.expectedMessage))); err != nil {
				t.Errorf("Testcase %d: Expected '%s' message: %s", i, testCase.expectedMessage, err)
			}
			if _, err := c.ExpectString("Prompted password: "); err != nil {
				t.Errorf("Testcase %d: Expected prompt error: %s", i, err)
			}
			if _, err := c.Send("foo\n"); err != nil {
				t.Errorf("Testcase %d: Failed to send value to fake console: %s", i, err)
			}
		}()

		var value string
		AskPasswordIfMissing(&value, "Prompted password", testCase.min, testCase.max)
		if value != "foo" {
			t.Errorf("Expected 'foo', got '%s' value", value)
		}
	}
}

func TestComputeImage(t *testing.T) {
	data := [][]string{
		{"registry:5000/path/to/image:foo", "registry:5000/path/to/image:foo", "bar"},
		{"registry:5000/path/to/image:bar", "registry:5000/path/to/image", "bar"},
		{"registry/path/to/image:foo", "registry/path/to/image:foo", "bar"},
		{"registry/path/to/image:bar", "registry/path/to/image", "bar"},
		{"registry:5000/path/to/image-migration-14-16:foo", "registry:5000/path/to/image:foo", "bar", "-migration-14-16"},
		{"registry:5000/path/to/image-migration-14-16:bar", "registry:5000/path/to/image", "bar", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:foo", "registry/path/to/image:foo", "bar", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:bar", "registry/path/to/image", "bar", "-migration-14-16"},
	}

	for i, testCase := range data {
		result := testCase[0]
		image := testCase[1]
		tag := testCase[2]
		appendToImage := testCase[3:]

		actual, err := ComputeImage(image, tag, appendToImage...)

		if err != nil {
			t.Errorf("Testcase %d: Unexpected error while computing image with %s, %s, %s: %s", i, image, tag, appendToImage, err)
		}
		if actual != result {
			t.Errorf("Testcase %d: Expected %s got %s when computing image with %s, %s, %s", i, result, actual, image, tag, appendToImage)
		}
	}
}

func TestComputeImageError(t *testing.T) {
	_, err := ComputeImage("registry:path/to/image:tag:tag", "bar")
	if err == nil {
		t.Error("Expected error, got none")
	}
}
