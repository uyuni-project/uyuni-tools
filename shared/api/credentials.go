// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Store API credentials for future API use.
func StoreLoginCreds(client *APIClient) error {
	if client.AuthCookie.Value == "" {
		return fmt.Errorf(L("not logged in, session cookie is missing"))
	}
	// Future: Add support for more servers if needed in the future
	auth := []authStorage{
		{
			Session: client.AuthCookie.Value,
			Server:  client.Details.Server,
			CApath:  client.Details.CApath,
		},
	}

	authData, err := json.Marshal(auth)
	if err != nil {
		return utils.Errorf(err, L("unable to create credentials json"))
	}

	err = os.WriteFile(getAPICredsFile(), authData, 0600)
	if err != nil {
		return utils.Errorf(err, L("unable to write credentials store %s"), getAPICredsFile())
	}
	return nil
}

// Remove stored API credentials.
func RemoveLoginCreds() error {
	//Future: Multi-server support will need some parsing here
	return os.Remove(getAPICredsFile())
}

// Fills ConnectionDetails with cached credentials and returns true.
// In case cached credentials are not present, asks for password and returns false.
func getLoginCredentials(conn *ConnectionDetails) error {
	// Load stored credentials if no user was specified
	if IsAlreadyLoggedIn() && conn.User == "" {
		if err := loadLoginCreds(conn); err != nil {
			log.Warn().Err(err).Msg(L("Cannot load stored credentials"))
			if err := RemoveLoginCreds(); err != nil {
				return utils.Errorf(err, L("Failed to remove stored credentials!"))
			}
		} else {
			// We have connection cookie
			conn.InSession = true
			return nil
		}
	}

	// If user name provided, but no password and not loaded
	if conn.User != "" {
		utils.AskIfMissing(&conn.User, L("API server user"), 0, 0, nil)
	}
	if conn.Password == "" {
		utils.AskPasswordIfMissing(&conn.Password, L("API server password"), 0, 0)
	}

	if conn.User == "" || conn.Password == "" {
		return fmt.Errorf(L("No credentials provided"))
	}

	return nil
}

// Read stored session and server details.
func loadLoginCreds(connection *ConnectionDetails) error {
	data, err := os.ReadFile(getAPICredsFile())
	if err != nil {
		return utils.Errorf(err, L("unable to read credentials file %s"), getAPICredsFile())
	}
	authStore := []authStorage{}
	err = json.Unmarshal(data, &authStore)
	if err != nil {
		return utils.Errorf(err, L("unable to decode credentials file"))
	}

	if len(authStore) == 0 {
		return fmt.Errorf(L("no credentials loaded"))
	}

	// Currently we support storing data only to one server
	// Future: add support for more servers if wanted

	authData := authStore[0]

	if connection.Server != "" && connection.Server != authData.Server {
		return fmt.Errorf(L("specified api server does not match with stored credentials"))
	}
	connection.Server = authData.Server
	if authData.CApath != "" {
		connection.CApath = authData.CApath
	}

	connection.Cookie = authData.Session

	return nil
}

// Returns true if credentials file already exists.
// Does not check for credentials validity.
func IsAlreadyLoggedIn() bool {
	return utils.FileExists(getAPICredsFile())
}

func getAPICredsFile() string {
	return path.Join(utils.GetUserConfigDir(), api_credentials_store)
}
