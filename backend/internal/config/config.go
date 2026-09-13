package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
}

func init() {
	// Minimal .env parser to keep local development working without external dependencies
	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) > 0 && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					if os.Getenv(key) == "" {
						os.Setenv(key, val)
					}
				}
			}
		}
	}
}

func LoadConfig() Config {
	uri := os.Getenv("MONGODB_URI")
	db := os.Getenv("MONGODB_DATABASE")

	return Config{
		MongoURI:      uri,
		MongoDatabase: db,
	}
}
