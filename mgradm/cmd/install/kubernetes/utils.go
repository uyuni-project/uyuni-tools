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
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(_ *types.GlobalFlags,
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

	if err := shared_utils.IsValidFQDN(fqdn); err != nil {
		return err
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
	if flags.SSL.UseExisting() {
		if err := kubernetes.DeployExistingCertificate(&flags.Helm, &flags.SSL); err != nil {
			return err
		}
	} else {
		sslArgs, err := kubernetes.DeployCertificate(
			&flags.Helm, &flags.SSL, clusterInfos.GetKubeconfig(), fqdn,
			flags.Image.PullPolicy,
		)

		if err != nil {
			return shared_utils.Errorf(err, L("cannot deploy certificate"))
		}
		helmArgs = append(helmArgs, sslArgs...)
	}

	// Create a secret using SCC credentials if any are provided
	helmArgs, err = shared_kubernetes.AddSCCSecret(helmArgs, flags.Helm.Uyuni.Namespace, &flags.SCC)
	if err != nil {
		return err
	}

	// Deploy Uyuni and wait for it to be up
	if err := kubernetes.Deploy(
		cnx, flags.Image.Registry, &flags.Image, &flags.HubXmlrpc, &flags.Helm,
		clusterInfos, fqdn, flags.Debug.Java, false, helmArgs...,
	); err != nil {
		return shared_utils.Errorf(err, L("cannot deploy uyuni"))
	}

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}

	if err := install_shared.RunSetup(cnx, &flags.InstallFlags, args[0], envs); err != nil {
		namespace, err := cnx.GetNamespace("")
		if err != nil {
			return shared_utils.Errorf(err, L("failed to stop service"))
		}
		if stopErr := shared_kubernetes.Stop(namespace, shared_kubernetes.ServerApp); stopErr != nil {
			log.Error().Err(stopErr).Msg(L("failed to stop service"))
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
