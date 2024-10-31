// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// SslCertGenerationFlags stores informations to generate an SSL Certificate.
type SslCertGenerationFlags struct {
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
}

// SslPair is a type for SSL Cert and Key.
type SslPair struct {
	Cert string
	Key  string
}
