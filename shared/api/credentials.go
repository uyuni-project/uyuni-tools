// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
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
func StoreLoginCreds(ctx context.Context, connection *ConnectionDetails) error {
	// check login is valid
	if !checkCredentials(connection) {
		return fmt.Errorf(L("Failed to validate credentials"))
	}

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
		log.Error().Msgf(L("Unable to create credentials json"))
		return err
	}

	err = os.WriteFile(getAPICredsFile(), authData, 0)
	if err != nil {
		log.Error().Msgf(L("Unable to write credentials store"))
		return err
	}

	log.Info().Msg(L("Login credentials verified and stored"))
	return nil
}

// Remove stored API credentials.
func RemoveLoginCreds(ctx context.Context) error {
	if err := os.Remove(getAPICredsFile()); err != nil {
		return err
	}
	log.Info().Msg(L("Successfully logged out"))
	return nil
}

// Read and decrypt stored login credentials.
func LoadLoginCreds(ctx context.Context, connection *ConnectionDetails) error {
	data, err := os.ReadFile(getAPICredsFile())
	if err != nil {
		return utils.Errorf(err, L("Unable to read credentials file"))
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

func checkCredentials(connection *ConnectionDetails) bool {
	_, err := Init(connection)
	return err == nil
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
