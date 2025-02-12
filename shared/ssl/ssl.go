// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// OrderCas generates the server certificate with the CA chain.
//
// Returns the certificate chain and the root CA.
func OrderCas(chain *types.CaChain, serverPair *types.SSLPair) ([]byte, []byte, error) {
	if err := CheckPaths(chain, serverPair); err != nil {
		return []byte{}, []byte{}, err
	}

	// Extract all certificates and their data
	certs, err := readCertificates(chain.Root)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	for _, caPath := range chain.Intermediate {
		intermediateCerts, err := readCertificates(caPath)
		if err != nil {
			return []byte{}, []byte{}, err
		}
		certs = append(certs, intermediateCerts...)
	}
	serverCerts, err := readCertificates(serverPair.Cert)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	certs = append(certs, serverCerts...)

	serverCert, err := findServerCert(certs)
	if err != nil {
		return []byte{}, []byte{}, errors.New(L("Failed to find a non-CA certificate"))
	}

	// Map all certificates using their hashes
	mapBySubjectHash := map[string]certificate{}
	if serverCert.subjectHash != "" {
		mapBySubjectHash[serverCert.subjectHash] = *serverCert
	}

	for _, caCert := range certs {
		if caCert.subjectHash != "" {
			mapBySubjectHash[caCert.subjectHash] = caCert
		}
	}

	// Sort from server certificate to RootCA
	return sortCertificates(mapBySubjectHash, serverCert.subjectHash)
}

type certificate struct {
	content      []byte
	subject      string
	subjectHash  string
	issuer       string
	issuerHash   string
	startDate    time.Time
	endDate      time.Time
	subjectKeyID string
	authKeyID    string
	isCa         bool
	isRoot       bool
}

func findServerCert(certs []certificate) (*certificate, error) {
	for _, cert := range certs {
		if !cert.isCa {
			return &cert, nil
		}
	}
	return nil, errors.New(L("expected to find a certificate, got none"))
}

func readCertificates(path string) ([]certificate, error) {
	fd, err := os.Open(path)
	if err != nil {
		return []certificate{}, utils.Errorf(err, L("Failed to read certificate file %s"), path)
	}

	certs := []certificate{}
	for {
		log.Debug().Msgf("Running openssl x509 on %s", path)
		cmd := exec.Command("openssl", "x509")
		cmd.Stdin = fd
		out, err := cmd.Output()

		if err != nil {
			// openssl got an invalid certificate or the end of the file
			break
		}

		// Extract data from the certificate
		cert, err := extractCertificateData(out)
		if err != nil {
			return []certificate{}, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// Extract data from the certificate to help ordering and verifying it.
func extractCertificateData(content []byte) (certificate, error) {
	args := []string{"x509", "-noout", "-subject", "-subject_hash", "-startdate", "-enddate",
		"-issuer", "-issuer_hash", "-ext", "subjectKeyIdentifier,authorityKeyIdentifier,basicConstraints"}
	log.Debug().Msg("Running command openssl " + strings.Join(args, " "))
	cmd := exec.Command("openssl", args...)

	log.Trace().Msgf("Extracting data from certificate:\n%s", string(content))

	reader := bytes.NewReader(content)
	cmd.Stdin = reader

	out, err := cmd.Output()
	if err != nil {
		return certificate{}, utils.Error(err, L("Failed to extract data from certificate"))
	}
	lines := strings.Split(string(out), "\n")

	cert := certificate{content: content}

	const timeLayout = "Jan 2 15:04:05 2006 MST"

	nextVal := ""
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "subject=") {
			cert.subject = strings.SplitN(line, "=", 2)[1]
		} else if strings.HasPrefix(line, "issuer=") {
			cert.issuer = strings.SplitN(line, "=", 2)[1]
		} else if strings.HasPrefix(line, "notBefore=") {
			date := strings.SplitN(line, "=", 2)[1]
			cert.startDate, err = time.Parse(timeLayout, date)
			if err != nil {
				return cert, utils.Errorf(err, L("Failed to parse start date: %s\n"), date)
			}
		} else if strings.HasPrefix(line, "notAfter=") {
			date := strings.SplitN(line, "=", 2)[1]
			cert.endDate, err = time.Parse(timeLayout, date)
			if err != nil {
				return cert, utils.Errorf(err, L("Failed to parse end date: %s\n"), date)
			}
		} else if strings.HasPrefix(line, "X509v3 Subject Key Identifier") {
			nextVal = "subjectKeyId"
		} else if strings.HasPrefix(line, "X509v3 Authority Key Identifier") {
			nextVal = "authKeyId"
		} else if strings.HasPrefix(line, "X509v3 Basic Constraints") {
			nextVal = "basicConstraints"
		} else if strings.HasPrefix(line, "    ") {
			if nextVal == "subjectKeyId" {
				cert.subjectKeyID = strings.ToUpper(strings.TrimSpace(line))
			} else if nextVal == "authKeyId" && strings.HasPrefix(line, "    keyid:") {
				cert.authKeyID = strings.ToUpper(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]))
			} else if nextVal == "basicConstraints" && strings.Contains(line, "CA:TRUE") {
				cert.isCa = true
			} else {
				// Unhandled extension value
				continue
			}
		} else if cert.subjectHash == "" {
			// subject_hash comes first without key to identify it
			cert.subjectHash = strings.TrimSpace(line)
		} else {
			// second issue_hash without key to identify this value
			cert.issuerHash = strings.TrimSpace(line)
		}
	}

	if cert.subject == cert.issuer {
		cert.isRoot = true
		// Some Root CAs might not have their authorityKeyIdentifier set to themself
		if cert.isCa && cert.authKeyID == "" {
			cert.authKeyID = cert.subjectKeyID
		}
	} else {
		cert.isRoot = false
	}
	return cert, nil
}

// Prepare the certificate chain starting by the server up to the root CA.
// Returns the certificate chain and the root CA.
func sortCertificates(mapBySubjectHash map[string]certificate, serverCertHash string) ([]byte, []byte, error) {
	if len(mapBySubjectHash) == 0 {
		return []byte{}, []byte{}, errors.New(L("No CA found"))
	}

	cert := mapBySubjectHash[serverCertHash]
	issuerHash := cert.issuerHash
	_, found := mapBySubjectHash[issuerHash]
	if issuerHash == "" || !found {
		return []byte{}, []byte{}, errors.New(L("No CA found for server certificate"))
	}

	sortedChain := bytes.NewBuffer(mapBySubjectHash[serverCertHash].content)
	var rootCa []byte

	for {
		cert, found = mapBySubjectHash[issuerHash]
		if !found {
			return []byte{}, []byte{}, fmt.Errorf(L("Missing CA with subject hash %s"), issuerHash)
		}

		nextHash := cert.issuerHash
		if nextHash == issuerHash {
			// Found Root CA, we can exit
			rootCa = cert.content
			break
		}
		issuerHash = nextHash
		sortedChain.Write(cert.content)
	}
	return sortedChain.Bytes(), rootCa, nil
}

// CheckPaths ensures that all the passed path exists and the required files are available.
func CheckPaths(chain *types.CaChain, serverPair *types.SSLPair) error {
	mandatoryFile(chain.Root, "root CA")
	for _, ca := range chain.Intermediate {
		if err := optionalFile(ca); err != nil {
			return err
		}
	}
	if err := mandatoryFile(serverPair.Cert, L("server certificate is required")); err != nil {
		return err
	}
	if err := mandatoryFile(serverPair.Key, L("server key is required")); err != nil {
		return err
	}
	return nil
}

func mandatoryFile(file string, msg string) error {
	if file == "" {
		return errors.New(msg)
	}
	return optionalFile(file)
}

func optionalFile(file string) error {
	if file != "" && !utils.FileExists(file) {
		return fmt.Errorf(L("%s file is not accessible"), file)
	}
	return nil
}

// Converts an SSL key to RSA.
func GetRsaKey(keyContent string, password string) []byte {
	// Kubernetes only handles RSA private TLS keys, convert and strip password
	caPassword := password
	utils.AskPasswordIfMissing(&caPassword, L("Source server SSL CA private key password"), 0, 0)

	// Convert the key file to RSA format for kubectl to handle it
	cmd := exec.Command("openssl", "rsa", "-passin", "env:pass")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to open openssl rsa process input stream"))
	}
	if _, err := io.WriteString(stdin, keyContent); err != nil {
		log.Fatal().Err(err).Msg(L("Failed to write openssl key content to input stream"))
	}

	cmd.Env = append(cmd.Env, "pass="+caPassword)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to convert CA private key to RSA"))
	}
	return out
}

// StripTextFromCertificate removes the optional text part of an x509 certificate.
func StripTextFromCertificate(certContent string) []byte {
	cmd := exec.Command("openssl", "x509")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to open openssl x509 process input stream"))
	}
	if _, err := io.WriteString(stdin, certContent); err != nil {
		log.Fatal().Err(err).Msg(L("Failed to write SSL certificate to input stream"))
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Msg(L("failed to strip text part from CA certificate"))
	}
	return out
}
