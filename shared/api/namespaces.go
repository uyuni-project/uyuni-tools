// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"errors"

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

func (c *APIClient) validateEndpoint(endpoint string) error {
	if c.AuthCookie == nil || endpoint == "" {
		return nil
	}
	return ValidateNamespace(c, endpoint)
}

func (c *APIClient) loadValidNamespaces() error {
	if c.validNamespaces != nil {
		return nil
	}

	res, err := c.Get("access/listNamespaces")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var response APIResponse[[]NamespaceAccess]
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Success {
		return errors.New(response.Message)
	}

	c.validNamespaces = map[string]struct{}{}
	for _, namespace := range response.Result {
		if namespace.Namespace != "" {
			c.validNamespaces[namespace.Namespace] = struct{}{}
		}
	}
	return nil
}
