// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
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
	err := os.MkdirAll(tempDir, 0755)
	testutils.AssertNoError(t, "failed to make temporary directory", err)
	err = checkPermissions(tempDir, 0005|0050|0500)
	testutils.AssertNoError(t, "failed to check permissions", err)
}

func TestValidateYamlFiles(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{"httpd.yaml", "ssh.yaml", "config.yaml"}
	for _, file := range testFiles {
		filePath := path.Join(tempDir, file)
		_, err := os.Create(filePath)
		testutils.AssertNoError(t, "Failed to create test file "+filePath, err)
	}

	// Test: when all files are present and have correct permissions
	testutils.AssertNoError(t, "Expected no error", validateInstallYamlFiles(tempDir))

	// Test: Missing file scenario, remove one file and expect an error
	os.Remove(path.Join(tempDir, "httpd.yaml"))
	testutils.AssertError(t, "missing required configuration file:", validateInstallYamlFiles(tempDir))
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
				"Environment=HTTPD_EXTRA_CONF=-v/my/apache.conf:/etc/apache2/conf.d/zz_apache_tuning.conf:ro",
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
	testutils.AssertNoError(t, "Unexpected error", err)
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
	testutils.WriteFile(t, configPath, configData)

	httpdPath := path.Join(tempDir, "httpd.yaml")
	httpdData := "httpd:\n  server_crt: SERVER_CERT_CONTENT\n" +
		"  server_key: SERVER_KEY_CONTENT\n  other_httpd_key: other_httpd_value\n"
	testutils.WriteFile(t, httpdPath, httpdData)

	testutils.AssertNoError(t, "ExtractSecrets failed", ExtractSecrets())

	// Verify secrets
	testutils.AssertEquals(t, "Wrong CA secret", "CA_CERT_CONTENT", secrets[podman.CASecret])
	testutils.AssertEquals(t, "Wrong server cert secret", "SERVER_CERT_CONTENT", secrets[podman.ProxySSLCertSecret])
	testutils.AssertEquals(t, "Wrong server key secret", "SERVER_KEY_CONTENT", secrets[podman.ProxySSLKeySecret])

	// Verify files are updated (keys removed)
	data, _ := os.ReadFile(configPath)
	var config map[string]interface{}
	err := yaml.Unmarshal(data, &config)
	testutils.AssertNoError(t, "Failed to unmarshal config.yaml", err)
	_, ok := config["ca_crt"]
	testutils.AssertTrue(t, "ca_crt was not removed from config.yaml", !ok)
	testutils.AssertEquals(t, "other_key was modified", "other_value", config["other_key"])

	data, _ = os.ReadFile(httpdPath)
	var httpdConfig map[string]interface{}
	err = yaml.Unmarshal(data, &httpdConfig)
	testutils.AssertNoError(t, "Failed to unmarshal httpd.yaml", err)

	httpd := ensureStringMap(httpdConfig["httpd"])
	_, ok = httpd["server_crt"]
	testutils.AssertTrue(t, "server_crt was not removed from httpd.yaml", !ok)
	_, ok = httpd["server_key"]
	testutils.AssertTrue(t, "server_key was not removed from httpd.yaml", !ok)
	testutils.AssertEquals(t, "other_httpd_key was modified", "other_httpd_value", httpd["other_httpd_key"])
}
