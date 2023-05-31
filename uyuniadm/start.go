package main

//TODO required go version >= 18
//zypper in libbtrfs-devel libgpgme-devel device-mapper-devel gpgme libassuan >= 2.5.3
//systemct start podman

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/spf13/cobra"
)

//TODO avoid global variable
var admin manager

var enableTftp string
var reportDB database

func startContainer(context context.Context, name string) {

	if err := containers.Start(context, name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Container started.")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start command",
	Run: func(cmd *cobra.Command, arg []string) {
		cmd.OutOrStdout()
		cmd.OutOrStderr()
		setEnvStartParam()
		fmt.Println("Starting container")

		opt := readJsonOpt(configJsonFilename)

		//TODO to allow uyuniadm to handle more than one container, we should find another way to pass the env variable
		//option from argument and specific for start container
		opt.Hostname = uyuniContainer.hostname
		opt.Name = uyuniContainer.name
		opt.ConmonPIDFile = uyuniContainer.name + ".pid"
		opt.CIDFile = uyuniContainer.name + ".ctr-id"

		getImage(ctx, image.String())
		createContainer(opt, image.String())

		startContainer(ctx, opt.Name)
	},
}

/*TODO to allow uyuniadm to handle more than one container, we should find another way to pass the env variable
*  - do not hardcore env variable
*  - Possible solution: option.json as single point of truth for env variable:
*    - we parse all the env variable from option.json
*    - each run of uyuniadm would save the variable with a unique prefix
*    - before passing options to createContainer, we assign values of env variable to container (removing prefix)
*  - this function should be uniq for all commands
 */

func setEnvStartParam() {
	os.Setenv("MANAGER_USER", admin.user)
	os.Setenv("MANAGER_PASS", admin.password)
	os.Setenv("MANAGER_ADMIN_EMAIL", admin.email)
	os.Setenv("MANAGER_DB_NAME", admin.db.name)
	os.Setenv("MANAGER_DB_HOST", admin.db.host)
	os.Setenv("MANAGER_DB_PORT", admin.db.port)
	os.Setenv("MANAGER_DB_PROTOCOL", admin.db.port)
	os.Setenv("REPORT_DB_HOST", reportDB.host)
	os.Setenv("REPORT_DB_PASS", reportDB.password)
	os.Setenv("MANAGER_ENABLE_TFTP", enableTftp)
	os.Setenv("USE_EXISTING_CERTS", cert.useExistingCertificate)
	os.Setenv("CERT_CITY", cert.city)
	os.Setenv("CERT_COUNTRY", cert.country)
	os.Setenv("CERT_STATE", cert.state)
	os.Setenv("CERT_O", cert.org)
	os.Setenv("CERT_OU", cert.ou)
	os.Setenv("CERT_PASS", cert.pwd)
	os.Setenv("CERT_EMAIL", cert.email)
}
