package loader

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Load(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, relativePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// load ignore patterns from .bazignore
	ignorePatterns, err := loadIgnorePatterns(".bazignore")
	if err != nil {
		return nil, err
	}

	log.Info().Int("ignorePatterns", len(ignorePatterns)).Msg("Number of ignore patterns found")
	log.Trace().Str("ignorePatterns", strings.Join(ignorePatterns, ", ")).Msg("Ignore patterns found")

	// add cursor rules to ignore patterns
	ignorePatterns = append(ignorePatterns, ".cursor/rules")

	// filter files based on ignore patterns
	filteredFiles := []string{}
	for _, file := range files {
		ignore := false
		for _, pattern := range ignorePatterns {
			matched, err := filepath.Match(pattern, file)
			if err != nil {
				return nil, err
			}
			if matched || strings.HasPrefix(file, pattern) {
				ignore = true
				break
			}
		}
		if !ignore {
			filteredFiles = append(filteredFiles, file)
		}
	}

	log.Info().Int("filteredFiles", len(filteredFiles)).Msg("Number of files after filtering")
	log.Trace().Str("filteredFiles", strings.Join(filteredFiles, ", ")).Msg("Filtered files")

	return filteredFiles, nil
}

func loadIgnorePatterns(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	patterns := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}
