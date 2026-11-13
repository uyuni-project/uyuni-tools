// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

const unsupportedFunctionError = "This function is currently not supported by the server."

// ValidateNamespace checks if the server advertises the requested API namespace.
func ValidateNamespace(client *APIClient, namespace string) error {
	if namespace == "" {
		return nil
	}
	if err := client.loadValidNamespaces(); err != nil {
		return err
	}
	if _, ok := client.validNamespaces[namespace]; ok {
		return nil
	}
	return errors.New(L(unsupportedFunctionError))
}

func (c *APIClient) validateEndpoint(path string) error {
	if c.AuthCookie == nil || path == "" {
		return nil
	}
	endpoint, _, _ := strings.Cut(path, "?")
	return ValidateNamespace(c, endpoint)
}

func (c *APIClient) loadValidNamespaces() error {
	if c.validNamespaces != nil {
		return nil
	}

	res, err := c.Get("api/getApiCallList")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var response APIResponse[map[string]map[string]Endpoint]
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Success {
		return errors.New(response.Message)
	}

	c.validNamespaces = map[string]struct{}{}
	for namespace, endpoints := range response.Result {
		for _, endpoint := range endpoints {
			c.validNamespaces[fmt.Sprintf("%s/%s", strings.ReplaceAll(namespace, ".", "/"), endpoint.Name)] = struct{}{}
		}
	}
	return nil
}
