package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/plutov/gitprint/api/pkg/auth"
	"github.com/plutov/gitprint/api/pkg/git"
	"github.com/plutov/gitprint/api/pkg/http/response"
	"github.com/plutov/gitprint/api/pkg/log"
)

func (h *Handler) githubURL(c echo.Context) error {
	url := h.Services.GithubAuth.GetRedirectURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) githubCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	token, err := h.Services.GithubAuth.GetAccessToken(code, state)
	if err != nil {
		return response.BadRequestDefaultMessage(c)
	}

	ghClient := git.NewClient(token)
	user, err := ghClient.GetCurrentUser()
	if err != nil {
		return response.InternalError(c, "unable to get current user")
	}

	jwtToken, err := auth.FillJWT(user)
	if err != nil {
		log.WithError(err).Error("failed to fill jwt token")
		return response.InternalError(c, "unable to create a session")
	}

	go h.Stats.SaveStats(fmt.Sprintf("user:%s,timestamp:%d", user.Username, time.Now().UTC().Unix()))

	return response.Ok(c, echo.Map{
		"jwt_token": jwtToken,
	})
}
