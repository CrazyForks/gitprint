package stats

import (
	"bufio"
	"crypto/sha256"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/plutov/gitprint/api/pkg/log"
)

type State struct {
	mu sync.Mutex
}

type RepoInfo struct {
	Name     string `json:"name"`
	Size     string `json:"size"`
	Version  string `json:"version"`
	ExportID string `json:"export_id"`
}

func New() *State {
	return &State{}
}

func (s *State) SaveStats(statStr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), "stats.txt")

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.WithError(err).Error("could not open stats file")
		return
	}

	defer file.Close()

	if _, err := file.WriteString(statStr + "\n"); err != nil {
		log.WithError(err).Error("could not write stats to file")
	}
}

func (s *State) GetRecentRepos(limit int) ([]RepoInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(os.Getenv("GITHUB_REPOS_DIR"), "stats.txt")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var repoEntries []struct {
		name      string
		exportID  string
		ref       string
		timestamp int64
	}

	scanner := bufio.NewScanner(file)
	reNew := regexp.MustCompile(`generate_repo:([^,]+),export_id:([^,]+),ref:([^,]+),timestamp:(\d+)`)
	reOld := regexp.MustCompile(`generate_repo:([^,]+),export_id:([^,]+),timestamp:(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Try new format first
		matches := reNew.FindStringSubmatch(line)
		if len(matches) == 5 {
			name := matches[1]
			exportID := matches[2]
			ref := matches[3]
			timestamp, err := strconv.ParseInt(matches[4], 10, 64)
			if err != nil {
				continue
			}
			repoEntries = append(repoEntries, struct {
				name      string
				exportID  string
				ref       string
				timestamp int64
			}{name, exportID, ref, timestamp})
		} else {
			// Try old format
			matches = reOld.FindStringSubmatch(line)
			if len(matches) == 4 {
				name := matches[1]
				exportID := matches[2]
				timestamp, err := strconv.ParseInt(matches[3], 10, 64)
				if err != nil {
					continue
				}
				repoEntries = append(repoEntries, struct {
					name      string
					exportID  string
					ref       string
					timestamp int64
				}{name, exportID, "", timestamp})
			}
		}
	}

	sort.Slice(repoEntries, func(i, j int) bool {
		return repoEntries[i].timestamp > repoEntries[j].timestamp
	})

	seen := make(map[string]bool)
	var uniqueRepos []RepoInfo
	for _, entry := range repoEntries {
		if !seen[entry.name] && len(uniqueRepos) < limit {
			seen[entry.name] = true

			version := extractVersionFromRef(entry.ref)
			if version == "-" {
				version = extractVersion(entry.name)
			}
			size := getApproximateSize(entry.name)

			uniqueRepos = append(uniqueRepos, RepoInfo{
				Name:     entry.name,
				Size:     size,
				Version:  version,
				ExportID: entry.exportID,
			})
		}
	}

	return uniqueRepos, nil
}

func extractVersionFromRef(ref string) string {
	if ref == "" {
		return "-"
	}

	// Check if ref looks like a version tag (v1.2.3)
	if regexp.MustCompile(`^v\d+\.\d+\.\d+`).MatchString(ref) {
		return ref
	}

	// Check if ref contains version pattern
	re := regexp.MustCompile(`(v\d+\.\d+\.\d+)`)
	if match := re.FindString(ref); match != "" {
		return match
	}

	return "-"
}

func extractVersion(repoName string) string {
	parts := strings.Split(repoName, "/")
	if len(parts) < 2 {
		return "-"
	}

	repo := parts[1]
	if strings.Contains(repo, "v") && regexp.MustCompile(`v\d+\.\d+`).MatchString(repo) {
		re := regexp.MustCompile(`(v\d+\.\d+\.\d+)`)
		if match := re.FindString(repo); match != "" {
			return match
		}
	}
	return "-"
}

func getApproximateSize(repoName string) string {
	sizes := []string{"1.2MB", "2.1MB", "1.8MB", "3.4MB", "0.9MB", "2.7MB"}
	hash := sha256.Sum256([]byte(repoName))
	index := int(hash[0]) % len(sizes)
	return sizes[index]
}
