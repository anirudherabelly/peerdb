package peerdbenv

import (
	"os"
	"strconv"

	"golang.org/x/exp/constraints"
)

// GetEnv returns the value of the environment variable with the given name
// and a boolean indicating whether the environment variable exists.
func getEnv(name string) (string, bool) {
	val, exists := os.LookupEnv(name)
	return val, exists
}

// GetEnvInt returns the value of the environment variable with the given name
// or defaultValue if the environment variable is not set or is not a valid
// integer value.
func getEnvInt(name string, defaultValue int) int {
	val, ok := getEnv(name)
	if !ok {
		return defaultValue
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return i
}

// getEnvUint32 returns the value of the environment variable with the given name
// or defaultValue if the environment variable is not set or is not a valid
// uint32 value.
func getEnvUint[T constraints.Unsigned](name string, defaultValue T) T {
	val, ok := getEnv(name)
	if !ok {
		return defaultValue
	}

	// widest bit size, truncate later
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return T(i)
}

// getEnvBool returns the value of the environment variable with the given name
// or defaultValue if the environment variable is not set or is not a valid
// boolean value.
func getEnvBool(name string, defaultValue bool) bool {
	val, ok := getEnv(name)
	if !ok {
		return defaultValue
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return b
}

// GetEnvString returns the value of the environment variable with the given name
// or defaultValue if the environment variable is not set.
func getEnvString(name string, defaultValue string) string {
	val, ok := getEnv(name)
	if !ok {
		return defaultValue
	}

	return val
}
