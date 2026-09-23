package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL        string
	SupabaseKey        string
	SupabaseProjectRef string
	Port               string
}

func loadConfig() (Config, error) {
	_ = godotenv.Load()

	config := Config{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Port:        os.Getenv("PORT"),
	}

	if config.SupabaseURL == "" {
		return Config{}, fmt.Errorf("SUPABASE_URL is required")
	}

	if config.SupabaseKey == "" {
		return Config{}, fmt.Errorf("SUPABASE_KEY is required")
	}

	parsedURL, err := url.Parse(config.SupabaseURL)
	if err != nil {
		return Config{}, fmt.Errorf("invalid SUPABASE_URL: %w", err)
	}

	host := parsedURL.Hostname()
	const suffix = ".supabase.co"

	if !strings.HasSuffix(host, suffix) {
		return Config{}, fmt.Errorf("SUPABASE_URL must be a Supabase project URL")
	}

	config.SupabaseProjectRef = strings.TrimSuffix(host, suffix)

	if config.SupabaseProjectRef == "" {
		return Config{}, fmt.Errorf("could not determine Supabase project reference")
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	return config, nil
}
