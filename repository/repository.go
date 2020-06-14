package repository

import (
	"context"

	"github.com/jinzhu/gorm"

	"go-template/repository/account"
)

// Repository .
type Repository struct {
	Account account.Repository
}

// New .
func New(getClient func(ctx context.Context) *gorm.DB) *Repository {
	return &Repository{
		Account: account.NewPG(getClient),
	}
}
