// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetFilename returns the filename to save a configuration file to.
// If an output filename is not specified, then the filename is based on the proxy name.
func GetFilename(output string, proxyName string) string {
	filename := output
	if filename == "" {
		filename = strings.Split(proxyName, ".")[0] + "-config"
	}
	return filename + ".tar.gz"
}

// Prompt for password.
func PromptForPassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter %s: ", caPassword)
	password, _ := reader.ReadString('\n')
	return password
}
