// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"errors"
	"os"
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
	// Deploy the SSL Certificate secret and CA configmap
	serverCrt, rootCaCrt := ssl.OrderCas(&sslFlags.Ca, &sslFlags.Server)
	serverKey := utils.ReadFile(sslFlags.Server.Key)

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	secretPath := filepath.Join(tempDir, "secret.yaml")
	log.Info().Msg(L("Creating SSL server certificate secret"))
	tlsSecretData := templates.TLSSecretTemplateData{
		Namespace:   namespace,
		Name:        CertSecretName,
		Certificate: base64.StdEncoding.EncodeToString(serverCrt),
		Key:         base64.StdEncoding.EncodeToString(serverKey),
		RootCa:      base64.StdEncoding.EncodeToString(rootCaCrt),
	}

	if err = utils.WriteTemplateToFile(tlsSecretData, secretPath, 0500, true); err != nil {
		return utils.Errorf(err, L("Failed to generate uyuni-crt secret definition"))
	}
	err = utils.RunCmd("kubectl", "apply", "-f", secretPath)
	if err != nil {
		return utils.Errorf(err, L("Failed to create uyuni-crt TLS secret"))
	}

	// Copy the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	return createCaConfig(namespace, rootCaCrt)
}

// DeployReusedCa deploys an existing SSL CA using an already installed cert-manager.
func DeployReusedCa(namespace string, ca *types.SSLPair) error {
	log.Info().Msg(L("Creating cert-manager issuer for existing CA"))
	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	issuerPath := filepath.Join(tempDir, "issuer.yaml")

	issuerData := templates.ReusedCaIssuerTemplateData{
		Namespace:   namespace,
		Key:         ca.Key,
		Certificate: ca.Cert,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		return utils.Errorf(err, L("failed to generate issuer definition"))
	}

	err = utils.RunCmd("kubectl", "apply", "-f", issuerPath)
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create issuer"))
	}

	return nil
}

// DeployGenerateCa deploys a new SSL CA using cert-manager.
func DeployGeneratedCa(
	namespace string,
	sslFlags *cmd_utils.InstallSSLFlags,
	fqdn string,
) error {
	log.Info().Msg(L("Creating SSL certificate issuer"))
	tempDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return utils.Errorf(err, L("failed to create temporary directory"))
	}
	defer os.RemoveAll(tempDir)

	issuerPath := filepath.Join(tempDir, "issuer.yaml")

	issuerData := templates.GeneratedCaIssuerTemplateData{
		Namespace: namespace,
		Country:   sslFlags.Country,
		State:     sslFlags.State,
		City:      sslFlags.City,
		Org:       sslFlags.Org,
		OrgUnit:   sslFlags.OU,
		Email:     sslFlags.Email,
		Fqdn:      fqdn,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		return utils.Errorf(err, L("failed to generate issuer definition"))
	}

	err = utils.RunCmd("kubectl", "apply", "-f", issuerPath)
	if err != nil {
		return utils.Errorf(err, L("Failed to create issuer"))
	}

	return nil
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
			return utils.Errorf(err, L("cannot run helm upgrade"))
		}
	}

	// Wait for cert-manager to be ready
	err := kubernetes.WaitForDeployments("", "cert-manager-webhook")
	if err != nil {
		return utils.Errorf(err, L("cannot deploy"))
	}

	return nil
}

func extractCaCertToConfig(namespace string) error {
	// TODO Replace with [trust-manager](https://cert-manager.io/docs/projects/trust-manager/) to automate this
	const jsonPath = "-o=jsonpath={.data.ca\\.crt}"

	log.Info().Msg(L("Extracting CA certificate to a configmap"))
	// Skip extracting if the configmap is already present
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "configmap", "uyuni-ca", jsonPath, "-n", namespace,
	)
	log.Info().Msgf(L("CA cert: %s"), string(out))
	if err == nil && len(out) > 0 {
		log.Info().Msg(L("uyuni-ca configmap already existing, skipping extraction"))
		return nil
	}

	out, err = utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "secret", "-n", namespace, "uyuni-ca", jsonPath,
	)
	if err != nil {
		return utils.Errorf(err, L("Failed to get uyuni-ca certificate"))
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return utils.Errorf(err, L("Failed to base64 decode CA certificate"))
	}

	return createCaConfig(namespace, decoded)
}

func createCaConfig(namespace string, ca []byte) error {
	configMap := core.ConfigMap{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      "uyuni-ca",
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, ""),
		},
		Data: map[string]string{
			"ca.crt": string(ca),
		},
	}
	return kubernetes.Apply([]runtime.Object{&configMap}, L("failed to create the SSH migration ConfigMap"))
}

// HasIssuer returns true if the issuer is defined.
//
// False will be returned in case of errors or if the issuer resource doesn't exist on the cluster.
func HasIssuer(namespace string, name string) bool {
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "issuer", "-n", namespace, name, "-o", "name")
	return err == nil && strings.TrimSpace(string(out)) != ""
}
