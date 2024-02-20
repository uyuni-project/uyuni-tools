// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"strings"
	"testing"
)

func TestReadCertificatesRootCa(t *testing.T) {
	actual := readCertificates("testdata/chain1/root-ca.crt")
	if len(actual) != 1 {
		t.Errorf("readCertificates got %d certificates; want 1", len(actual))
	}

	if !actual[0].isRoot {
		t.Error("CA should be root")
	}
}

func TestReadCertificatesNoCa(t *testing.T) {
	actual := readCertificates("testdata/chain1/server.crt")
	if len(actual) != 1 {
		t.Errorf("readCertificates got %d certificates; want 1", len(actual))
	}

	if actual[0].isCa {
		t.Error("Shouldn't be a CA certificate")
	}
}

func TestReadCertificatesMultiple(t *testing.T) {
	actual := readCertificates("testdata/chain1/intermediate-ca.crt")
	if len(actual) != 2 {
		t.Errorf("readCertificates got %d certificates; want 2", len(actual))
	}

	content := string(actual[0].content)
	if !strings.HasPrefix(content, "-----BEGIN CERTIFICATE-----\nMIIEXjCCA0agA") ||
		!strings.HasSuffix(content, "nrUN5m7Y0taw4qrOVOZRmGXu\n-----END CERTIFICATE-----\n") {
		t.Errorf("Wrong certificate content:\n%s", content)
	}

	if actual[1].subject != "C = DE, ST = STATE, O = ORG, OU = ORGUNIT, CN = TeamCA" {
		t.Errorf("Wrong certificate subject: %s", actual[1].subject)
	}

	if actual[1].subjectHash != "85a51924" {
		t.Errorf("Wrong subject hash: %s", actual[1].subjectHash)
	}

	if actual[0].issuer != "C = DE, ST = STATE, L = CITY, O = ORG, OU = ORGUNIT, CN = RootCA" {
		t.Errorf("Wrong certificate issuer: %s", actual[0].issuer)
	}

	if actual[0].issuerHash != "e96ab651" {
		t.Errorf("Wrong issuer hash: %s", actual[0].issuerHash)
	}

	if actual[0].isRoot {
		t.Error("CA shouldn't be root")
	}

	if !actual[0].isCa {
		t.Error("Should be a CA")
	}

	if actual[1].subjectKeyId != "62:00:25:E4:EE:70:E5:37:2D:1E:9E:AE:4E:B7:3E:FC:62:08:BF:27" {
		t.Errorf("Wrong subject key id: %s", actual[1].subjectKeyId)
	}

	if actual[0].authKeyId != "6E:6D:4B:35:22:23:3E:13:18:A5:93:61:0E:9C:BE:1E:D2:B8:1B:D4" {
		t.Errorf("Wrong auth key id: %s", actual[0].authKeyId)
	}
}

func TestOrderCas(t *testing.T) {
	chain := CaChain{Root: "testdata/chain1/root-ca.crt", Intermediate: []string{"testdata/chain1/intermediate-ca.crt"}}
	server := SslPair{Cert: "testdata/chain1/server.crt", Key: "testdata/chain1/server.key"}

	certs, rootCa := OrderCas(&chain, &server)
	ordered := strings.Split(string(certs), "-----BEGIN CERTIFICATE-----\n")

	if ordered[0] != "" {
		t.Errorf("Found unknown content before first certificate: %s", ordered[0])
	}
	onlyCerts := ordered[1:]

	expected := []struct {
		Begin string
		End   string
	}{
		{Begin: "MIIEdDCCA1ygAwIBAgIUZ2P1Ka9Eun", End: "JtS8rmkQpYyJciifX0PxYzTg=="},
		{Begin: "MIIETzCCAzegAwIBAgIUZ2P1Ka9Eun", End: "s3DjcCbkzyTUCKh9Po4\nmoUf"},
		{Begin: "MIIEXjCCA0agAwIBAgIUZ2P1Ka9Eunnv3dy/", End: "nrUN5m7Y0taw4qrOVOZRmGXu"},
	}

	// Do not count the empty first item
	if len(onlyCerts) != len(expected) {
		t.Errorf("Wrong number of certificates in the chain: got %d; want %d", len(onlyCerts), len(expected))
	}

	for i, data := range expected {
		if !strings.HasPrefix(onlyCerts[i], data.Begin) ||
			!strings.HasSuffix(onlyCerts[i], data.End+"\n-----END CERTIFICATE-----\n") {
			t.Errorf("Invalid certificate #%d, got:\n:%s", i, onlyCerts[i])
		}
	}

	rootCert := string(rootCa)
	if !strings.HasPrefix(rootCert, "-----BEGIN CERTIFICATE-----\nMIIEVjCCAz6gAwIBAgIUSZYESIXLDe") ||
		!strings.HasSuffix(rootCert, "5c7cfxV\nkABuj9PJxnNnFQ==\n-----END CERTIFICATE-----\n") {
		t.Errorf("Invalid root CA certificate, got:\n:%s", rootCert)
	}
}

func TestFindServerCertificate(t *testing.T) {
	certsList := readCertificates("testdata/chain2/spacewalk.crt")
	actual, err := findServerCert(certsList)

	if err != nil {
		t.Error("Expected to find a certificate, got none")
	}

	if actual.subjectHash != "78b716a6" {
		t.Errorf("Wrong subject hash, got %s", actual.subjectHash)
	}
}

// Test a CA chain with all the chain in the server certificate file.
func TestOrderCasChain2(t *testing.T) {
	chain := CaChain{Root: "testdata/chain2/RHN-ORG-TRUSTED-SSL-CERT", Intermediate: []string{}}
	server := SslPair{Cert: "testdata/chain2/spacewalk.crt", Key: "testdata/chain2/spacewalk.key"}

	certs, rootCa := OrderCas(&chain, &server)
	ordered := strings.Split(string(certs), "-----BEGIN CERTIFICATE-----\n")

	if ordered[0] != "" {
		t.Errorf("Found unknown content before first certificate: %s", ordered[0])
	}
	onlyCerts := ordered[1:]

	expected := []struct {
		Begin string
		End   string
	}{
		{Begin: "MIIEejCCA2KgAwIBAgIUEbWzxg57E", End: "Ur+fgZpBNvbkjD8b+S0ECQA6Dg=="},
		{Begin: "MIIETzCCAzegAwIBAgIUEbWzxg57E", End: "TT2Sljt0YfkmWfdXA\nwOUt"},
		{Begin: "MIIEXjCCA0agAwIBAgIUEbWzxg57E", End: "ivyvRvlwCUNstG6u8Y7IxHHn"},
	}

	// Do not count the empty first item
	if len(onlyCerts) != len(expected) {
		t.Errorf("Wrong number of certificates in the chain: got %d; want %d", len(onlyCerts), len(expected))
	}

	for i, data := range expected {
		if !strings.HasPrefix(onlyCerts[i], data.Begin) ||
			!strings.HasSuffix(onlyCerts[i], data.End+"\n-----END CERTIFICATE-----\n") {
			t.Errorf("Invalid certificate #%d, got:\n:%s", i, onlyCerts[i])
		}
	}

	rootCert := string(rootCa)
	if !strings.HasPrefix(rootCert, "-----BEGIN CERTIFICATE-----\nMIIEVjCCAz6gAwIBAgIUA12e94NK") ||
		!strings.HasSuffix(rootCert, "AQKotV5y5qBInw==\n-----END CERTIFICATE-----\n") {
		t.Errorf("Invalid root CA certificate, got:\n:%s", rootCert)
	}
}

func TestGetRsaKey(t *testing.T) {
	actual := string(GetRsaKey("testdata/RootCA.key", "secret"))
	if !strings.HasPrefix(actual, "-----BEGIN PRIVATE KEY-----\nMIIEugIBADANBgkqhkiG9w0BAQEFAAS") ||
		!strings.HasSuffix(actual, "DKY9SmW6QD+RJwbMc4M=\n-----END PRIVATE KEY-----\n") {
		t.Errorf("Unexpected generated RSA key: %s", actual)
	}
}
