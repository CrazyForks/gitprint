package builder

import (
	"testing"

	"github.com/google/go-github/v65/github"
)

func TestGenerateDocument(t *testing.T) {
	repo := &github.Repository{
		FullName:        github.String("testdata/testrepo"),
		StargazersCount: github.Int(10000),
		ForksCount:      github.Int(500),
	}

	tests := []struct {
		repository *github.Repository
		outputDir  string
		isNilErr   bool
	}{
		{repo, "notfound", false},
		{repo, "./testdata/testrepo", true},
	}

	for _, tt := range tests {
		t.Run(tt.outputDir, func(t *testing.T) {
			_, err := GenerateDocument(tt.repository, tt.outputDir)
			if tt.isNilErr && err != nil {
				t.Errorf("expecting nil error, got %v", err)
			}
		})
	}
}
