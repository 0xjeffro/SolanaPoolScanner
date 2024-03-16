package utils

import (
	"SolanaPoolScanner/db"
	"os"
)

func CheckEnv() {
	necessaryEnv := []string{
		"API_TOKEN", "DSN",
		"TICKER",
	}
	for _, env := range necessaryEnv {
		if os.Getenv(env) == "" {
			panic("Missing environment variable: " + env)
		}
	}
}

func Init() {
	CheckEnv()
	db.CreateTable()
}
