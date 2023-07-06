package podman

import (
	"log"
	"os"
	"regexp"
)

func ReadConfig() []byte {
	configPath := getServiceConfig()
	config, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read %s config file: %s\n", configPath, err)
	}
	return config
}

func WriteConfig(config []byte) {
	configPath := getServiceConfig()
	if err := os.WriteFile(configPath, config, 0644); err != nil {
		log.Fatalf("Failed to save configuration changes to %s: %s\n", configPath, err)
	}
}

func UpdateConfigValue(config []byte, key string, value string) []byte {
	var matcher = regexp.MustCompile(key + ` ?=.*`)
	return matcher.ReplaceAll(config, []byte(key+"="+value))
}

func getServiceConfig() string {
	// TODO Adjust for othe distros
	return "/etc/sysconfig/uyuni-server-systemd-services"
}
