package env

import (
	"bufio"
	"os"
	"strings"
)

// Load reads the .env file and sets environment variables
func Load(filename string) error {
	// Open .env file
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If .env doesn't exist, try sample.env
			file, err = os.Open("sample.env")
			if err != nil {
				return nil // Neither file exists, return without error
			}
		} else {
			return err
		}
	}
	defer file.Close()

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Only set if not already set in environment
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
