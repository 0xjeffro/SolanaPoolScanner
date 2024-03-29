package workers

import (
	"SolanaPoolScanner/db"
	"SolanaPoolScanner/types"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"math"
	"net/http"
	"os"
	"time"
)

// PageLimit the maximum value is dictated by the birdeye API,
// see: https://docs.birdeye.so/reference/get_defi-txs-token
const PageLimit = 50
const MaxWorkers = 2
const MAXSLEEP = 60 * 1000

var LeaderStatus types.Leader
var WorkerStatus types.Worker

var taskC chan int
var taskFinished chan bool
var noMoreTask chan bool

func StartWorkers() {
	token := os.Getenv("API_TOKEN")
	tokenAddr := os.Getenv("TOKEN_ADDR")

	// init channels
	taskC = make(chan int, MaxWorkers)
	taskFinished = make(chan bool, MaxWorkers)
	noMoreTask = make(chan bool, MaxWorkers)

	go taskGenerator()

	for i := 0; i < MaxWorkers; i++ {
		go worker(i, tokenAddr, token)
	}
}

func taskGenerator() {
	offset := 0
	for {
		if len(taskC) == 0 {
			select {
			case <-taskFinished:
				for i := 0; i < MaxWorkers; i++ {
					noMoreTask <- true
				}
				return
			default:
				zap.S().Info("+++++ Generating tasks +++++")
				for i := 0; i < MaxWorkers; i++ {
					taskC <- offset
					offset += PageLimit
				}
			}
		}
	}
}

func worker(workerID int, tokenAddr string, apiToken string) {
	defer func() {
		zap.S().Info("Worker ", workerID, " is exiting")
		if err := recover(); err != nil {
			zap.S().Error("Worker ", workerID, " failed with error: ", err)
		}
	}()
	for { // keep fetching tasks from the task pipeline
		var rsp types.Response
		var err error
		select {
		case <-noMoreTask:
			return
		case offset := <-taskC:
			fmt.Println("Worker", workerID, "received", offset)
			for { // retry on error
				rsp, _, err = getTokenTrade(tokenAddr, offset, PageLimit, "swap", apiToken)

				if err == nil && rsp.Success {
					db.Insert(rsp)
					fmt.Println("Worker", workerID, "offset", offset, "done")
					if rsp.Data.HasNext == false {
						taskFinished <- true
					}
					break
				} else {
					zap.S().Error("Worker ", workerID, " offset ", offset, " failed", "retrying")
					time.Sleep(500 * time.Millisecond)
				}
			}
		}
	}
}

func StartLeader() {
	go leader(0, os.Getenv("TOKEN_ADDR"), os.Getenv("API_TOKEN"))
}

func leader(leaderID int, tokenAddr string, apiToken string) {
	defer func() {
		LeaderStatus.Active = false
		if err := recover(); err != nil {
			LeaderStatus.ExitMsg = fmt.Sprintf("%s", err)
		}
	}()

	sleepTime := 1
	lastDataMap := make(map[string]bool)
	lastHit := 0
	lastAvgSleepTime := 1
	loopCount := 0

	LeaderStatus.ID = leaderID
	LeaderStatus.Active = true

	for {
		loopCount += 1
		zap.S().Info("Loop ", loopCount, " ...", " LastHits:", lastHit, " LastSleep:", sleepTime, "ms", " AvgSleep:", lastAvgSleepTime, "ms")
		rsp, _, err := getTokenTrade(tokenAddr, PageLimit*leaderID, PageLimit, "swap", apiToken)

		LeaderStatus.LastAPICallAt = time.Now()
		if err != nil {
			LeaderStatus.APICallCountFail += 1
		} else {
			LeaderStatus.APICallCountSuccess += 1
		}

		if err == nil {
			hit := 0
			newLastDataMap := make(map[string]bool)
			for _, item := range rsp.Data.Items {
				newLastDataMap[item.TxHash] = true
				if val, ok := lastDataMap[item.TxHash]; ok && val == true {
					hit += 1
				}
			}
			if hit < PageLimit/5 {
				sleepTime = 1
			} else if PageLimit/5 <= hit && hit < PageLimit/5*3 {
				sleepTime = int(math.Max(float64(sleepTime/2), 1))
			} else {
				sleepTime = int(math.Min(float64(sleepTime+1000), MAXSLEEP))
			}
			lastDataMap = newLastDataMap
			lastAvgSleepTime = int((float64(loopCount)-1.0)/float64(loopCount)*float64(lastAvgSleepTime) + float64(sleepTime)/float64(loopCount))
			db.Insert(rsp)
		} else {
			zap.S().Error("Leader ", leaderID, err)
		}
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
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
