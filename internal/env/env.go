package env

import (
	"os"
	"strconv"
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
