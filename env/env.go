package env

import "os"

func Get(name string) string {
	return os.Getenv(name)
}

func GetEnvironmentName() string {
	if env := os.Getenv("APP_ENV"); len(env) > 0 {
		return env
	}

	return "local"
}
