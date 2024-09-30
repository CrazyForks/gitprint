package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/plutov/gitprint/api/pkg/files"
	"github.com/plutov/gitprint/api/pkg/git"
	"github.com/plutov/gitprint/api/pkg/http/response"
)

func (h *Handler) downloadRepo(c echo.Context) error {
	token := c.QueryParam("token")
	owner := c.QueryParam("owner")
	repo := c.QueryParam("repo")
	ref := c.QueryParam("ref")

	ghClient := git.NewClient(token)

	res, err := ghClient.DownloadRepo(owner, repo, ref)
	if err != nil {
		return response.InternalError(c, "unable to download repo")
	}

	extracted, err := files.ExtractAndFilterFiles(res.OutputFile)
	if err != nil {
		return response.InternalError(c, "unable to extract and filter files")
	}

	return response.Ok(c, extracted)
}
