package main

import (
	"SolanaPoolScanner/db"
	"SolanaPoolScanner/types"
	"SolanaPoolScanner/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const PageLimit = 50 // 50 items per page

func main() {
	utils.Init()
	token := os.Getenv("API_TOKEN")
	tokenAddr := "26V8Ge3Y2sD1g8VspAsK2NddztywKch6jiTE9b1XNKGr"

	nCoLeader := 1
	nWorkers := 10

	for i := 0; i < nCoLeader; i++ {
		go coLeader(i, tokenAddr, token)
	}

	for j := 0; j < nWorkers; j++ {
		go worker(j, PageLimit*nWorkers, tokenAddr, j*PageLimit, token)
	}

	for {
		leader(tokenAddr, token)
	}
}

func leader(tokenAddr string, apiToken string) {
	for {
		fmt.Println("Leader")
		rsp, data, err := getTokenTrade(tokenAddr, 0, PageLimit, "swap", apiToken)
		if err == nil && rsp.Data.HasNext {
			db.Insert(rsp, data)
		}
	}
}

func coLeader(coLeaderID int, tokenAddr string, apiToken string) {
	for {
		fmt.Println("CoLeader ", coLeaderID)
		rsp, data, err := getTokenTrade(tokenAddr, PageLimit*(coLeaderID+1), PageLimit, "swap", apiToken)
		if err == nil && rsp.Data.HasNext {
			db.Insert(rsp, data)
		}
	}
}

func worker(workerID int, stepSize int, tokenAddr string, offset int, apiToken string) {
	for {
		fmt.Println("Worker ", workerID, " offset ", offset)

		var rsp types.Response
		var data string
		var err error
		for {
			rsp, data, err = getTokenTrade(tokenAddr, offset, PageLimit, "swap", apiToken)
			if err == nil && rsp.Success {
				break
			} else {
				fmt.Println("Worker ", workerID, " offset ", offset, " failed", "retrying")
				time.Sleep(3 * time.Second)
			}
		}

		if rsp.Data.HasNext {
			db.Insert(rsp, data)
			offset += stepSize
		} else {
			break
		}
	}
}

func getTokenTrade(addr string, offset int, lim int, txType string, token string) (types.Response, string, error) {
	url := fmt.Sprintf("https://public-api.birdeye.so/defi/txs/token?address=%s&offset=%d&limit=%d&tx_type=%s", addr, offset, lim, txType)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-chain", "solana")
	req.Header.Add("X-API-KEY", token)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode, res.Body)
		return types.Response{}, "", errors.New("Error: " + string(body))
	}

	jsonData := types.Response{}
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		fmt.Println(err)
	}
	return jsonData, string(body), nil
}
