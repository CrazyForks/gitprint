package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/plutov/gitprint/api/pkg/files"
	"github.com/plutov/gitprint/api/pkg/http/response"
)

func (h *Handler) downloadExportFile(c echo.Context) error {
	exportID := c.QueryParam("export_id")
	ext := c.QueryParam("ext")

	if exportID == "" || ext == "" {
		return response.BadRequest(c, "missing export_id or ext")
	}

	var path string
	switch ext {
	case "html":
		path = files.GetExportHTMLFile(exportID)
		return c.File(path)
	case "pdf":
		path = files.GetExportPDFFile(exportID)
		return c.Attachment(path, fmt.Sprintf("%s.pdf", exportID))
	default:
		return response.BadRequest(c, "invalid ext")
	}
}
