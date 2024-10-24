// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// StoreLoginCreds stores the API credentials for future API use.
func StoreLoginCreds(client *APIClient) error {
	if client.AuthCookie.Value == "" {
		return errors.New(L("not logged in, session cookie is missing"))
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

// RemoveLoginCreds removes the stored API credentials.
func RemoveLoginCreds() error {
	//Future: Multi-server support will need some parsing here
	return os.Remove(getAPICredsFile())
}

// Asks for not provided ConnectionDetails or errors out.
func getLoginCredentials(conn *ConnectionDetails) error {
	// If user name provided, but no password and not loaded
	utils.AskIfMissing(&conn.Server, L("API server URL"), 0, 0, nil)
	utils.AskIfMissing(&conn.User, L("API server user"), 0, 0, nil)
	utils.AskPasswordIfMissingOnce(&conn.Password, L("API server password"), 0, 0)

	if conn.User == "" || conn.Password == "" {
		return errors.New(L("No credentials provided"))
	}

	return nil
}

// Fills ConnectionDetails with cached credentials if possible.
func getStoredConnectionDetails(conn *ConnectionDetails) {
	if IsAlreadyLoggedIn() && conn.User == "" {
		if err := loadLoginCreds(conn); err != nil {
			log.Warn().Err(err).Msg(L("Cannot load stored credentials"))
			if err := RemoveLoginCreds(); err != nil {
				log.Warn().Err(err).Msg(L("Failed to remove stored credentials!"))
			}
		} else {
			// We have connection cookie
			conn.InSession = true
		}
	}
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
		return errors.New(L("no credentials loaded"))
	}

	// Currently we support storing data only to one server
	// Future: add support for more servers if wanted

	authData := authStore[0]

	if connection.Server != "" && connection.Server != authData.Server {
		return errors.New(L("specified api server does not match with stored credentials"))
	}
	connection.Server = authData.Server
	if authData.CApath != "" {
		connection.CApath = authData.CApath
	}

	connection.Cookie = authData.Session

	return nil
}

// IsAlreadyLoggedIn returns true if credentials file already exists.
//
// Does not check for credentials validity.
func IsAlreadyLoggedIn() bool {
	return utils.FileExists(getAPICredsFile())
}

func getAPICredsFile() string {
	return path.Join(utils.GetUserConfigDir(), api_credentials_store)
}
