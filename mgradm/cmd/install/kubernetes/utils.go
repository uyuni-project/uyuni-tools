// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	install_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(globalFlags *types.GlobalFlags,
	flags *kubernetesInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	flags.CheckParameters(cmd, "kubectl")
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)

	fqdn := args[0]

	if !shared_utils.IsValidFQDN(fqdn) {
		return fmt.Errorf(L("%s is not a valid FDQN"), fqdn)
	}

	helmArgs := []string{"--set", "timezone=" + flags.TZ}
	if flags.Mirror != "" {
		// TODO Handle claims for multi-node clusters
		helmArgs = append(helmArgs, "--set", "mirror.hostPath="+flags.Mirror)
	}
	if flags.Debug.Java {
		helmArgs = append(helmArgs, "--set", "exposeJavaDebug=true")
	}

	// Check the kubernetes cluster setup
	clusterInfos, err := shared_kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	// Deploy the SSL CA or server certificate
	ca := ssl.SslPair{}
	sslArgs, err := kubernetes.DeployCertificate(&flags.Helm, &flags.Ssl, "", &ca, clusterInfos.GetKubeconfig(), fqdn,
		flags.Image.PullPolicy)
	if err != nil {
		return shared_utils.Errorf(err, L("cannot deploy certificate"))
	}
	helmArgs = append(helmArgs, sslArgs...)

	// Deploy Uyuni and wait for it to be up
	if err := kubernetes.Deploy(cnx, globalFlags.Registry, &flags.Image, &flags.Helm, &flags.Ssl,
		clusterInfos, fqdn, flags.Debug.Java, helmArgs...,
	); err != nil {
		return shared_utils.Errorf(err, L("cannot deploy uyuni"))
	}

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}

	if err := install_shared.RunSetup(cnx, &flags.InstallFlags, args[0], envs); err != nil {
		if stopErr := shared_kubernetes.Stop(shared_kubernetes.ServerFilter); stopErr != nil {
			log.Error().Msgf(L("Failed to stop service: %v"), stopErr)
		}
		return err
	}

	// The CA needs to be added to the database for Kickstart use.
	err = adm_utils.ExecCommand(zerolog.DebugLevel, cnx,
		"/usr/bin/rhn-ssl-dbstore", "--ca-cert=/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT")
	if err != nil {
		return shared_utils.Errorf(err, L("error storing the SSL CA certificate in database"))
	}
	return nil
}
