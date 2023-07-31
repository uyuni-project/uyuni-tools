package install

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func waitForSystemStart(viper *viper.Viper, globalFlags *types.GlobalFlags) {
	// Setup the systemd service configuration options
	image := fmt.Sprintf("%s:%s", viper.GetString("image"), viper.GetString("tag"))
	podman.GenerateSystemdService(viper.GetString("tz"), image, viper.GetStringSlice("podman.arg"), globalFlags.Verbose)

	log.Println("Waiting for the server to start...")
	// Start the service
	if err := exec.Command("systemctl", "enable", "--now", "uyuni-server").Run(); err != nil {
		log.Fatalf("Failed to enable uyuni-server systemd service: %s\n", err)
	}

	utils.WaitForServer()
}

func pullImage(viper *viper.Viper) {
	image := fmt.Sprintf("%s:%s", viper.GetString("image"), viper.GetString("tag"))
	log.Printf("Running podman pull %s\n", image)
	cmd := exec.Command("podman", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to pull image: %s\n", err)
	}
}

func installForPodman(viper *viper.Viper, globalFlags *types.GlobalFlags, cmd *cobra.Command, args []string) {
	pullImage(viper)

	waitForSystemStart(viper, globalFlags)

	env := []string{
		"MANAGER_USER=" + viper.GetString("db.user"),
		"MANAGER_PASS=" + viper.GetString("db.password"),
		"SCC_USER=" + viper.GetString("scc.user"),
		"SCC_PASS=" + viper.GetString("scc.password"),
		"REPORT_DB_USER=" + viper.GetString("reportdb.user"),
		"REPORT_DB_PASS=" + viper.GetString("reportdb.password"),
		"EXTERNALDB_ADMIN_USER=" + viper.GetString("db.admin.user"),
		"EXTERNALDB_ADMIN_PASS=" + viper.GetString("db.admin.password"),
	}

	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		args[0],
	}

	localDb := utils.Contains(localHostValues, viper.GetString("db.host"))

	dbHost := viper.GetString("db.host")
	reportdbHost := viper.GetString("reportdb.host")

	if localDb {
		// For now the setup script expects the localhost value for local DB
		// but the FQDN is required for the report db even if it's local
		dbHost = "localhost"
		if viper.GetString("reportdb.host") == "" {
			reportdbHost = args[0]
		}
	}

	execArgs := []string{
		"exec",
		"-e", "UYUNI_FQDN=" + args[0],
		"-e", "MANAGER_USER",
		"-e", "MANAGER_PASS",
		"-e", "MANAGER_ADMIN_EMAIL=" + viper.GetString("email"),
		"-e", "MANAGER_MAIL_FROM=" + viper.GetString("emailFrom"),
		"-e", "MANAGER_ENABLE_TFTP=" + boolToString(viper.GetBool("enableTftp")),
		"-e", "LOCAL_DB=" + boolToString(localDb),
		"-e", "MANAGER_DB_NAME=" + viper.GetString("db.name"),
		"-e", "MANAGER_DB_HOST=" + dbHost,
		"-e", "MANAGER_DB_PORT=" + strconv.Itoa(viper.GetInt("db.port")),
		"-e", "MANAGER_DB_PROTOCOL=" + viper.GetString("db.protocol"),
		"-e", "REPORT_DB_NAME=" + viper.GetString("reportdb.name"),
		"-e", "REPORT_DB_HOST=" + reportdbHost,
		"-e", "REPORT_DB_PORT=" + strconv.Itoa(viper.GetInt("reportdb.port")),
		"-e", "REPORT_DB_USER",
		"-e", "REPORT_DB_PASS",
		"-e", "EXTERNALDB_ADMIN_USER",
		"-e", "EXTERNALDB_ADMIN_PASS",
		"-e", "EXTERNALDB_PROVIDER=" + viper.GetString("db.provider"),
		"-e", "ISS_PARENT=" + viper.GetString("issParent"),
		"-e", "MIRROR_PATH=" + viper.GetString("mirrorPath"),
		"-e", "ACTIVATE_SLP=N", // Deprecated, will be removed soon
		"-e", "SCC_USER",
		"-e", "SCC_PASS",
	}

	if viper.GetBool("cert.useexisting") {
		// TODO Get existing certificates path and mount them
		// Set CA_CERT, SERVER_CERT, SERVER_KEY or run the rhn-ssl-check tool in a container
		// The SERVER_CERT needs to get the intermediate keys
	} else {
		execArgs = append(execArgs,
			"-e", "CERT_O="+viper.GetString("cert.org"),
			"-e", "CERT_OU="+viper.GetString("cert.ou"),
			"-e", "CERT_CITY="+viper.GetString("cert.city"),
			"-e", "CERT_STATE="+viper.GetString("cert.state"),
			"-e", "CERT_COUNTRY="+viper.GetString("cert.country"),
			"-e", "CERT_EMAIL="+viper.GetString("cert.email"),
			"-e", "CERT_CNAMES="+strings.Join(append([]string{args[0]}, viper.GetStringSlice("cert.cnames")...), ","),
			"-e", "CERT_PASS",
		)
		env = append(env, "CERT_PASS="+viper.GetString("cert.password"))
	}

	execArgs = append(execArgs, "uyuni-server", "/usr/lib/susemanager/bin/mgr-setup", "-s", "-n")
	if globalFlags.Verbose {
		fmt.Printf("> Running: %s %s\n", "podman", strings.Join(execArgs, " "))
	}

	execCmd := exec.Command("podman", execArgs...)
	execCmd.Env = append(execCmd.Environ(), env...)
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout

	if err := execCmd.Run(); err != nil {
		log.Fatalf("Failed to setup the server: %s\n", err)
	}

	log.Println("Server set up")
}

func boolToString(value bool) string {
	if value {
		return "Y"
	}
	return "N"
}
