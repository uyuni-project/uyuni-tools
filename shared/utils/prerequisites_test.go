// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"net"
	"testing"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestCheckPort(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Could not get a free port: %v", err)
	}

	// Should succeed
	err = CheckPort(port)
	if err != nil {
		t.Errorf("Expected port %d to be free, got error: %v", port, err)
	}

	// Bind to it intentionally
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("Could not bind to port: %v", err)
	}
	defer l.Close()

	// Should fail
	err = CheckPort(port)
	if err == nil {
		t.Errorf("Expected port %d to be in use and fail, but checking succeeded", port)
	}
}
