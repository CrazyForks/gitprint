package builder

import (
	"bytes"
	"testing"
)

func TestGenerateHTML(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		isNilErr bool
		content  string
	}{
		{"empty", &Document{}, true, `<!doctype html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <title>TODO</title>
    </head>
    <body>
        HELLO
    </body>
</html>
`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := GenerateHTML(w, tt.doc, tt.name)
			if tt.isNilErr && err != nil {
				t.Errorf("expecting nil error, got %v", err)
			}

			if tt.isNilErr && w.String() != tt.content {
				t.Errorf("expecting %s, got %s", tt.content, w.String())
			}
		})
	}
}
