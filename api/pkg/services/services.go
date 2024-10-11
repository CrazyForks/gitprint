package services

import (
	"github.com/plutov/gitprint/api/pkg/git"
	"github.com/plutov/gitprint/api/pkg/stats"
)

type Services struct {
	GithubAuth          *git.Auth
	GenerateRateLimiter *git.TTLMap
	Stats               *stats.State
}

func InitServices() (Services, error) {
	svc := Services{
		GithubAuth:          git.NewAuth(),
		GenerateRateLimiter: git.NewTTLMap(1000, 3600),
		Stats:               stats.New(),
	}

	return svc, nil
}
