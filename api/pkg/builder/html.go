package builder

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"

	"github.com/plutov/gitprint/api/pkg/files"
	"github.com/plutov/gitprint/api/pkg/log"
)

func GenerateHTML(w io.Writer, doc *Document, exportID string) error {
	logCtx := log.With("exportID", exportID)
	logCtx.Info("generating html")

	t, err := template.ParseFiles("./templates/base.html")
	if err != nil {
		logCtx.WithError(err).Error("failed to parse template")
		return err
	}

	// markdown to html
	for i, node := range doc.Nodes {
		if node.Type == NodeTypeFile && strings.HasSuffix(strings.ToLower(node.Title), ".md") {
			doc.Nodes[i].ContentFile.ContentHTML = template.HTML(MarkdownToHTML(node.ContentFile.Content))
			doc.Nodes[i].ContentFile.Content = ""
			doc.Nodes[i].ContentFile.IsMarkdown = true
		}
	}

	if err := t.Execute(w, doc); err != nil {
		logCtx.WithError(err).Error("failed to execute template")
		return err
	}

	logCtx.Info("html generated")
	return nil
}

func MarkdownToHTML(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	unsafeHTML := string(markdown.Render(doc, renderer))

	// sanitize HTML
	return bluemonday.UGCPolicy().Sanitize(unsafeHTML)
}

func GenerateAndSaveHTMLFile(doc *Document, exportID string) (string, error) {
	logCtx := log.With("exportID", exportID)
	logCtx.Info("saving html file")

	output := files.GetExportHTMLFile(exportID)

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		logCtx.WithError(err).Error("failed to create output directory")
		return "", err
	}
	o, err := os.Create(output)
	if err != nil {
		logCtx.WithError(err).Error("failed to create output file")
		return "", err
	}
	defer o.Close()

	if err := GenerateHTML(o, doc, exportID); err != nil {
		logCtx.WithError(err).Error("failed to generate html")
		return "", err
	}

	logCtx.Info("html file saved")
	return output, nil
}
