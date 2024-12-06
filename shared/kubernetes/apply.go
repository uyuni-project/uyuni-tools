// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

// Apply runs kubectl apply for the provided objects.
//
// The message should be a user-friendly localized message to provide in case of error.
func Apply[T runtime.Object](objects []T, message string) error {
	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	// Run the job
	definitionPath := path.Join(tempDir, "definition.yaml")
	if err := YamlFile(objects, definitionPath); err != nil {
		return err
	}

	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "kubectl", "apply", "-f", definitionPath); err != nil {
		return utils.Errorf(err, message)
	}
	return nil
}

// YamlFile generates a YAML file from a list of kubernetes objects.
func YamlFile[T runtime.Object](objects []T, path string) error {
	printer := printers.YAMLPrinter{}
	file, err := os.Create(path)
	if err != nil {
		return utils.Errorf(err, L("failed to create %s YAML file"), path)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msgf(L("failed to close %s YAML file"), path)
		}
	}()

	for _, obj := range objects {
		err = printer.PrintObj(obj, file)
		if err != nil {
			return utils.Errorf(err, L("failed to write PVC to file"))
		}
	}

	return nil
}
