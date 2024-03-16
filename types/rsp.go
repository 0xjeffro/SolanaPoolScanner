package types

import (
	"encoding/json"
)

type Quote struct {
	Symbol         string      `json:"symbol"`
	Decimals       int         `json:"decimals"`
	Address        string      `json:"address"`
	Amount         json.Number `json:"amount"`
	FeeInfo        *string     `json:"feeInfo"`
	UiAmount       float64     `json:"uiAmount"`
	Price          *float64    `json:"price"`
	NearestPrice   *float64    `json:"nearestPrice"`
	ChangeAmount   json.Number `json:"changeAmount"`
	UiChangeAmount float64     `json:"uiChangeAmount"`
}

type Base struct {
	Symbol         string      `json:"symbol"`
	Decimals       int         `json:"decimals"`
	Address        string      `json:"address"`
	Amount         json.Number `json:"amount"`
	UiAmount       float64     `json:"uiAmount"`
	Price          *float64    `json:"price"`
	NearestPrice   *float64    `json:"nearestPrice"`
	ChangeAmount   json.Number `json:"changeAmount"`
	UiChangeAmount float64     `json:"uiChangeAmount"`
}

type Item struct {
	Quote         Quote    `json:"quote"`
	Base          Base     `json:"base"`
	BasePrice     *float64 `json:"basePrice"`
	QuotePrice    *float64 `json:"quotePrice"`
	TxHash        string   `json:"txHash"`
	Source        string   `json:"source"`
	BlockUnixTime int      `json:"blockUnixTime"`
	TxType        string   `json:"txType"`
	Owner         string   `json:"owner"`
	Side          string   `json:"side"`
	Alias         *string  `json:"alias"`
	PricePair     float64  `json:"pricePair"`
	From          Base     `json:"from"`
	To            Quote    `json:"to"`
	TokenPrice    float64  `json:"tokenPrice"`
	PoolId        string   `json:"poolId"`
}

type Data struct {
	Items   []Item `json:"items"`
	HasNext bool   `json:"hasNext"`
}

type Response struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}
