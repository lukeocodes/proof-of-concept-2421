package filetree

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type FileStatus struct {
	Path      string
	Processed bool
}

type FileTree struct {
	Root      string
	Tree      string
	Files     []FileStatus
	Ignores   []string
	RulesOnly bool
}

// NewFileTree creates a new FileTree instance
func NewFileTree(root string) (*FileTree, error) {
	ft := &FileTree{
		Root:      root,
		Files:     make([]FileStatus, 0),
		RulesOnly: false,
	}
	if err := ft.loadIgnoreFile(); err != nil {
		return nil, err
	}
	// Always ignore the rules directory unless specifically requested
	ft.Ignores = append(ft.Ignores, ".cursor/rules")
	if err := ft.buildTree(); err != nil {
		return nil, err
	}
	return ft, nil
}

// NewRulesFileTree creates a new FileTree instance for the rules directory only
func NewRulesFileTree(root string) (*FileTree, error) {
	ft := &FileTree{
		Root:      filepath.Join(root, ".cursor/rules"),
		Files:     make([]FileStatus, 0),
		RulesOnly: true,
	}
	if err := ft.loadIgnoreFile(); err != nil {
		return nil, err
	}
	if err := ft.buildTree(); err != nil {
		return nil, err
	}
	return ft, nil
}

// loadIgnoreFile reads .lukeignore file
func (ft *FileTree) loadIgnoreFile() error {
	ignoreFile := filepath.Join(ft.Root, ".lukeignore")
	// If we're in rules-only mode, look for the ignore file in the parent directory
	if ft.RulesOnly {
		ignoreFile = filepath.Join(filepath.Dir(filepath.Dir(ft.Root)), ".lukeignore")
	}

	file, err := os.Open(ignoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			// Convert pattern to use OS-specific path separator
			pattern = filepath.FromSlash(pattern)
			ft.Ignores = append(ft.Ignores, pattern)
		}
	}
	return scanner.Err()
}

// matchPattern checks if a path matches a pattern
func (ft *FileTree) matchPattern(pattern, path string) bool {
	// Handle directory-specific patterns
	if strings.HasSuffix(pattern, string(os.PathSeparator)) {
		fi, err := os.Stat(path)
		if err != nil || !fi.IsDir() {
			return false
		}
		pattern = pattern[:len(pattern)-1]
	}

	// Handle patterns starting with "/"
	if strings.HasPrefix(pattern, string(os.PathSeparator)) {
		// Match from root directory
		pattern = pattern[1:]
		path, _ = filepath.Rel(ft.Root, path)
		matched, _ := filepath.Match(pattern, path)
		return matched
	}

	// Handle patterns with "**"
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		path, _ = filepath.Rel(ft.Root, path)

		// Check if path starts with first part
		if len(parts[0]) > 0 && !strings.HasPrefix(path, parts[0]) {
			return false
		}

		// Check if path ends with last part
		if len(parts[1]) > 0 && !strings.HasSuffix(path, parts[1]) {
			return false
		}

		return true
	}

	// Handle simple patterns
	path, _ = filepath.Rel(ft.Root, path)
	matched, _ := filepath.Match(pattern, filepath.Base(path))
	if matched {
		return true
	}

	// Try matching against full relative path
	matched, _ = filepath.Match(pattern, path)
	return matched
}

// shouldIgnore checks if a path should be ignored
func (ft *FileTree) shouldIgnore(path string) bool {
	// If we're in rules-only mode, only process files in the rules directory
	if ft.RulesOnly {
		relPath, err := filepath.Rel(ft.Root, path)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return true
		}
		return false
	}

	// Check against ignore patterns
	for _, pattern := range ft.Ignores {
		if ft.matchPattern(pattern, path) {
			return true
		}
	}
	return false
}

// buildTree creates an ASCII representation of the directory structure
func (ft *FileTree) buildTree() error {
	var builder strings.Builder
	err := filepath.Walk(ft.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ft.shouldIgnore(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(ft.Root, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Add file to flat array if it's not a directory
		if !info.IsDir() {
			ft.Files = append(ft.Files, FileStatus{
				Path:      relPath,
				Processed: false,
			})
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		prefix := strings.Repeat("|   ", depth)
		if depth > 0 {
			prefix = prefix[:len(prefix)-4] + "+-- "
		}

		// Add a slash suffix for directories
		name := info.Name()
		if info.IsDir() {
			name += "/"
		}

		builder.WriteString(prefix + name + "\n")
		return nil
	})

	if err != nil {
		return err
	}

	ft.Tree = builder.String()
	return nil
}

// MarkFileProcessed marks a file as processed
func (ft *FileTree) MarkFileProcessed(path string) {
	for i := range ft.Files {
		if ft.Files[i].Path == path {
			ft.Files[i].Processed = true
			return
		}
	}
}

// GetUnprocessedFiles returns a list of files that haven't been processed
func (ft *FileTree) GetUnprocessedFiles() []string {
	var unprocessed []string
	for _, file := range ft.Files {
		if !file.Processed {
			unprocessed = append(unprocessed, file.Path)
		}
	}
	return unprocessed
}

// GetProcessedFiles returns a list of files that have been processed
func (ft *FileTree) GetProcessedFiles() []string {
	var processed []string
	for _, file := range ft.Files {
		if file.Processed {
			processed = append(processed, file.Path)
		}
	}
	return processed
}
