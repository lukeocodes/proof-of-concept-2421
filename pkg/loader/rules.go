package loader

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func LoadRules() ([]string, error) {
	var ruleDirectory string = ".cursor/rules"

	var rules []string
	err := filepath.Walk(ruleDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(ruleDirectory, path)
			if err != nil {
				return err
			}
			rules = append(rules, relativePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Info().Int("rules", len(rules)).Msg("Number of rules loaded")
	log.Trace().Str("rules", strings.Join(rules, ", ")).Msg("Rules loaded")

	return rules, nil
}
