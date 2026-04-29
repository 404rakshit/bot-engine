package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RequireEnv(envName string) (string, error) {
	if value, exists := os.LookupEnv(envName); exists {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable %q is not set", envName)
}

// GetEnvOrDefault fetches the env var, or falls back to a provided default.
// It uses your robust type-switching logic to support almost any fallback type.
func GetEnvOrDefault(envName string, fallback any) string {
	if value, exists := os.LookupEnv(envName); exists {
		return value
	}

	switch v := fallback.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		return fmt.Sprint(v)
	default:
		// Panicking here is appropriate because passing an unsupported
		// type is a developer error that should be caught at startup.
		panic(fmt.Sprintf("unsupported fallback type %T for env %q", fallback, envName))
	}
}

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
