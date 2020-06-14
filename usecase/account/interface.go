package account

import (
	"context"
	"go-template/model"
)

// IUsecase . Xu ly logic
// Khai bao cac ham se dung ben usecase
type IUsecase interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	SignIn(ctx context.Context, account *model.Account) (*SignInResponse, error)
	EditAccount(ctx context.Context, account *model.Account) (*model.Account, error)
}
