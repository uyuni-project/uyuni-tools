// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// CertificateInfo holds information about an SSL certificate for support purposes.
type CertificateInfo struct {
	Path        string
	Subject     string
	Issuer      string
	NotBefore   time.Time
	NotAfter    time.Time
	IsCA        bool
	IsExpired   bool
	ExpiresIn   time.Duration
	Fingerprint string
	Error       string
}

// CollectSSLCertInfo collects SSL certificate information from the container for debugging.
func CollectSSLCertInfo(dir string, execFunc func(command string, args ...string) ([]byte, error)) (string, error) {
	sslInfoFile, err := os.Create(path.Join(dir, "ssl-certificates-info"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), sslInfoFile.Name())
	}
	defer sslInfoFile.Close()

	var output bytes.Buffer
	output.WriteString("=== SSL Certificate Information ===\n")
	output.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	// Certificate paths to check
	certPaths := []struct {
		name string
		path string
	}{
		{"Root CA", CAContainerPath},
		{"DB Root CA", DBCAContainerPath},
		{"Server Certificate", ServerCertPath},
	}

	for _, cert := range certPaths {
		output.WriteString(fmt.Sprintf("--- %s (%s) ---\n", cert.name, cert.path))
		info := getCertificateInfo(cert.path, execFunc)
		output.WriteString(formatCertificateInfo(info))
		output.WriteString("\n")
	}

	// Check certificate chain validity
	output.WriteString("--- Certificate Chain Validation ---\n")
	chainValid := validateCertificateChain(execFunc)
	output.WriteString(chainValid)
	output.WriteString("\n")

	// Check certificate expiry warnings
	output.WriteString("--- Expiry Summary ---\n")
	expirySummary := getExpirySummary(certPaths, execFunc)
	output.WriteString(expirySummary)

	_, err = sslInfoFile.WriteString(output.String())
	if err != nil {
		return "", err
	}

	return sslInfoFile.Name(), nil
}

// getCertificateInfo retrieves certificate information using openssl.
func getCertificateInfo(
	certPath string,
	execFunc func(command string, args ...string) ([]byte, error),
) CertificateInfo {
	info := CertificateInfo{Path: certPath}

	// Check if file exists
	_, err := execFunc("test", "-f", certPath)
	if err != nil {
		info.Error = "Certificate file not found"
		return info
	}

	// Get certificate details using openssl
	out, err := execFunc("openssl", "x509", "-in", certPath, "-noout",
		"-subject", "-issuer", "-dates", "-fingerprint", "-ext", "basicConstraints")
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read certificate: %v", err)
		return info
	}

	parseCertificateOutput(string(out), &info)

	// Calculate expiry status
	now := time.Now()
	if !info.NotAfter.IsZero() {
		info.ExpiresIn = info.NotAfter.Sub(now)
		info.IsExpired = now.After(info.NotAfter)
	}

	return info
}

// parseCertificateOutput parses the openssl x509 output.
func parseCertificateOutput(output string, info *CertificateInfo) {
	lines := strings.Split(output, "\n")
	const timeLayout = "Jan  2 15:04:05 2006 GMT"
	const timeLayoutAlt = "Jan 2 15:04:05 2006 GMT"

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "subject="):
			info.Subject = strings.TrimPrefix(line, "subject=")
		case strings.HasPrefix(line, "issuer="):
			info.Issuer = strings.TrimPrefix(line, "issuer=")
		case strings.HasPrefix(line, "notBefore="):
			dateStr := strings.TrimPrefix(line, "notBefore=")
			if t, err := time.Parse(timeLayout, dateStr); err == nil {
				info.NotBefore = t
			} else if t, err := time.Parse(timeLayoutAlt, dateStr); err == nil {
				info.NotBefore = t
			}
		case strings.HasPrefix(line, "notAfter="):
			dateStr := strings.TrimPrefix(line, "notAfter=")
			if t, err := time.Parse(timeLayout, dateStr); err == nil {
				info.NotAfter = t
			} else if t, err := time.Parse(timeLayoutAlt, dateStr); err == nil {
				info.NotAfter = t
			}
		case strings.Contains(line, "Fingerprint="):
			// Handle SHA1, SHA256, and other fingerprint formats
			info.Fingerprint = strings.SplitN(line, "=", 2)[1]
		case strings.Contains(line, "CA:TRUE"):
			info.IsCA = true
		}
	}
}

// formatCertificateInfo formats certificate info for output.
func formatCertificateInfo(info CertificateInfo) string {
	var buf bytes.Buffer

	if info.Error != "" {
		buf.WriteString(fmt.Sprintf("ERROR: %s\n", info.Error))
		return buf.String()
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\n", info.Subject))
	buf.WriteString(fmt.Sprintf("Issuer: %s\n", info.Issuer))
	buf.WriteString(fmt.Sprintf("Valid From: %s\n", info.NotBefore.Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf("Valid Until: %s\n", info.NotAfter.Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf("Is CA: %t\n", info.IsCA))
	buf.WriteString(fmt.Sprintf("Fingerprint: %s\n", info.Fingerprint))

	if info.IsExpired {
		buf.WriteString("*** WARNING: CERTIFICATE IS EXPIRED ***\n")
	} else if info.ExpiresIn < 30*24*time.Hour {
		buf.WriteString(fmt.Sprintf("*** WARNING: Certificate expires in %d days ***\n",
			int(info.ExpiresIn.Hours()/24)))
	}

	return buf.String()
}

// validateCertificateChain checks if the certificate chain is valid.
func validateCertificateChain(execFunc func(command string, args ...string) ([]byte, error)) string {
	var buf bytes.Buffer

	// Verify the server certificate against the CA
	out, err := execFunc("openssl", "verify", "-CAfile", CAContainerPath, ServerCertPath)
	if err != nil {
		buf.WriteString(fmt.Sprintf("Chain validation FAILED: %v\n", err))
		buf.WriteString(fmt.Sprintf("Output: %s\n", string(out)))
	} else {
		outStr := string(out)
		if strings.Contains(outStr, ": OK") {
			buf.WriteString("Chain validation: OK\n")
		} else {
			buf.WriteString(fmt.Sprintf("Chain validation result: %s\n", outStr))
		}
	}

	return buf.String()
}

// getExpirySummary provides a summary of certificate expiry status.
func getExpirySummary(
	certPaths []struct{ name, path string },
	execFunc func(command string, args ...string) ([]byte, error),
) string {
	var buf bytes.Buffer
	hasIssues := false

	for _, cert := range certPaths {
		info := getCertificateInfo(cert.path, execFunc)
		if info.Error != "" {
			buf.WriteString(fmt.Sprintf("- %s: ERROR - %s\n", cert.name, info.Error))
			hasIssues = true
			continue
		}

		daysUntilExpiry := int(info.ExpiresIn.Hours() / 24)
		if info.IsExpired {
			buf.WriteString(fmt.Sprintf("- %s: *** EXPIRED *** (expired %d days ago)\n",
				cert.name, -daysUntilExpiry))
			hasIssues = true
		} else if daysUntilExpiry < 30 {
			buf.WriteString(fmt.Sprintf("- %s: *** WARNING *** expires in %d days\n",
				cert.name, daysUntilExpiry))
			hasIssues = true
		} else if daysUntilExpiry < 90 {
			buf.WriteString(fmt.Sprintf("- %s: expires in %d days\n", cert.name, daysUntilExpiry))
		} else {
			buf.WriteString(fmt.Sprintf("- %s: OK (expires in %d days)\n", cert.name, daysUntilExpiry))
		}
	}

	if !hasIssues {
		buf.WriteString("\nAll certificates are valid and not expiring soon.\n")
	} else {
		buf.WriteString("\n*** ATTENTION: Certificate issues detected. Please review above. ***\n")
	}

	return buf.String()
}

// CollectSSLCertInfoFromHost collects SSL certificate info by running openssl locally.
// This is used when we cannot exec into a container.
func CollectSSLCertInfoFromHost(dir string) (string, error) {
	sslInfoFile, err := os.Create(path.Join(dir, "ssl-certificates-info"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), sslInfoFile.Name())
	}
	defer sslInfoFile.Close()

	certPaths := []struct{ name, path string }{
		{"Root CA", "/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT"},
		{"Server Certificate", "/etc/pki/tls/certs/spacewalk.crt"},
	}

	var output bytes.Buffer
	output.WriteString("=== SSL Certificate Information ===\n")
	output.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	for _, cert := range certPaths {
		output.WriteString(fmt.Sprintf("--- %s (%s) ---\n", cert.name, cert.path))
		info := getCertificateInfoLocal(cert.path)
		output.WriteString(formatCertificateInfo(info))
		output.WriteString("\n")
	}

	_, err = sslInfoFile.WriteString(output.String())
	if err != nil {
		return "", err
	}

	return sslInfoFile.Name(), nil
}

// getCertificateInfoLocal retrieves certificate info from a local file.
func getCertificateInfoLocal(certPath string) CertificateInfo {
	info := CertificateInfo{Path: certPath}

	if !utils.FileExists(certPath) {
		info.Error = "Certificate file not found"
		return info
	}

	cmd := exec.Command("openssl", "x509", "-in", certPath, "-noout",
		"-subject", "-issuer", "-dates", "-fingerprint", "-ext", "basicConstraints")
	out, err := cmd.CombinedOutput()
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read certificate: %v", err)
		log.Debug().Err(err).Msgf("Failed to read certificate %s", certPath)
		return info
	}

	parseCertificateOutput(string(out), &info)

	// Calculate expiry status
	now := time.Now()
	if !info.NotAfter.IsZero() {
		info.ExpiresIn = info.NotAfter.Sub(now)
		info.IsExpired = now.After(info.NotAfter)
	}

	return info
}
