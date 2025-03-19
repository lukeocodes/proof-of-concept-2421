package rules

import (
	"concept/pkg/mdc"

	"os"
	"path/filepath"
)

// Rules represents a collection of rules with methods to match files
type Rules struct {
	rules []mdc.Mdc
}

// New creates a new Rules instance from a slice of file paths
func New(filePaths []string) (*Rules, error) {
	rules := make([]mdc.Mdc, 0, len(filePaths))

	for _, path := range filePaths {
		rule, err := parseRuleFile(filepath.Join(".cursor/rules", path))
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return &Rules{rules: rules}, nil
}

// GetMatchingRules returns all rules that match the given file path
func (r *Rules) GetMatchingRules(filePath string) []mdc.Mdc {
	matching := make([]mdc.Mdc, 0)

	for _, rule := range r.rules {
		for _, ruleGlob := range rule.Globs {
			if ruleGlob.Match(filePath) {
				matching = append(matching, rule)
			}
		}
	}

	return matching
}

// parseRuleFile reads a markdown file with TOML frontmatter and returns a Rule
func parseRuleFile(filePath string) (mdc.Mdc, error) {
	var rule mdc.Mdc

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return rule, err
	}

	// Parse the frontmatter (MDC)
	mdc, err := mdc.Unmarshal(content)
	if err != nil {
		return rule, err
	}

	mdc.Path = filePath

	return *mdc, nil
}

// GetAllRules returns all rules
func (r *Rules) GetAllRules() []mdc.Mdc {
	return r.rules
}
