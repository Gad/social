package env

import (
	"os"
	"strconv"
	"time"
)
// GetString retrieves a string value from the environment variables.
// If the variable is not set, it returns the provided default value.
func GetString(key, fallback string) string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

// GetInt retrieves an integer value from the environment variables.
// If the variable is not set, it returns the provided default value.
func GetInt(key string, fallback int) int {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	val2int, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return val2int
}

// GetDurantion retrieves an time.Duration value from the environment variables.
// If the variable is not set, it returns the provided default value.
func GetDuration(durationString string, fallback time.Duration) time.Duration {

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return fallback
	}
	return duration
}

// GetBool retrieves a boolean value from the environment variables.
// If the variable is not set, it returns the provided default value.
func GetBool(key string, fallback bool) bool {

	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return parsed
}
