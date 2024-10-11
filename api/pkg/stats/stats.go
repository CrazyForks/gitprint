package stats

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/plutov/gitprint/api/pkg/log"
)

type State struct {
    mu sync.Mutex
}

func New() *State {
    return &State{}
}

func (s *State) SaveStats(statStr string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    path := filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), "stats.txt")


    file,err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
        if err != nil {
		  log.WithError(err).Error("could not open stats file")
          return
	  }

	  defer file.Close()

    if _, err := file.WriteString(statStr + "\n"); err != nil {
        log.WithError(err).Error("could not write stats to file")
    }
}
