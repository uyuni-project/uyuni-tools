// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func runGet(_ *types.GlobalFlags, flags *apiFlags, cmd *cobra.Command, args []string) error {
	log.Debug().Msgf("Running GET command %s", args[0])
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil && (client.Details.User != "" || client.Details.InSession) {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}
	path := args[0]
	options := args[1:]

	res, err := api.Get[interface{}](client, fmt.Sprintf("%s?%s", path, strings.Join(options, "&")))
	if err != nil {
		return utils.Errorf(err, L("error in query '%s'"), path)
	}

	// Check if the result is binary data by examining content type or trying JSON marshaling
	if err := outputResult(cmd, res); err != nil {
		return err
	}

	return nil
}

// outputResult handles outputting the API response appropriately based on content type.
// For JSON/text data, it pretty-prints the JSON.
// For binary data, it writes raw bytes to stdout.
func outputResult(cmd *cobra.Command, res *api.APIResponse[interface{}]) error {
	// Try to marshal as JSON first
	out, err := json.MarshalIndent(res.Result, "", "  ")
	if err != nil {
		// If JSON marshaling fails, treat as binary data
		log.Debug().Msg(L("Result is not JSON-serializable, treating as binary data"))

		// Convert result to bytes and write directly
		if res.Result != nil {
			if data, ok := res.Result.([]byte); ok {
				_, err := cmd.OutOrStdout().Write(data)
				return err
			}
			// Fallback: try to write as string
			_, err := fmt.Fprint(cmd.OutOrStdout(), res.Result)
			return err
		}
		return nil
	}

	// Check if output looks like binary (contains null bytes or non-printable chars)
	if containsBinaryData(out) {
		log.Debug().Msg(L("Detected binary data in response"))
		_, err := cmd.OutOrStdout().Write(out)
		return err
	}

	// Output JSON with newline
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(out))
	return err
}

// containsBinaryData checks if the data contains binary content.
func containsBinaryData(data []byte) bool {
	// Check for null bytes which indicate binary data
	for _, b := range data {
		if b == 0 {
			return true
		}
		// Check for high frequency of non-printable characters
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			return true
		}
	}
	return false
}

// isTextContentType checks if the content type indicates text-based data.
func isTextContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	// Text-based content types
	textTypes := []string{
		"application/json",
		"text/plain",
		"text/html",
		"text/xml",
		"application/xml",
		"application/javascript",
	}

	for _, t := range textTypes {
		if mediaType == t {
			return true
		}
	}

	return false
}
