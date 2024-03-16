package db

import (
	"SolanaPoolScanner/types"
	"fmt"
)

func Insert(rep types.Response, jsonData string) {
	rows := make([]Trade, 0)
	for _, item := range rep.Data.Items {
		txHash := item.TxHash
		owner := item.Owner
		estPrice := 0.0
		uiChangeAmount := 0.0
		blockUnixTime := item.BlockUnixTime

		if item.Side == "sell" {
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
		} else {
			if item.Quote.Price != nil {
				estPrice = *item.QuotePrice
			} else if item.Quote.NearestPrice != nil {
				estPrice = *item.Quote.NearestPrice
			} else if item.Base.Price != nil || item.Base.NearestPrice != nil {
				if item.Base.Price != nil {
					estPrice = *item.Base.Price * item.Base.UiAmount / item.Quote.UiAmount
				} else {
					estPrice = *item.Base.NearestPrice * item.Base.UiAmount / item.Quote.UiAmount
				}
			}
			uiChangeAmount = item.Quote.UiChangeAmount
		}

		row := Trade{
			TxHash:         txHash,
			Owner:          owner,
			EstPrice:       estPrice,
			UiChangeAmount: uiChangeAmount,
			BlockUnixTime:  blockUnixTime,
			Data:           jsonData,
		}
		rows = append(rows, row)
	}

	// Insert rows into database
	conn := GetConnection()
	err := conn.Create(&rows).Error
	if len(rows) == 0 {
		fmt.Println("No rows to insert", jsonData)
	}
	if err != nil {
		fmt.Println("Error inserting rows: ", err)
	}
}
