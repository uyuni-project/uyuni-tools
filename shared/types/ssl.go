// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// SSLCertGenerationFlags stores informations to generate an SSL Certificate.
type SSLCertGenerationFlags struct {
	Cnames   []string `mapstructure:"cname"`
	Country  string
	State    string
	City     string
	Org      string
	OU       string
	Password string
	Email    string
}

// CaChain is a type to store CA Chain.
type CaChain struct {
	Root         string
	Intermediate []string
	// Key is the CA key file in the case of a migration of a self-generate CA.
	Key string
}

// SSLPair is a type for SSL Cert and Key.
type SSLPair struct {
	Cert string
	Key  string
}
