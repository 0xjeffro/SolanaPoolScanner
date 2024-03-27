package db

import (
	"SolanaPoolScanner/types"
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

func Insert(rep types.Response) {
	rows := make([]Trade, 0)
	for _, item := range rep.Data.Items {
		txHash := item.TxHash
		owner := item.Owner
		estPrice := 0.0
		uiChangeAmount := 0.0
		blockUnixTime := item.BlockUnixTime

		jsonItem, _ := json.Marshal(item)

		if item.Base.Price != nil {
			estPrice = *item.BasePrice
		} else if item.Base.NearestPrice != nil {
			estPrice = *item.Base.NearestPrice
		} else if item.Quote.Price != nil || item.Quote.NearestPrice != nil {
			if item.Quote.Price != nil {
				estPrice = *item.Quote.Price * item.Quote.UiAmount / item.Base.UiAmount
			} else {
				estPrice = *item.Quote.NearestPrice * item.Quote.UiAmount / item.Base.UiAmount
			}
		}
		uiChangeAmount = item.Base.UiChangeAmount

		row := Trade{
			TxHash:         txHash,
			Owner:          owner,
			EstPrice:       estPrice,
			UiChangeAmount: uiChangeAmount,
			BlockUnixTime:  blockUnixTime,
			Data:           string(jsonItem),
		}
		rows = append(rows, row)
	}

	// Insert rows into database
	conn := GetConnection()
	if len(rows) != 0 {
		err := conn.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
		//err := conn.Create(&rows).Error
		if err != nil {
			zap.S().Error("Error inserting rows: ", err)
		}
	}
}
