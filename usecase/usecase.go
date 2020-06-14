package usecase

import (
	"go-template/repository"
	"go-template/usecase/account"
)

// Usecase .
type Usecase struct {
	Account account.IUsecase
}

// New .
func New(repo *repository.Repository) *Usecase {
	return &Usecase{
		Account: account.New(repo.Account),
	}
}
