package builder

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateHTML(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		isNilErr bool
		contains []string
	}{
		{"simple", &Document{
			Title: "plutov/plutov",
			Nodes: []DocumentNode{
				DocumentNode{
					Type: NodeTypeMeta,
					ContentMeta: &ContentMeta{
						FullName: "plutov/plutov",
					},
				},
			},
		}, true, []string{"<h1>plutov/plutov</h1>"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := GenerateHTML(w, tt.doc, tt.name)
			if tt.isNilErr && err != nil {
				t.Errorf("expecting nil error, got %v", err)
			}

			if tt.isNilErr {
				for _, c := range tt.contains {
					if !strings.Contains(w.String(), c) {
						t.Errorf("expecting to contain %s", c)
					}
				}
			}
		})
	}
}
