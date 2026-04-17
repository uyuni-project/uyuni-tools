// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/templates"
	pxyutils "github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"gopkg.in/yaml.v2"
)

func TestCheckDirPermissions(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := checkPermissions(tempDir, 0005|0050|0500); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateYamlFiles(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{"httpd.yaml", "ssh.yaml", "config.yaml"}
	for _, file := range testFiles {
		filePath := path.Join(tempDir, file)
		if _, err := os.Create(filePath); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Test: when all files are present and have correct permissions
	if err := validateInstallYamlFiles(tempDir); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test: Missing file scenario, remove one file and expect an error
	os.Remove(path.Join(tempDir, "httpd.yaml"))
	if err := validateInstallYamlFiles(tempDir); err == nil {
		t.Errorf("Expected an error due to missing httpd.yaml, but got none")
	}
}

func TestGetSystemID(t *testing.T) {
	// event output
	systemid := `<?xml version=\"1.0\"?><params><param><value><struct><member><name>username</name>` +
		`<value><string>admin</string></value></member><member><name>os_release</name><value><string>6.1</string>` +
		`</value></member><member><name>operating_system</name><value><string>SL-Micro</string></value></member>` +
		`<member><name>architecture</name><value><string>x86_64-redhat-linux</string></value></member><member>` +
		`<name>system_id</name><value><string>ID-1000010001</string></value></member><member><name>type</name><value>` +
		`<string>REAL</string></value></member><member><name>fields</name><value><array><data><value>` +
		`<string>system_id</string></value><value><string>os_release</string></value><value><string>operating_system` +
		`</string></value><value><string>architecture</string></value><value><string>username</string></value><value>` +
		`<string>type</string></value></data></array></value></member><member><name>checksum</name><value>` +
		`<string>1aaa4427328cfd7fbd613693802e0920d9f1c1ea2b3d31a869ed1ac3fbfe4174</string></value></member></struct>` +
		`</value></param></params>`

	event := `suse/systemid/generated {"data": "` + systemid + `", "_stamp": "2025-08-04T12:04:29.403745"}`

	// create custom runners
	contextRunner = testutils.FakeContextRunnerGenerator(event, nil)
	newRunner = testutils.FakeRunnerGenerator("", nil)

	received, err := getSystemIDEvent()
	testutils.AssertNoError(t, "error during obtaining systemid", err)
	testutils.AssertEquals(t, "received event differs", []byte(event), received)

	receivedSystemid, err := parseSystemIDEvent(received)
	testutils.AssertNoError(t, "error during event decoding", err)
	// unquote raw string before comparing.
	unquotedSystemid, _ := strconv.Unquote(`"` + systemid + `"`)
	testutils.AssertEquals(t, "received systemid differs", unquotedSystemid, receivedSystemid)
}

func TestGenerateSystemdService(t *testing.T) {
	driver := testutils.FakeSystemdDriver{
		ReloadDaemonError: nil,
	}
	systemd := podman.NewSystemdWithDriver(&driver)

	const httpProxyConfig = "/fake/proxy"
	const ipv6 = true
	var callsCount = 0

	systemdGenerator = func(template utils.Template, service string, image string, config string) error {
		callsCount++

		// Assert that the calls are correct
		switch service {
		case "httpd":
			testutils.AssertEquals(t, "Wrong httpd image", "fake/httpd:latest", image)
			testutils.AssertStringContains(t, "Missing httpd tuning configuration", config,
				"Environment=HTTPD_EXTRA_CONF=-v/my/apache.conf:/etc/apache2/conf.d/apache_tuning.conf:ro",
			)
			httpTemplate, ok := template.(templates.HttpdTemplateData)
			testutils.AssertTrue(t, "httpd template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong http service data",
				templates.HttpdTemplateData{
					Volumes:       utils.ProxyHttpdVolumes,
					HTTPProxyFile: httpProxyConfig,
				},
				httpTemplate,
			)

		case "salt-broker":
			testutils.AssertEquals(t, "Wrong salt-broker image", "fake/salt-broker:latest", image)
			saltBrokerTemplate, ok := template.(templates.SaltBrokerTemplateData)
			testutils.AssertTrue(t, "salt broker template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong salt broker service data",
				templates.SaltBrokerTemplateData{
					HTTPProxyFile: httpProxyConfig,
				},
				saltBrokerTemplate,
			)

		case "squid":
			testutils.AssertEquals(t, "Wrong squid image", "fake/squid:latest", image)
			testutils.AssertStringContains(t, "Missing squid tuning configuration", config,
				"Environment=SQUID_EXTRA_CONF=-v/my/squid.conf:/etc/squid/conf.d/squid_tuning.conf:ro",
			)
			squidTemplate, ok := template.(templates.SquidTemplateData)
			testutils.AssertTrue(t, "squid template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong squid service data",
				templates.SquidTemplateData{
					Volumes:       utils.ProxySquidVolumes,
					HTTPProxyFile: httpProxyConfig,
				},
				squidTemplate,
			)

		case "ssh":
			testutils.AssertEquals(t, "Wrong ssh image", "fake/ssh:latest", image)
			testutils.AssertStringContains(t, "Missing squid tuning configuration", config,
				"Environment=SSH_EXTRA_CONF=-v/my/sshd.conf:/etc/ssh/sshd_config.d/10-tuning.conf:ro",
			)
			sshTemplate, ok := template.(templates.SSHTemplateData)
			testutils.AssertTrue(t, "ssh template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong ssh service data",
				templates.SSHTemplateData{
					HTTPProxyFile: httpProxyConfig,
				},
				sshTemplate,
			)

		case "tftpd":
			testutils.AssertEquals(t, "Wrong tftp image", "fake/tftp:latest", image)
			tftpdTemplate, ok := template.(templates.TFTPDTemplateData)
			testutils.AssertTrue(t, "tftpd template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong tftpd service data",
				templates.TFTPDTemplateData{
					Volumes:       utils.ProxyTftpdVolumes,
					HTTPProxyFile: httpProxyConfig,
				},
				tftpdTemplate,
			)

		case "pod":
			testutils.AssertEquals(t, "The pod should not have an image", "", image)
			testutils.AssertStringContains(t, "Missing podman extra args", config,
				"Environment=\"PODMAN_EXTRA_ARGS=--arg1 --arg2\"",
			)
			podTemplate, ok := template.(templates.PodTemplateData)
			testutils.AssertTrue(t, "pod template of unexpected type", ok)
			testutils.AssertEquals(t, "Wrong pod data", templates.PodTemplateData{
				Ports: []types.PortMap{
					{
						Port:    22,
						Exposed: 8022,
					},
					utils.NewPortMap(4505),
					utils.NewPortMap(4506),
					utils.NewPortMap(443),
					utils.NewPortMap(80),
					{
						Exposed:  69,
						Port:     69,
						Protocol: "udp",
					},
				},
				HTTPProxyFile: httpProxyConfig,
				Network:       podman.UyuniNetwork,
				IPV6Enabled:   ipv6,
			}, podTemplate)

		default:
			t.Errorf("Unexpected systemd service generation: %s", service)
		}
		return nil
	}

	err := GenerateSystemdService(systemd, "fake/httpd:latest", "fake/salt-broker:latest",
		"fake/squid:latest", "fake/ssh:latest",
		"fake/tftp:latest", &PodmanProxyFlags{
			Podman: podman.PodmanFlags{Args: []string{"--arg1", "--arg2"}},
			ProxyImageFlags: pxyutils.ProxyImageFlags{
				Tuning: pxyutils.Tuning{
					Httpd: "/my/apache.conf",
					Squid: "/my/squid.conf",
					SSH:   "/my/sshd.conf",
				},
			},
		},
		ipv6, httpProxyConfig)
	testutils.AssertTrue(t, fmt.Sprintf("Unexpected error: %v", err), err == nil)
	testutils.AssertEquals(t, "Unexpected number of services generated", 6, callsCount)

	// Restore the mocked variables
	systemdGenerator = generateSystemdFile
}

func TestExtractSecrets(t *testing.T) {
	tempDir := t.TempDir()
	oldProxyConfigDir := proxyConfigDir
	proxyConfigDir = tempDir
	defer func() { proxyConfigDir = oldProxyConfigDir }()

	secrets := make(map[string]string)
	oldCreateSecret := createSecret
	createSecret = func(name string, value string) error {
		secrets[name] = value
		return nil
	}
	defer func() { createSecret = oldCreateSecret }()

	configPath := path.Join(tempDir, "config.yaml")
	configData := "ca_crt: CA_CERT_CONTENT\nother_key: other_value\n"
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatal(err)
	}

	httpdPath := path.Join(tempDir, "httpd.yaml")
	httpdData := "httpd:\n  server_crt: SERVER_CERT_CONTENT\n" +
		"  server_key: SERVER_KEY_CONTENT\n  other_httpd_key: other_httpd_value\n"
	if err := os.WriteFile(httpdPath, []byte(httpdData), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ExtractSecrets(); err != nil {
		t.Errorf("ExtractSecrets failed: %v", err)
	}

	// Verify secrets
	testutils.AssertEquals(t, "Wrong CA secret", "CA_CERT_CONTENT", secrets[podman.CASecret])
	testutils.AssertEquals(t, "Wrong server cert secret", "SERVER_CERT_CONTENT", secrets[podman.ProxySSLCertSecret])
	testutils.AssertEquals(t, "Wrong server key secret", "SERVER_KEY_CONTENT", secrets[podman.ProxySSLKeySecret])

	// Verify files are updated (keys removed)
	data, _ := os.ReadFile(configPath)
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to unmarshal config.yaml: %v", err)
	}
	if _, ok := config["ca_crt"]; ok {
		t.Error("ca_crt was not removed from config.yaml")
	}
	testutils.AssertEquals(t, "other_key was modified", "other_value", config["other_key"])

	data, _ = os.ReadFile(httpdPath)
	var httpdConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &httpdConfig); err != nil {
		t.Fatalf("Failed to unmarshal httpd.yaml: %v", err)
	}

	httpd := ensureStringMap(httpdConfig["httpd"])
	if _, ok := httpd["server_crt"]; ok {
		t.Error("server_crt was not removed from httpd.yaml")
	}
	if _, ok := httpd["server_key"]; ok {
		t.Error("server_key was not removed from httpd.yaml")
	}
	testutils.AssertEquals(t, "other_httpd_key was modified", "other_httpd_value", httpd["other_httpd_key"])
}
