package main

import (
	"SolanaPoolScanner/service"
	"SolanaPoolScanner/utils"
	"SolanaPoolScanner/workers"
	"os"
)

var MODE = os.Getenv("MODE")

func main() {
	utils.Init()
	if MODE == "" || MODE == "0" {
		workers.StartLeader()
	} else {
		workers.StartWorkers()
	}
	service.StartGin()
}
