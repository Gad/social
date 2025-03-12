package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {

	if val, ok := os.LookupEnv(key); !ok {
		return fallback
	} else {
		return val
	}

}

func GetInt(key string, fallback int) int {

	if val, ok := os.LookupEnv(key); !ok {
		return fallback
	} else {
		val2int, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return val2int
	}
}

func GetDuration(durationString string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return fallback
	}
	return duration
}
