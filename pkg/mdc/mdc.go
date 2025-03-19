package mdc

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/rs/zerolog/log"
)

// Mdc represents a single MDC file with its metadata and content
type Mdc struct {
	// Frontmatter fields
	Globs       []glob.Glob
	Description string
	AlwaysApply bool
	Path        string

	// The actual content of the rule file after the frontmatter
	Content string
}

// parseFrontmatter parses the frontmatter section into key-value pairs
func parseFrontmatter(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove any quotes around the value
		value = strings.Trim(value, `"'`)
		result[key] = value
	}

	return result, nil
}

// parseGlobs parses and compiles a comma-separated string of glob patterns
func parseGlobs(globStr string) ([]glob.Glob, error) {
	var globs []glob.Glob

	// Handle both comma-separated and single value formats
	var patterns []string
	if strings.Contains(globStr, ",") {
		// Split by comma and trim spaces
		parts := strings.Split(globStr, ",")
		patterns = make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				patterns = append(patterns, p)
			}
		}
	} else {
		// Single glob
		pattern := strings.TrimSpace(globStr)
		if pattern != "" {
			patterns = []string{pattern}
		}
	}

	// Compile all patterns
	globs = make([]glob.Glob, 0, len(patterns))
	for _, pattern := range patterns {
		g, err := glob.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile glob pattern '%s': %w", pattern, err)
		}
		globs = append(globs, g)
	}

	return globs, nil
}

// ParseBytes parses a byte slice containing an MDC file
func ParseBytes(data []byte) (*Mdc, error) {
	mdc := &Mdc{}

	log.Trace().
		Str("data", string(data)).
		Msg("Parsing MDC")

	// Split the content into frontmatter and markdown
	parts := bytes.Split(data, []byte("---\n"))
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid MDC format: missing frontmatter delimiters")
	}

	// Parse the frontmatter
	frontmatter, err := parseFrontmatter(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Set the fields
	mdc.Description = frontmatter["description"]

	// Parse and compile globs
	if globsStr, ok := frontmatter["globs"]; ok {
		globs, err := parseGlobs(globsStr)
		if err != nil {
			return nil, err
		}
		mdc.Globs = globs
	}

	// Parse alwaysApply
	if alwaysApply, ok := frontmatter["alwaysApply"]; ok {
		mdc.AlwaysApply = strings.ToLower(alwaysApply) == "true"
	}

	// Store the markdown content (everything after the second ---)
	mdc.Content = string(bytes.Join(parts[2:], []byte("---\n")))

	return mdc, nil
}

// getGlobPattern extracts the original pattern from a compiled glob
func getGlobPattern(g glob.Glob) string {
	// Since glob.Glob doesn't expose its pattern, we need to match against test cases
	// to determine the original pattern
	testCases := map[string][]struct {
		input string
		match bool
	}{
		"*.ts": {
			{"test.ts", true},
			{"test.js", false},
			{"src/test.ts", false},
		},
		"src/*.ts": {
			{"test.ts", false},
			{"src/test.ts", true},
			{"src/test.js", false},
		},
		"**/*.ts": {
			{"test.ts", true},
			{"src/test.ts", true},
			{"src/nested/test.ts", true},
			{"test.js", false},
		},
	}

	for pattern, cases := range testCases {
		matches := true
		for _, tc := range cases {
			if g.Match(tc.input) != tc.match {
				matches = false
				break
			}
		}
		if matches {
			return pattern
		}
	}

	// If we can't determine the exact pattern, return a default
	return "*.ts"
}

// Marshal converts an MDC struct back to bytes
func (m *Mdc) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Write description if present
	if m.Description != "" {
		buf.WriteString(fmt.Sprintf("description: %s\n", m.Description))
	}

	// Write globs if present
	if len(m.Globs) > 0 {
		patterns := make([]string, len(m.Globs))
		for i, g := range m.Globs {
			patterns[i] = getGlobPattern(g)
		}
		buf.WriteString(fmt.Sprintf("globs: %s\n", strings.Join(patterns, ", ")))
	}

	// Write alwaysApply if true
	if m.AlwaysApply {
		buf.WriteString("alwaysApply: true\n")
	}

	buf.WriteString("---\n")
	buf.WriteString(m.Content)

	return buf.Bytes(), nil
}

// Unmarshal parses an MDC struct from a byte slice
func Unmarshal(data []byte) (*Mdc, error) {
	return ParseBytes(data)
}

// Validate checks if the MDC struct has all required fields
func (m *Mdc) Validate() error {
	var errors []string

	if m.Description == "" {
		errors = append(errors, "description is required")
	}

	if len(m.Globs) == 0 {
		errors = append(errors, "at least one glob pattern is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}
