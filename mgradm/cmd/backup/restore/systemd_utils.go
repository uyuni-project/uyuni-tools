// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restore

import (
	"archive/tar"
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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
				hasError = errors.Join(hasError, err)
				continue
			}
		case tar.TypeReg:
			fh, err := os.Create(header.Name)
			if err != nil {
				log.Warn().Err(err).Msgf(L("Unable to create %s"), header.Name)
				hasError = errors.Join(hasError, err)
				continue
			}
			if _, err := io.Copy(fh, tr); err != nil {
				log.Warn().Err(err).Msgf(L("Unable to restore content of %s"), header.Name)
				hasError = errors.Join(hasError, err)
				fh.Close()
				os.Remove(header.Name)
				continue
			}
			fh.Close()
		default:
			log.Warn().Msgf(L("Unknown filetype of %s"), header.Name)
			continue
		}

		if err := restoreFileAttributes(header.Name, header); err != nil {
			log.Warn().Err(err).Msgf(L("Unable to restore file details for %s"), header.Name)
			hasError = errors.Join(hasError, err)
		}
	}
	return hasError
}

func restoreFileAttributes(filename string, th *tar.Header) error {
	var e error
	e = errors.Join(e, os.Chmod(filename, th.FileInfo().Mode()))
	e = errors.Join(e, os.Chown(filename, th.Uid, th.Gid))
	e = errors.Join(e, os.Chtimes(filename, th.AccessTime, th.ModTime))
	return e
}
