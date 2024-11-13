// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

// Models/Schemas for the proxy API.

// ProxyConfigRequest is the request schema for the proxy/containerConfig endpoint when
// user has proxy certificates.
type ProxyConfigRequest struct {
	ProxyName       string
	ProxyPort       int
	Server          string
	MaxCache        int
	Email           string
	RootCA          string
	ProxyCrt        string
	ProxyKey        string
	IntermediateCAs []string
}

// ProxyConfigGenerateRequest is the request schema for the proxy/containerConfig endpoint when
// user wants to generate proxy certificates.
type ProxyConfigGenerateRequest struct {
	ProxyName  string
	ProxyPort  int
	Server     string
	MaxCache   int
	Email      string
	CaCrt      string
	CaKey      string
	CaPassword string
	Cnames     []string
	Country    string
	State      string
	City       string
	Org        string
	OrgUnit    string
	SSLEmail   string
}
