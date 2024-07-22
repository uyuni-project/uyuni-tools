// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"golang.org/x/crypto/nacl/secretbox"
)

// Store API credentials for future API use.
func StoreLoginCreds(connection *ConnectionDetails) error {
	// encrypt login data
	urldata := getFixedServerString(connection.Server)
	encodedURL := make([]byte, base64.URLEncoding.EncodedLen(len(urldata)))
	base64.URLEncoding.Encode(encodedURL, urldata)

	encUser, err := encryptMsg(connection.User, [32]byte(encodedURL))
	if err != nil {
		return err
	}
	encPassword, err := encryptMsg(connection.Password, [32]byte(encodedURL))
	if err != nil {
		return err
	}

	// store encrypted credentials into separate credentials config storage
	auth := authStorage{
		User:     encUser,
		Password: encPassword,
		Server:   connection.Server,
	}

	authData, err := json.Marshal(auth)
	if err != nil {
		return utils.Errorf(err, L("Unable to create credentials json"))
	}

	err = os.WriteFile(getAPICredsFile(), authData, 0600)
	if err != nil {
		return utils.Errorf(err, L("Unable to write credentials store %s"), getAPICredsFile())
	}
	return nil
}

// Remove stored API credentials.
func RemoveLoginCreds() error {
	return os.Remove(getAPICredsFile())
}

// Check if login credentials are valid.
func (c *APIClient) ValidateCreds() bool {
	err := c.Login()
	return err == nil
}

// Fills ConnectionDetails with cached credentials and returns true.
// In case cached credentials are not present, asks for password and returns false.
func getLoginCredentials(conn *ConnectionDetails) error {
	// Load stored credentials if no user was specified
	cachedCredentials := false
	if utils.FileExists(getAPICredsFile()) && conn.User == "" {
		if err := loadLoginCreds(conn); err != nil {
			log.Warn().Err(err).Msg(L("Cannot load stored credentials"))
			if err := RemoveLoginCreds(); err != nil {
				log.Warn().Err(err).Msg(L("Failed to remove stored credentials!"))
			}
			return err
		} else {
			cachedCredentials = true
		}
	}

	// If user name provided, but no password and not loaded
	if conn.User != "" {
		if conn.Password == "" {
			utils.AskPasswordIfMissing(&conn.Password, L("API server password"), 0, 0)
		}
	}

	conn.Cached = cachedCredentials
	return nil
}

// Read and decrypt stored login credentials.
func loadLoginCreds(connection *ConnectionDetails) error {
	data, err := os.ReadFile(getAPICredsFile())
	if err != nil {
		return utils.Errorf(err, L("Unable to read credentials file %s"), getAPICredsFile())
	}
	authData := authStorage{}
	err = json.Unmarshal(data, &authData)
	if err != nil {
		return utils.Errorf(err, L("Unable to decode credentials file"))
	}
	connection.Server = authData.Server

	urldata := getFixedServerString(connection.Server)
	encodedURL := make([]byte, base64.URLEncoding.EncodedLen(len(urldata)))
	base64.URLEncoding.Encode(encodedURL, urldata)

	decUser, err := decryptMsg(authData.User, [32]byte(encodedURL))
	if err != nil {
		return err
	}
	decPassword, err := decryptMsg(authData.Password, [32]byte(encodedURL))
	if err != nil {
		return err
	}
	connection.Password = string(decPassword)
	connection.User = string(decUser)
	return nil
}

func getAPICredsFile() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			xdgConfigHome = path.Join(home, ".config")
		}
	}
	return path.Join(xdgConfigHome, api_credentials_store)
}

func getFixedServerString(server string) []byte {
	res := make([]byte, 32)
	for i := range res {
		res[i] = '_'
	}
	if len(server) > 32 {
		copy(res, server[:32])
	} else {
		copy(res, server)
	}
	return res
}

func encryptMsg(message string, secret [32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	return secretbox.Seal(nonce[:], []byte(message), &nonce, &secret), nil
}

func decryptMsg(message []byte, secret [32]byte) (string, error) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], message[:24])
	decrypted, err := secretbox.Open(nil, message[24:], &decryptNonce, &secret)
	if !err {
		return "", fmt.Errorf(L("Decoding of secret failed"))
	}
	return string(decrypted), nil
}
