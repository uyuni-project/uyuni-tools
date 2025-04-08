// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restore

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func restoreSystemdConfiguration(backupSource string, flags *shared.Flagpole) error {
	backupFile, err := os.Open(backupSource)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	var hasError error

	tr := tar.NewReader(backupFile)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if flags.DryRun {
			log.Info().Msgf(L("Would restore %s"), header.Name)
			continue
		}

		log.Debug().Msgf("Restoring systemd file %s", header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(header.Name, header.FileInfo().Mode()); err != nil {
				log.Warn().Msgf(L("Unable to create directory %s"), header.Name)
				hasError = utils.JoinErrors(hasError, err)
				continue
			}
		case tar.TypeReg:
			if err := restoreSystemdFile(header, tr); err != nil {
				hasError = utils.JoinErrors(hasError, err)
				continue
			}
		default:
			log.Warn().Msgf(L("Unknown filetype of %s"), header.Name)
			continue
		}

		if err := restoreFileAttributes(header.Name, header); err != nil {
			log.Warn().Err(err).Msgf(L("Unable to restore file details for %s"), header.Name)
			hasError = utils.JoinErrors(hasError, err)
		}
	}
	return hasError
}

func restoreSystemdFile(header *tar.Header, tr *tar.Reader) error {
	fh, err := os.Create(header.Name)
	if err != nil {
		log.Warn().Err(err).Msgf(L("Unable to create %s"), header.Name)
		return err
	}
	if _, err := io.Copy(fh, tr); err != nil {
		log.Warn().Err(err).Msgf(L("Unable to restore content of %s"), header.Name)
		fh.Close()
		os.Remove(header.Name)
		return err
	}
	fh.Close()
	return nil
}

func restoreFileAttributes(filename string, th *tar.Header) error {
	var e error
	e = utils.JoinErrors(e, os.Chmod(filename, th.FileInfo().Mode()))
	e = utils.JoinErrors(e, os.Chown(filename, th.Uid, th.Gid))
	e = utils.JoinErrors(e, os.Chtimes(filename, th.AccessTime, th.ModTime))
	return e
}

func generateDefaltSystemdServices(flags *shared.Flagpole) error {
	if flags.DryRun {
		log.Info().Msg(L("Would generate default systemd services"))
		return nil
	}
	// Generate minimum set - uyuni-db and uyuni-server services - like we do on default install
	serverImage := fmt.Sprintf("%s%s:%s", utils.ServerImage.Registry, utils.ServerImage.Name, utils.ServerImage.Tag)
	dbImage := fmt.Sprintf("%s%s:%s",
		utils.PostgreSQLImage.Registry,
		utils.PostgreSQLImage.Name,
		utils.PostgreSQLImage.Tag)

	return utils.JoinErrors(
		podman.GenerateSystemdService(systemd, "", serverImage, false, "", []string{}),
		pgsql.GeneratePgsqlSystemdService(systemd, dbImage),
		systemd.ReloadDaemon(false),
	)
}
