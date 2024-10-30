// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeployExistingCertificate execute a deploy of an existing certificate.
func DeployExistingCertificate(namespace string, sslFlags *cmd_utils.InstallSSLFlags) error {
	if err := createTLSCertificate(
		namespace, kubernetes.CASecretName, kubernetes.CAConfigName, &sslFlags.Ca, &sslFlags.Server,
	); err != nil {
		return err
	}

	// Handle the DB certificate
	var dbCA *types.CaChain
	var dbPair *types.SSLPair
	if sslFlags.UseProvidedDB() {
		dbCA = &sslFlags.DB.CA
		dbPair = &sslFlags.DB.SSLPair
	} else if sslFlags.DB.IsDefined() && !sslFlags.DB.CA.IsThirdParty() {
		//let's use the CA already present in the server
		dbCA = &sslFlags.Ca
		dbPair = &sslFlags.DB.SSLPair
	} else {
		return errors.New(L("Database SSL certificate and key have to be defined"))
	}

	return createTLSCertificate(namespace, kubernetes.DBCertSecretName, kubernetes.DBCAConfigName, dbCA, dbPair)
}

func createTLSCertificate(
	namespace string,
	secretName string,
	caConfigName string,
	ca *types.CaChain,
	certPair *types.SSLPair,
) error {
	// Deploy the SSL Certificate secret and CA ConfigMap
	serverCrt, rootCaCrt, err := ssl.OrderCas(ca, certPair)
	if err != nil {
		return err
	}
	serverKey := utils.ReadFile(certPair.Key)

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	secretPath := filepath.Join(tempDir, "secret.yaml")
	log.Info().Msg(L("Creating SSL server certificate secret"))
	tlsSecretData := templates.TLSSecretTemplateData{
		Namespace:   namespace,
		Name:        secretName,
		Certificate: base64.StdEncoding.EncodeToString(serverCrt),
		Key:         base64.StdEncoding.EncodeToString(serverKey),
		RootCa:      base64.StdEncoding.EncodeToString(rootCaCrt),
	}

	if err = utils.WriteTemplateToFile(tlsSecretData, secretPath, 0500, true); err != nil {
		return utils.Errorf(err, L("Failed to generate %s secret definition"), secretName)
	}
	err = utils.RunCmd("kubectl", "apply", "-f", secretPath)
	if err != nil {
		return utils.Errorf(err, L("Failed to create %s TLS secret"), secretName)
	}

	// Copy the CA cert into a ConfigMap for containers who shouldn't see the key
	return createCAConfig(namespace, caConfigName, rootCaCrt)
}

// DeployReusedCA deploys an existing SSL CA using an already installed cert-manager.
func DeployReusedCA(namespace string, ca *types.SSLPair, fqdn string) error {
	log.Info().Msg(L("Creating cert-manager issuer for existing CA"))
	return templates.NewReusedCAIssuerTemplate(namespace, fqdn, ca.Cert, ca.Key).Apply()
}

// DeployGenerateCA deploys a new SSL CA using cert-manager.
func DeployGeneratedCA(
	namespace string,
	sslFlags *cmd_utils.InstallSSLFlags,
	fqdn string,
) error {
	log.Info().Msg(L("Creating SSL certificate issuer"))

	return templates.NewGeneratedCAIssuerTemplate(
		namespace,
		fqdn,
		sslFlags.Country,
		sslFlags.State,
		sslFlags.City,
		sslFlags.Org,
		sslFlags.OU,
		sslFlags.Email,
	).Apply()
}

// Wait for issuer to be ready.
func waitForIssuer(namespace string, name string) error {
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(
			zerolog.DebugLevel, "kubectl", "get",
			"-o=jsonpath={.status.conditions[*].type}",
			"-n", namespace,
			"issuer", name,
		)
		if err == nil && string(out) == "Ready" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New(L("Issuer didn't turn ready after 60s"))
}

// InstallCertManager deploys the cert-manager helm chart with the CRDs.
func InstallCertManager(kubernetesFlags *cmd_utils.KubernetesFlags, kubeconfig string, imagePullPolicy string) error {
	if ready, err := kubernetes.IsDeploymentReady("", "cert-manager"); err != nil {
		return err
	} else if !ready {
		log.Info().Msg(L("Installing cert-manager"))
		repo := ""
		chart := kubernetesFlags.CertManager.Chart
		version := kubernetesFlags.CertManager.Version
		namespace := kubernetesFlags.CertManager.Namespace

		args := []string{
			"--set", "crds.enabled=true",
			"--set", "crds.keep=true",
			"--set-json", "global.commonLabels={\"installedby\": \"mgradm\"}",
			"--set", "image.pullPolicy=" + string(kubernetes.GetPullPolicy(imagePullPolicy)),
		}
		extraValues := kubernetesFlags.CertManager.Values
		if extraValues != "" {
			args = append(args, "-f", extraValues)
		}

		// Use upstream chart if nothing defined
		if chart == "" {
			repo = "https://charts.jetstack.io"
			chart = "cert-manager"
		}
		// The installedby label will be used to only uninstall what we installed
		if err := kubernetes.HelmUpgrade(
			kubeconfig, namespace, true, repo, "cert-manager", chart, version, args...,
		); err != nil {
			return utils.Error(err, L("cannot run helm upgrade"))
		}
	}

	// Wait for cert-manager to be ready
	err := kubernetes.WaitForDeployments("", "cert-manager-webhook")
	if err != nil {
		return utils.Error(err, L("cannot deploy"))
	}

	return nil
}

func extractCACertToConfig(namespace string) error {
	// TODO Replace with [trust-manager](https://cert-manager.io/docs/projects/trust-manager/) to automate this
	const jsonPath = "-o=jsonpath={.data.ca\\.crt}"

	log.Info().Msg(L("Extracting CA certificate to a ConfigMap"))
	// Skip extracting if the configmap is already present
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "configmap", kubernetes.CAConfigName, jsonPath, "-n", namespace,
	)
	log.Info().Msgf(L("CA cert: %s"), string(out))
	if err == nil && len(out) > 0 {
		log.Info().Msgf(L("%s ConfigMap already existing, skipping extraction"), kubernetes.CAConfigName)
		return nil
	}

	out, err = utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "secret", "-n", namespace, kubernetes.CAConfigName, jsonPath,
	)
	if err != nil {
		return utils.Errorf(err, L("Failed to get %s certificate"), kubernetes.CAConfigName)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return utils.Error(err, L("failed to base64 decode CA certificate"))
	}

	// Copy the CA to a ConfigMap as the secret shouldn't be available to the server
	if err := createCAConfig(namespace, kubernetes.CAConfigName, decoded); err != nil {
		return err
	}
	// Also copy the CA to a separate ConfigMap as we would be expecting it for the setup and server containers
	return createCAConfig(namespace, kubernetes.DBCAConfigName, decoded)
}

func createCAConfig(namespace string, name string, ca []byte) error {
	configMap := core.ConfigMap{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, ""),
		},
		Data: map[string]string{
			"ca.crt": string(ca),
		},
	}
	return kubernetes.Apply([]runtime.Object{&configMap}, fmt.Sprintf(L("failed to create the %s ConfigMap"), name))
}

// HasIssuer returns true if the issuer is defined.
//
// False will be returned in case of errors or if the issuer resource doesn't exist on the cluster.
func HasIssuer(namespace string, name string) bool {
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "issuer", "-n", namespace, name, "-o", "name")
	return err == nil && strings.TrimSpace(string(out)) != ""
}
