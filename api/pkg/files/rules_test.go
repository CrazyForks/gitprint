package files

import (
	"testing"
)

func TestIsAllowedFile(t *testing.T) {
	tests := []struct {
		filename string
		result   bool
	}{
		{"test.txt", true},
		{"test.go", true},
		{"README.MD", true},
		{"LICENSE", true},
		{"config/Dockerfile", true},
		{"test", false},
		{".test.go", false},
		{"hello.min.css", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := IsAllowedFile(tt.filename)
			if got != tt.result {
				t.Errorf("expecting %t, got %t", tt.result, got)
			}
		})
	}
}

func TestIsAllowedDir(t *testing.T) {
	tests := []struct {
		dir    string
		result bool
	}{
		{"pkg", true},
		{"node_modules", false},
		{"vendor", false},
		{"dist", false},
		{".github", false},
		{".GITIGNORE", false},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			got := IsAllowedDir(tt.dir)
			if got != tt.result {
				t.Errorf("expecting %t, got %t", tt.result, got)
			}
		})
	}
}
