// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestReadCertificatesRootCa(t *testing.T) {
	actual, err := readCertificates("testdata/chain1/root-ca.crt")
	testutils.AssertEquals(t, "error not nil", nil, err)
	testutils.AssertEquals(t, "Didn't get the expected certificates count", 1, len(actual))
	testutils.AssertTrue(t, "CA should be root", actual[0].isRoot)
}

func TestReadCertificatesNoCa(t *testing.T) {
	actual, err := readCertificates("testdata/chain1/server.crt")
	testutils.AssertEquals(t, "error not nil", nil, err)
	testutils.AssertEquals(t, "Didn't get the expected certificates count", 1, len(actual))
	testutils.AssertTrue(t, "Shouldn't be a CA certificate", !actual[0].isCa)
}

func TestReadCertificatesMultiple(t *testing.T) {
	actual, err := readCertificates("testdata/chain1/intermediate-ca.crt")
	testutils.AssertEquals(t, "error not nil", nil, err)
	testutils.AssertEquals(t, "Didn't get the expected certificates count", 2, len(actual))
	if len(actual) != 2 {
		t.Errorf("readCertificates got %d certificates; want 2", len(actual))
	}

	content := string(actual[0].content)
	if !strings.HasPrefix(content, "-----BEGIN CERTIFICATE-----\nMIIEXjCCA0agA") ||
		!strings.HasSuffix(content, "nrUN5m7Y0taw4qrOVOZRmGXu\n-----END CERTIFICATE-----\n") {
		t.Errorf("Wrong certificate content:\n%s", content)
	}

	testutils.AssertEquals(t, "Wrong certificate subject",
		"C=DE, ST=STATE, O=ORG, OU=ORGUNIT, CN=TeamCA",
		canonicalizeOpenSSLOutput(actual[1].subject),
	)

	testutils.AssertEquals(t, "Wrong subject hash", "85a51924", actual[1].subjectHash)

	testutils.AssertEquals(t, "Wrong certificate issuer",
		"C=DE, ST=STATE, L=CITY, O=ORG, OU=ORGUNIT, CN=RootCA",
		canonicalizeOpenSSLOutput(actual[0].issuer),
	)

	testutils.AssertEquals(t, "Wrong issuer hash", "e96ab651", actual[0].issuerHash)
	testutils.AssertTrue(t, "CA shouldn't be root", !actual[0].isRoot)
	testutils.AssertTrue(t, "Should be a CA", actual[0].isCa)

	testutils.AssertEquals(t, "Wrong subject key id",
		"62:00:25:E4:EE:70:E5:37:2D:1E:9E:AE:4E:B7:3E:FC:62:08:BF:27", actual[1].subjectKeyID,
	)

	testutils.AssertEquals(t, "Wrong auth key id",
		"6E:6D:4B:35:22:23:3E:13:18:A5:93:61:0E:9C:BE:1E:D2:B8:1B:D4", actual[0].authKeyID,
	)
}

// canonicalizeOpenSSLOutput standardizes openSSL test output across its versions.
func canonicalizeOpenSSLOutput(value string) string {
	return strings.ReplaceAll(value, " = ", "=")
}

func TestOrderCas(t *testing.T) {
	chain := types.CaChain{
		Root:         "testdata/chain1/root-ca.crt",
		Intermediate: []string{"testdata/chain1/intermediate-ca.crt"},
	}
	server := types.SSLPair{Cert: "testdata/chain1/server.crt", Key: "testdata/chain1/server.key"}

	certs, rootCa, err := OrderCas(&chain, &server)
	testutils.AssertEquals(t, "error not nil", nil, err)
	ordered := strings.Split(string(certs), "-----BEGIN CERTIFICATE-----\n")

	testutils.AssertEquals(t, "Found unknown content before first certificate", "", ordered[0])
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
	testutils.AssertEquals(t, "Wrong number of certificates in the chain", len(expected), len(onlyCerts))

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
	certsList, err := readCertificates("testdata/chain2/spacewalk.crt")
	testutils.AssertEquals(t, "error not nil", nil, err)

	actual, err := findServerCert(certsList)

	testutils.AssertEquals(t, "Expected to find a certificate, got none", nil, err)
	testutils.AssertEquals(t, "Wrong subject hash", "78b716a6", actual.subjectHash)
}

// Test a CA chain with all the chain in the server certificate file.
func TestOrderCasChain2(t *testing.T) {
	chain := types.CaChain{Root: "testdata/chain2/RHN-ORG-TRUSTED-SSL-CERT", Intermediate: []string{}}
	server := types.SSLPair{Cert: "testdata/chain2/spacewalk.crt", Key: "testdata/chain2/spacewalk.key"}

	certs, rootCa, err := OrderCas(&chain, &server)
	testutils.AssertEquals(t, "error not nil", nil, err)

	ordered := strings.Split(string(certs), "-----BEGIN CERTIFICATE-----\n")

	testutils.AssertEquals(t, "Found unknown content before first certificate", "", ordered[0])
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
	testutils.AssertEquals(t, "Wrong number of certificates in the chain", len(expected), len(onlyCerts))

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
	key := testutils.ReadFile(t, "testdata/RootCA.key")
	actual := string(GetRsaKey(key, "secret"))

	// This is what new openssl would generate
	matchingPKCS8 := strings.HasPrefix(actual, "-----BEGIN PRIVATE KEY-----\nMIIEugIBADANBgkqhkiG9w0BAQEFAAS") &&
		strings.HasSuffix(actual, "DKY9SmW6QD+RJwbMc4M=\n-----END PRIVATE KEY-----\n")

	// This is what older openssl would generate
	matchingPKCS1 := strings.HasPrefix(actual, "-----BEGIN RSA PRIVATE KEY-----\nMIIEoAIBAAKCAQEArqQvTR0") &&
		strings.HasSuffix(actual, "+3i4RXV4XtWHzmQymPUplukA/kScGzHOD\n-----END RSA PRIVATE KEY-----\n")

	if !matchingPKCS1 && !matchingPKCS8 {
		t.Errorf("Unexpected generated RSA key: %s", actual)
	}
}
