// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"regexp"
	"syscall"
	"testing"

	expect "github.com/Netflix/go-expect"
)

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

	go func() {
		if _, err := c.ExpectString("Prompted value: "); err != nil {
			t.Errorf("Expected prompt error: %s", err)
		}
		if _, err := c.Send("\n"); err != nil {
			t.Errorf("Failed to send empty line to fake console: %s", err)
		}
		if _, err := c.Expect(expect.Regexp(regexp.MustCompile("A value is required"))); err != nil {
			t.Errorf("Expected missing value message error :%s", err)
		}
		if _, err := c.ExpectString("Prompted value: "); err != nil {
			t.Errorf("Expected prompt error: %s", err)
		}
		if _, err := c.Send("foo\n"); err != nil {
			t.Errorf("Failed to send value to fake console: %s", err)
		}
	}()

	var value string
	AskIfMissing(&value, "Prompted value")
	if value != "foo" {
		t.Errorf("Expected 'foo', got '%s' value", value)
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

	go func() {
		if _, err := c.ExpectString("Prompted password: "); err != nil {
			t.Errorf("Expected prompt error: %s", err)
		}
		if _, err := c.Send("\n"); err != nil {
			t.Errorf("Failed to send empty line to fake console: %s", err)
		}
		if _, err := c.Expect(expect.Regexp(regexp.MustCompile("A value is required"))); err != nil {
			t.Errorf("Expected missing value message error :%s", err)
		}
		if _, err := c.ExpectString("Prompted password: "); err != nil {
			t.Errorf("Expected prompt error: %s", err)
		}
		if _, err := c.Send("foo\n"); err != nil {
			t.Errorf("Failed to send value to fake console: %s", err)
		}
	}()

	var value string
	AskPasswordIfMissing(&value, "Prompted password")
	if value != "foo" {
		t.Errorf("Expected 'foo', got '%s' value", value)
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
