package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

var conn *gorm.DB

func InitDB() error {
	DSN := os.Getenv("DSN")
	var err error
	conn, err = gorm.Open(postgres.New(postgres.Config{
		DSN: DSN + " TimeZone=Asia/Shanghai",
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetConnection() *gorm.DB {
	return conn
}
