package workers

import (
	"SolanaPoolScanner/db"
	"SolanaPoolScanner/types"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

const PageLimit = 50 // 50 items per page

var WorkerStatus types.Status

func initStats(nLeaders int, nWorkers int) {
	WorkerStatus = types.Status{
		Leaders: make([]types.Leader, nLeaders),
		Workers: make([]types.Worker, nWorkers),
		Clock:   types.Clock{},
	}
}

func StartWorkers(nLeaders int, nWorkers int) {
	initStats(nLeaders, nWorkers)
	token := os.Getenv("API_TOKEN")
	tokenAddr := os.Getenv("TOKEN_ADDR")

	for i := 0; i < nLeaders; i++ {
		go leader(i, tokenAddr, token)
	}

	for i := 0; i < nWorkers; i++ {
		go worker(i, PageLimit*nWorkers, tokenAddr, i*PageLimit, token)
	}
	go clock()
}

func clock() {
	WorkerStatus.Clock.Active = true
	defer func() {
		WorkerStatus.Clock.Active = false
		if err := recover(); err != nil {
			WorkerStatus.Clock.ExitMsg = fmt.Sprintf("%s", err)
		}
	}()
	for {
		allCatchedUp := true
		activeWorkers := 0
		for i := range WorkerStatus.Workers {
			if WorkerStatus.Workers[i].Active {
				activeWorkers += 1
				if WorkerStatus.Workers[i].APICallCountSuccess < WorkerStatus.Clock.Clock {
					allCatchedUp = false
				}
			}
		}
		if allCatchedUp {
			WorkerStatus.Clock.Clock += 1
		}
		if activeWorkers == 0 {
			break
		}
	}
}

func leader(leaderID int, tokenAddr string, apiToken string) {
	defer func() {
		WorkerStatus.Leaders[leaderID].Active = false
		if err := recover(); err != nil {
			WorkerStatus.Leaders[leaderID].ExitMsg = fmt.Sprintf("%s", err)
		}
	}()

	WorkerStatus.Leaders[leaderID].ID = leaderID
	WorkerStatus.Leaders[leaderID].Active = true

	for {
		zap.S().Info("Leader ", leaderID)
		rsp, data, err := getTokenTrade(tokenAddr, PageLimit*leaderID, PageLimit, "swap", apiToken)

		WorkerStatus.Leaders[leaderID].LastAPICallAt = time.Now()
		if err != nil {
			WorkerStatus.Leaders[leaderID].APICallCountFail += 1
		} else {
			WorkerStatus.Leaders[leaderID].APICallCountSuccess += 1
		}

		if err == nil && rsp.Data.HasNext {
			db.Insert(rsp, data)
		}
	}
}

func worker(workerID int, stepSize int, tokenAddr string, offset int, apiToken string) {
	defer func() {
		WorkerStatus.Workers[workerID].Active = false
		if err := recover(); err != nil {
			WorkerStatus.Workers[workerID].ExitMsg = fmt.Sprintf("%s", err)
		}
	}()

	WorkerStatus.Workers[workerID].ID = workerID
	WorkerStatus.Workers[workerID].Active = true

	for {
		zap.S().Info("Worker ", workerID, " offset ", offset)

		var rsp types.Response
		var data string
		var err error
		for {
			if WorkerStatus.Workers[workerID].APICallCountSuccess < WorkerStatus.Clock.Clock {
				break
			}
		}
		for {
			rsp, data, err = getTokenTrade(tokenAddr, offset, PageLimit, "swap", apiToken)

			WorkerStatus.Workers[workerID].LastAPICallAt = time.Now()
			if err != nil {
				WorkerStatus.Workers[workerID].APICallCountFail += 1
			} else {
				WorkerStatus.Workers[workerID].APICallCountSuccess += 1
			}

			if err == nil && rsp.Success {
				break
			} else {
				zap.S().Error("Worker ", workerID, " offset ", offset, " failed", "retrying")
				time.Sleep(1 * time.Second)
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
		zap.S().Error("Error: ", res.StatusCode, res.Body)
		return types.Response{}, "", errors.New("Error: " + string(body))
	}

	jsonData := types.Response{}
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		fmt.Println(err)
	}
	return jsonData, string(body), nil
}
