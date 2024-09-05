// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
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
)

// Helm annotation to add in order to use cert-manager's uyuni CA issuer, in JSON format.
const ingressCertManagerAnnotation = "ingressSSLAnnotations={\"cert-manager.io/issuer\": \"uyuni-ca-issuer\"}"

// DeployExistingCertificate execute a deploy of an existing certificate.
func DeployExistingCertificate(
	helmFlags *cmd_utils.HelmFlags,
	sslFlags *cmd_utils.InstallSSLFlags,
) error {
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
		Namespace:   helmFlags.Uyuni.Namespace,
		Name:        "uyuni-cert",
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
	createCaConfig(helmFlags.Uyuni.Namespace, rootCaCrt)
	return nil
}

// DeployReusedCaCertificate deploys an existing SSL CA using cert-manager.
func DeployReusedCa(
	helmFlags *cmd_utils.HelmFlags,
	ca *types.SSLPair,
	kubeconfig string,
	imagePullPolicy string,
) ([]string, error) {
	helmArgs := []string{}

	// Install cert-manager if needed
	if err := installCertManager(helmFlags, kubeconfig, imagePullPolicy); err != nil {
		return []string{}, utils.Errorf(err, L("cannot install cert manager"))
	}

	log.Info().Msg(L("Creating cert-manager issuer for existing CA"))
	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return []string{}, err
	}
	defer cleaner()

	issuerPath := filepath.Join(tempDir, "issuer.yaml")

	issuerData := templates.ReusedCaIssuerTemplateData{
		Namespace:   helmFlags.Uyuni.Namespace,
		Key:         ca.Key,
		Certificate: ca.Cert,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		return []string{}, utils.Errorf(err, L("failed to generate issuer definition"))
	}

	err = utils.RunCmd("kubectl", "apply", "-f", issuerPath)
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create issuer"))
	}

	// Wait for issuer to be ready
	if err := waitForIssuer(helmFlags.Uyuni.Namespace, "uyuni-ca-issuer"); err != nil {
		return nil, err
	}
	helmArgs = append(helmArgs, "--set-json", ingressCertManagerAnnotation)

	// Copy the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	createCaConfig(helmFlags.Uyuni.Namespace, []byte(ca.Cert))

	return helmArgs, nil
}

// DeployGenerateCa deploys a new SSL CA using cert-manager.
func DeployCertificate(
	helmFlags *cmd_utils.HelmFlags,
	sslFlags *cmd_utils.InstallSSLFlags,
	kubeconfig string,
	fqdn string,
	imagePullPolicy string,
) ([]string, error) {
	helmArgs := []string{}

	// Install cert-manager if needed
	if err := installCertManager(helmFlags, kubeconfig, imagePullPolicy); err != nil {
		return []string{}, utils.Errorf(err, L("cannot install cert manager"))
	}

	log.Info().Msg(L("Creating SSL certificate issuer"))
	tempDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return []string{}, utils.Errorf(err, L("failed to create temporary directory"))
	}
	defer os.RemoveAll(tempDir)

	issuerPath := filepath.Join(tempDir, "issuer.yaml")

	issuerData := templates.GeneratedCaIssuerTemplateData{
		Namespace: helmFlags.Uyuni.Namespace,
		Country:   sslFlags.Country,
		State:     sslFlags.State,
		City:      sslFlags.City,
		Org:       sslFlags.Org,
		OrgUnit:   sslFlags.OU,
		Email:     sslFlags.Email,
		Fqdn:      fqdn,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		return []string{}, utils.Errorf(err, L("failed to generate issuer definition"))
	}

	err = utils.RunCmd("kubectl", "apply", "-f", issuerPath)
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create issuer"))
	}

	// Wait for issuer to be ready
	if err := waitForIssuer(helmFlags.Uyuni.Namespace, "uyuni-ca-issuer"); err != nil {
		return nil, err
	}
	helmArgs = append(helmArgs, "--set-json", ingressCertManagerAnnotation)

	// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	extractCaCertToConfig(helmFlags.Uyuni.Namespace)

	return helmArgs, nil
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

func installCertManager(helmFlags *cmd_utils.HelmFlags, kubeconfig string, imagePullPolicy string) error {
	if !kubernetes.IsDeploymentReady("", "cert-manager") {
		log.Info().Msg(L("Installing cert-manager"))
		repo := ""
		chart := helmFlags.CertManager.Chart
		version := helmFlags.CertManager.Version
		namespace := helmFlags.CertManager.Namespace

		args := []string{
			"--set", "crds.enabled=true",
			"--set", "crds.keep=true",
			"--set-json", "global.commonLabels={\"installedby\": \"mgradm\"}",
			"--set", "image.pullPolicy=" + kubernetes.GetPullPolicy(imagePullPolicy),
		}
		extraValues := helmFlags.CertManager.Values
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
	err := kubernetes.WaitForDeployment("", "cert-manager-webhook", "webhook")
	if err != nil {
		return utils.Errorf(err, L("cannot deploy"))
	}

	return nil
}

func extractCaCertToConfig(namespace string) {
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
		return
	}

	out, err = utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "secret", "uyuni-ca", jsonPath, "-n", namespace)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to get uyuni-ca certificate"))
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to base64 decode CA certificate"))
	}

	createCaConfig(namespace, decoded)
}

func createCaConfig(namespace string, ca []byte) {
	valueArg := "--from-literal=ca.crt=" + string(ca)
	if err := utils.RunCmd("kubectl", "create", "configmap", "uyuni-ca", valueArg, "-n", namespace); err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create uyuni-ca config map from certificate"))
	}
}
