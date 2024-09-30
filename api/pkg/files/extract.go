package files

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/plutov/gitprint/api/pkg/log"
)

type ExtractAndFilterResult struct {
	Files     int64
	OutputDir string
}

func ExtractAndFilterFiles(path string) (*ExtractAndFilterResult, error) {
	logCtx := log.With("path", path)
	logCtx.Info("extracting and filtering files")

	r, err := os.Open(path)
	if err != nil {
		logCtx.WithError(err).Error("failed to open file")
		return nil, err
	}
	defer r.Close()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		logCtx.WithError(err).Error("failed to create gzip reader")
		return nil, err
	}
	defer gzr.Close()

	res := &ExtractAndFilterResult{
		OutputDir: strings.Replace(path, ".tar.gz", "", 1),
	}
	// remove output dir if exists
	os.RemoveAll(res.OutputDir)

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			logCtx.With("res", *res).Info("extracted and filtered files")
			return res, nil

		// return any other error
		case err != nil:
			logCtx.WithError(err).Error("failed to read tar header")
			return nil, err

		// if the header is nil, just skip it
		case header == nil:
			continue
		}

		// remove root folder name but keep the hierarchy
		// eg. plutov-formulosity-xx123/README.md -> README.md
		// // eg. plutov-formulosity-xx123/src/main.go -> src/main.go
		parts := strings.Split(header.Name, string(filepath.Separator))
		if len(parts) > 0 {
			header.Name = strings.Join(parts[1:], "/")
		}
		target := filepath.Join(res.OutputDir, header.Name)

		// check the file type
		switch header.Typeflag {
		case tar.TypeReg:
			// skip empty and big files
			if header.Size == 0 || header.Size > MaxFileSize {
				continue
			}
			headerDir := filepath.Dir(header.Name)
			if headerDir != "." && !IsAllowedDir(headerDir) {
				continue
			}

			if !IsAllowedFile(header.Name) {
				continue
			}

			targetDir := filepath.Dir(target)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				logCtx.WithError(err).Error("failed to create directory")
				return nil, err
			}

			f, err := os.Create(target)
			if err != nil {
				logCtx.WithError(err).Error("failed to create file")
				return nil, err
			}

			if _, err := io.Copy(f, tr); err != nil {
				logCtx.WithError(err).Error("failed to copy file contents")
				return nil, err
			}

			f.Close()
			res.Files++
		}
	}
}
