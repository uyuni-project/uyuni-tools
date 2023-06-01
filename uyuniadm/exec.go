package main

//TODO required go version >= 18
//zypper in libbtrfs-devel libgpgme-devel device-mapper-devel gpgme libassuan >= 2.5.3
//systemct start podman

import (
	"fmt"
	"os"

	"context"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/spf13/cobra"
)

//TODO avoid global variable
var ctx context.Context
var uyuniContainer container
var cert certificate

var podmanSocket string
var image taggedReference
var configJsonFilename string

func getContext(socket string) context.Context {
	context, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return context
}

func init() {
	rootCmd.AddCommand(fooCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
  //TODO argument should be set somewhere else
	//start
	setFromArgument(startCmd, &configJsonFilename, "config", "/etc/uyuni-tools/options.json", "Config JSON Filename", true)
	setFromArgument(startCmd, &uyuniContainer.name, "name", "", "Container Name", true)
	setFromArgument(startCmd, &uyuniContainer.hostname, "hostname", "", "Hostname", true)
	setFromArgument(startCmd, &admin.user, "user", "spacewalk", "Manager User", true)
	setFromArgument(startCmd, &admin.password, "password", "spacewalk", "Manager Password", true)
	setFromArgument(startCmd, &admin.email, "email", "galaxy-noise@suse.de", "Manager E-Mail", true)
	setFromArgument(startCmd, &admin.db.name, "dbname", "susemanager", "Database Name", true)
	setFromArgument(startCmd, &admin.db.host, "dbhost", "localhost", "Database Host", true)
	setFromArgument(startCmd, &admin.db.port, "dbport", "5432", "Database Port", true)
	setFromArgument(startCmd, &admin.db.protocol, "dbprotocol", "tcp", "Database Protocol", true)
	setFromArgument(startCmd, &enableTftp, "tftp", "Y", "Enable TFTP", true)
	setFromArgument(startCmd, &reportDB.host, "reportdbhost", "", "Report Database Host", true)
	setFromArgument(startCmd, &reportDB.password, "reportdbpassword", "pythia_susemanager", "Report Database Password", true)
	setFromArgument(startCmd, &image.namedRepository, "namespace", "registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server", "Namespace Image", true)
	setFromArgument(startCmd, &image.tag, "tag", "latest", "Tag Image", true)
	setCertificate(startCmd, &cert)
	setFromArgument(startCmd, &podmanSocket, "podman_socket", "unix://run/podman/podman.sock", "Podman Socket", true)

	//stop
	setFromArgument(stopCmd, &uyuniContainer.name, "name", "uyuni-server", "Container Name", true)

	//migrate
	setFromArgument(migrateCmd, &configJsonFilename, "config", "/etc/uyuni-tools/options.json", "Config JSON Filename", true)
	setFromArgument(migrateCmd, &uyuniContainer.name, "name", "", "Container Name", true)
	setFromArgument(migrateCmd, &source.hostname, "server", "", "Hostname of the source server", true)
	setFromArgument(migrateCmd, &source.user, "user", "", "User used to login in source server", true)
	setFromArgument(migrateCmd, &image.namedRepository, "image", "registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server", "Image", true)
	setFromArgument(migrateCmd, &image.tag, "tag", "latest", "Tag Image", true)
	setCertificate(migrateCmd, &cert)
	setFromArgument(migrateCmd, &podmanSocket, "podman_socket", "unix://run/podman/podman.sock", "Podman Socket", true)
	setFromArgument(migrateCmd, &sshAuthSocket, "ssh_auth_socket", "", "SSH_AUTH_SOCKET value", true)

	ctx = getContext(podmanSocket)
}

func setCertificate(cmd *cobra.Command, cert *certificate) {
	//TODO set useExistingCertificate and the rest mutual exclusive
	setFromArgument(cmd, &cert.useExistingCertificate, "cert_use_existing", "N", "Use Existing Certificate", false)
	setFromArgument(cmd, &cert.country, "cert_country", "DE", "SSL Certificate Country", false)
	setFromArgument(cmd, &cert.state, "cert_state", "Bayern", "SSL Certificate State", false)
	setFromArgument(cmd, &cert.city, "cert_city", "Nuernberg", "SSL Certificate City", false)
	setFromArgument(cmd, &cert.org, "cert_org", "SUSE", "SSL Certificate Organization", false)
	setFromArgument(cmd, &cert.ou, "cert_ou", "SUSE", "SSL Certificate Organization Unit", false)
	setFromArgument(cmd, &cert.pwd, "cert_pwd", "spacewalk", "SSL Certificate Password", false)
	setFromArgument(cmd, &cert.email, "cert_email", "galaxy-noise@suse.de", "SSL Certificate E-Mail", false)

}

func setFromArgument(cmd *cobra.Command, value *string, arg string, defaultValue string, usage string, mandatory bool) {
	cmd.Flags().StringVar(value, arg, defaultValue, usage)
	if mandatory && len(defaultValue) == 0 {
		cmd.MarkFlagRequired(arg)
	}
	return
}
