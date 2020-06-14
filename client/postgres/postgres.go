package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"go-template/config"
	"go-template/util"
)

var db *gorm.DB

func init() {
	var err error
	cfg := config.GetConfig()

	db, err = gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.DBName, cfg.Postgres.Schema),
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected postgres db")

	db.DB().SetMaxIdleConns(100)

	if cfg.Debug {
		db = db.Debug()
	}
}

func GetClient(ctx context.Context) *gorm.DB {
	cloneDB := &gorm.DB{}
	*cloneDB = *db

	// use transaction per request
	if util.IsEnableTx(ctx) {
		tx := util.GetTx(ctx)
		return tx
	}

	return cloneDB
}

func Disconnect() {
	if db != nil {
		db.Close()
	}
}
