// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestGetContainerImageDefaultNamespace(t *testing.T) {
	utils.DefaultNamespace = "mynamespace"
	proxyFlags := ProxyImageFlags{
		Tag: "mytag",
	}
	imageName := proxyFlags.GetContainerImage("httpd")
	imageExpected := "mynamespace/proxy-httpd:mytag"
	if imageName != imageExpected {
		t.Errorf("Image name %s does not match expected %s", imageName, imageExpected)
	}
}

func TestGetContainerImageCustomRegistry(t *testing.T) {
	utils.DefaultNamespace = "mynamespace"
	proxyFlags := ProxyImageFlags{
		Registry: "mytestregistry.example.com",
		Tag:      "mytag",
	}
	imageName := proxyFlags.GetContainerImage("httpd")
	imageExpected := "mytestregistry.example.com/proxy-httpd:mytag"
	if imageName != imageExpected {
		t.Errorf("Image name %s does not match expected %s", imageName, imageExpected)
	}
}
