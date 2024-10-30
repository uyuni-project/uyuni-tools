// SPDX-FileCopyrightText: 2025 SUSE LLC
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

// IsThirdParty returns whether the CA chain is a third party one.
func (c *CaChain) IsThirdParty() bool {
	return c.IsDefined() && c.Key == ""
}

// IsDefined returns whether the CA chain is defined.
// At least the CA root certificate is available.
func (c *CaChain) IsDefined() bool {
	return c.Root != ""
}

// SSLPair is a type for SSL Cert and Key.
type SSLPair struct {
	Cert string
	Key  string
}

// IsDefined returns whether the SSL pair is defined.
func (p *SSLPair) IsDefined() bool {
	return p.Cert != "" && p.Key != ""
}
