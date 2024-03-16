package main

import (
	"SolanaPoolScanner/service"
	"SolanaPoolScanner/utils"
	"SolanaPoolScanner/workers"
	"os"
	"strconv"
)

func main() {
	utils.Init()

	nLeaders := 1
	nWorkers := 3
	var err error
	if os.Getenv("N_LEADERS") != "" {
		nLeaders, err = strconv.Atoi(os.Getenv("N_LEADERS"))
		if err != nil {
			panic(err)
		}
	}
	if os.Getenv("N_WORKERS") != "" {
		nWorkers, err = strconv.Atoi(os.Getenv("N_WORKERS"))
		if err != nil {
			panic(err)
		}
	}
	workers.StartWorkers(nLeaders, nWorkers)
	service.StartGin()
}
