package files

import (
	"os"
	"path/filepath"
)

const (
	publicInternalDir = "gitprint_public_internal"
)

func GetExportDir(exportID string) string {
	return filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), exportID)
}

func GetExportHTMLFile(exportID string) string {
	return filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), publicInternalDir, exportID) + ".html"
}

func GetExportPDFFile(exportID string) string {
	return filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), publicInternalDir, exportID) + ".pdf"
}
