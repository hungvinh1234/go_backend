package util

import (
	"context"

	"github.com/jinzhu/gorm"
)

var keyEnable string = "postgres_tx_enable"
var keyTx string = "postgres_tx"

func TxBegin(ctx context.Context, getClient func(ctx context.Context) *gorm.DB) context.Context {
	db := getClient(ctx)
	tx := db.Begin()
	ctx = SetTx(ctx, tx)

	ctx = context.WithValue(ctx, keyEnable, true)
	return ctx
}

func TxEnd(ctx context.Context, txFunc func() error) (context.Context, error) {
	tx := GetTx(ctx)
	var err error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit().Error // if Commit returns error update err with commit err
		}
	}()

	err = txFunc()
	ctx = context.WithValue(ctx, keyEnable, false)
	return ctx, err
}

func IsEnableTx(ctx context.Context) bool {
	tx_enable, ok := ctx.Value(keyEnable).(bool)
	return ok && tx_enable
}

func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(keyTx).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	ctx = context.WithValue(ctx, keyTx, tx)
	return ctx
}
