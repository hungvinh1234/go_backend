package account

import (
	"context"

	"github.com/jinzhu/gorm"

	"go-template/model"
)

func NewPG(getDB func(ctx context.Context) *gorm.DB) Repository {
	return pgRepository{getDB}
}

type pgRepository struct {
	getDB func(ctx context.Context) *gorm.DB
}

func (p pgRepository) GetByUserName(ctx context.Context, username string) (*model.Account, error) {
	db := p.getDB(ctx)
	account := model.Account{}

	err := db.Where("username = ?", username).First(&account).Error

	return &account, err
}

func (p pgRepository) CreateUser(ctx context.Context, account *model.Account) (*model.Account, error) {
	db := p.getDB(ctx)

	err := db.Create(account).Error

	return account, err
}

func (p pgRepository) GetByUserId(ctx context.Context, id int64) (*model.Account, error) {
	db := p.getDB(ctx)
	account := model.Account{}

	err := db.Where("id = ?", id).First(&account).Error

	return &account, err
}

func (p pgRepository) UpdateUser(ctx context.Context, account *model.Account) (*model.Account, error) {
	db := p.getDB(ctx)

	// vi vo tinh truoc do dc khoi tao la con tro roi
	err := db.Model(account).Updates(account).Error

	db.Model(account).Updates(map[string]interface{}{"is_admin": account.IsAdmin})
	//account o day da la con tro san roi
	return account, err
}

func (p pgRepository) GetUserList(ctx context.Context) (*[]model.Account, error) {
	db := p.getDB(ctx)
	accounts := []model.Account{}

	// vi vo tinh truoc do dc khoi tao la con tro roi
	err := db.Find(&accounts).Error
	//account o day da la con tro san roi
	return &accounts, err
}

func (p pgRepository) DeleteUser(ctx context.Context, account *model.Account) (*model.Account, error) {
	db := p.getDB(ctx)

	// vi vo tinh truoc do dc khoi tao la con tro roi
	err := db.Delete(account).Error
	//account o day da la con tro san roi
	return account, err
}

func (p pgRepository) GetByUserEmail(ctx context.Context, email string) (*model.Account, error) {
	db := p.getDB(ctx)
	account := model.Account{}

	err := db.Where("email = ?", email).First(&account).Error

	return &account, err
}
