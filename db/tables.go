package db

import "os"

type Trade struct {
	TxHash         string  `json:"txHash" gorm:"primaryKey;column:txHash;"`
	Owner          string  `json:"owner" gorm:"column:owner;"`
	EstPrice       float64 `json:"estPrice" gorm:"column:estPrice;"`
	UiChangeAmount float64 `json:"uiChangeAmount" gorm:"column:uiChangeAmount;"`
	BlockUnixTime  int     `json:"blockUnixTime" gorm:"column:blockUnixTime;"`
	Data           string  `json:"data" gorm:"column:data;"`
}

func (Trade) TableName() string {
	return os.Getenv("TICKER")
}

func CreateTable() {
	err := InitDB()
	if err != nil {
		panic(err)
	}
	db := GetConnection()

	if !db.Migrator().HasTable(&Trade{}) {
		err = db.Migrator().CreateTable(&Trade{})
		if err != nil {
			panic(err)
		}
	} else {
		panic("Table already exists")
	}
}
