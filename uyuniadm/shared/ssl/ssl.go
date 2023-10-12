package ssl

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type CaChain struct {
	Root         string
	Intermediate []string
}

type SslPair struct {
	Cert string
	Key  string
}

// Generate the server certificate with the CA chain.
// Returns the certificate chain and the root CA.
func OrderCas(chain *CaChain, serverPair *SslPair) ([]byte, []byte) {
	CheckPaths(chain, serverPair)

	// Extract all certificates and their data
	certs := readCertificates(chain.Root)
	for _, caPath := range chain.Intermediate {
		certs = append(certs, readCertificates(caPath)...)
	}
	serverCerts := readCertificates(serverPair.Cert)
	certs = append(certs, serverCerts...)

	serverCert := findServerCert(certs)
	if serverCert == nil {
		log.Fatal().Msg("Failed to find a non-CA certificate")
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
	subjectKeyId string
	authKeyId    string
	isCa         bool
	isRoot       bool
}

func findServerCert(certs []certificate) *certificate {
	for _, cert := range certs {
		if !cert.isCa {
			return &cert
		}
	}
	return nil
}

func readCertificates(path string) []certificate {
	fd, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to read certificate file %s", path)
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
		cert := extractCertificateData(out)
		certs = append(certs, cert)
	}
	return certs
}

// Extract data from the certificate to help ordering and verifying it.
func extractCertificateData(content []byte) certificate {
	args := []string{"x509", "-noout", "-subject", "-subject_hash", "-startdate", "-enddate",
		"-issuer", "-issuer_hash", "-ext", "subjectKeyIdentifier,authorityKeyIdentifier,basicConstraints"}
	log.Debug().Msg("Running command openssl " + strings.Join(args, " "))
	cmd := exec.Command("openssl", args...)

	log.Trace().Msgf("Extracting data from certificate:\n%s", string(content))

	reader := bytes.NewReader(content)
	cmd.Stdin = reader

	out, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to extract data from certificate")
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
				log.Fatal().Err(err).Msgf("Failed to parse start date: %s\n", date)
			}
		} else if strings.HasPrefix(line, "notAfter=") {
			date := strings.SplitN(line, "=", 2)[1]
			cert.endDate, err = time.Parse(timeLayout, date)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed to parse end date: %s\n", date)
			}
		} else if strings.HasPrefix(line, "X509v3 Subject Key Identifier") {
			nextVal = "subjectKeyId"
		} else if strings.HasPrefix(line, "X509v3 Authority Key Identifier") {
			nextVal = "authKeyId"
		} else if strings.HasPrefix(line, "X509v3 Basic Constraints") {
			nextVal = "basicConstraints"
		} else if strings.HasPrefix(line, "    ") {
			if nextVal == "subjectKeyId" {
				cert.subjectKeyId = strings.ToUpper(strings.TrimSpace(line))
			} else if nextVal == "authKeyId" && strings.HasPrefix(line, "    keyid:") {
				cert.authKeyId = strings.ToUpper(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]))
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
		if cert.isCa && cert.authKeyId == "" {
			cert.authKeyId = cert.subjectKeyId
		}
	} else {
		cert.isRoot = false
	}
	return cert
}

// Prepare the certificate chain starting by the server up to the root CA.
// Returns the certificate chain and the root CA.
func sortCertificates(mapBySubjectHash map[string]certificate, serverCertHash string) ([]byte, []byte) {

	if len(mapBySubjectHash) == 0 {
		log.Fatal().Msg("No CA found in hash")
	}

	cert := mapBySubjectHash[serverCertHash]
	issuerHash := cert.issuerHash
	_, found := mapBySubjectHash[issuerHash]
	if issuerHash == "" || !found {
		log.Fatal().Msg("No CA found for server certificate")
	}

	sortedChain := bytes.NewBuffer(mapBySubjectHash[serverCertHash].content)
	var rootCa []byte

	for {
		cert, found = mapBySubjectHash[issuerHash]
		if !found {
			log.Fatal().Msgf("Missing CA with subject hash %s", issuerHash)
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
	return sortedChain.Bytes(), rootCa
}

// Ensures that all the passed path exists and the required files are available.
func CheckPaths(chain *CaChain, serverPair *SslPair) {
	mandatoryFile(chain.Root, "root CA")
	for _, ca := range chain.Intermediate {
		optionalFile(ca)
	}
	mandatoryFile(serverPair.Cert, "server certificate")
	mandatoryFile(serverPair.Key, "server key")
}

func mandatoryFile(file string, msg string) {
	if file == "" {
		log.Fatal().Msgf("%s is required", msg)
	}
	optionalFile(file)
}

func optionalFile(file string) {
	if file != "" && !utils.FileExists(file) {
		log.Fatal().Msgf("%s file is not accessible", file)
	}
}

// Converts an SSL key to RSA.
func GetRsaKey(keyPath string, password string) []byte {
	// Kubernetes only handles RSA private TLS keys, convert and strip password
	caPassword := password
	utils.AskIfMissing(&caPassword, "Source server SSL CA private key password")

	// Convert the key file to RSA format for kubectl to handle it
	cmd := exec.Command("openssl", "rsa", "-in", keyPath, "-passin", "env:pass")
	cmd.Env = append(cmd.Env, "pass="+caPassword)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to convert CA private key to RSA")
	}
	return out
}
