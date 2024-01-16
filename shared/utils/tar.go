// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Extracts a tar.gz file.
func ExtractTarGz(tarballPath string, dstPath string) error {
	reader, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path, err := filepath.Abs(filepath.Join(dstPath, header.Name))
		if err != nil {
			return err
		}
		if !strings.HasPrefix(path, dstPath) {
			log.Warn().Msgf("Skipping extraction of %s in %s file as is resolves outside the target path",
				header.Name, tarballPath)
			continue
		}

		info := header.FileInfo()
		if info.IsDir() {
			log.Debug().Msgf("Creating folder %s", path)
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		log.Debug().Msgf("Extracting file %s", path)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err = io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}
