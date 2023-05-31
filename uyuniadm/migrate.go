package main

//TODO required go version >= 18
//zypper in libbtrfs-devel libgpgme-devel device-mapper-devel gpgme libassuan >= 2.5.3
//systemct start podman

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/spf13/cobra"
)

//TODO avoid global variable
var source remoteServer
var sshAuthSocket string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate command",
	Run: func(cmd *cobra.Command, arg []string) {
		cmd.OutOrStdout()
		cmd.OutOrStderr()
		setEnvMigrateParam()
		fmt.Println("Migrating server")
		opt := readJsonOpt(configJsonFilename)
		//TODO to allow uyuniadm to handle more than one container, we should find another way to pass the env variable
		//option from argument and specific for migrate container
		entrypoint := "/var/lib/uyuni-tools/setup-migration-container.sh"
		opt.Entrypoint = &entrypoint
		opt.Name = uyuniContainer.name
		opt.ConmonPIDFile = opt.Name + ".pid"
		opt.CIDFile = opt.Name + ".ctr-id"
		opt.Hostname = "tmp"
		//TODO we don't want to mount this folder
		opt.Volume = append(opt.Volume, filepath.Dir(sshAuthSocket)+":"+filepath.Dir(sshAuthSocket))
		//TODO we don't want to mount this folder
		opt.Volume = append(opt.Volume, "/var/lib/uyuni-tools/:/var/lib/uyuni-tools/")

		getImage(ctx, image.String())
		createContainer(opt, image.String())

		startContainer(ctx, opt.Name)
		if _, err := containers.Wait(ctx, opt.Name, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Server migrated")
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
func setEnvMigrateParam() {
	os.Setenv("REMOTE_USER", source.user)
	os.Setenv("UYUNI_FQDN", source.hostname)
	os.Setenv("CERT_CITY", cert.city)
	os.Setenv("CERT_COUNTRY", cert.country)
	os.Setenv("CERT_STATE", cert.state)
	os.Setenv("CERT_O", cert.org)
	os.Setenv("CERT_OU", cert.ou)
	os.Setenv("CERT_PASS", cert.pwd)
	os.Setenv("CERT_EMAIL", cert.email)
	os.Setenv("SSH_AUTH_SOCK", sshAuthSocket)
}
