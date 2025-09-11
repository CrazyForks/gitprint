package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/plutov/gitprint/api/pkg/http/response"
)

func (h *Handler) getRecentRepos(c echo.Context) error {
	repos, err := h.Stats.GetRecentRepos(6)
	if err != nil {
		return response.InternalError(c, "unable to get recent repos")
	}

	if len(repos)%2 != 0 && len(repos) > 0 {
		repos = repos[:len(repos)-1]
	}

	return response.Ok(c, map[string]interface{}{
		"repos": repos,
	})
}
