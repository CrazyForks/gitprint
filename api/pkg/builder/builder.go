package builder

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/google/go-github/v65/github"
	"github.com/plutov/gitprint/api/pkg/log"
)

const (
	NodeTypeMeta = "meta"
)

type DocumentNode struct {
}

type Document struct {
	Nodes []DocumentNode
}

func GenerateDocument(repo *github.Repository, outputDir string) (*Document, error) {
	logCtx := log.With("repo", repo.GetFullName(), "outputDir", outputDir)
	logCtx.Info("generating document")

	err := filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(path, d.Name(), "directory?", d.IsDir())
		return nil
	})
	if err != nil {
		logCtx.WithError(err).Error("unable to walk directories")
		return nil, err
	}

	return nil, nil
}
