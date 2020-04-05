package units

import "os"

// Env return env value with fallback value
func Env(item, fallback string) string {
	e := os.Getenv(item)
	if e == "" {
		return fallback
	}
	return e
}
