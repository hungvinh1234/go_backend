package account

import (
	"context"
	"errors"
	"log"
	"net/http"

	"go-template/config"
	"go-template/model"
	"go-template/repository/account"
	"go-template/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Viet noi dung ham va lay cac ham tuong tac ben db

type usecase struct {
	// cardRepo card.Repository
	accountRepo account.Repository
}

type SignInResponse struct {
	model.Account
	Token string
}

func New(accountRepo account.Repository) IUsecase {
	return &usecase{accountRepo}
}

//truyen vao 1 account chua gia tri lay tu form
func (e *usecase) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {

	//account de hung nguyen cai row ben ham GetByUser
	_, err := e.accountRepo.GetByUserName(ctx, account.Username)

	//err la khong tra ve ket qua mong muon
	if err != nil {
		//th 1 khong tim duoc ket qua mong muon
		if errors.Is(err, gorm.ErrRecordNotFound) {

			//ham lay cua nguoi ta nen ko func

			hashPass, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed hash password")
			}

			account.Password = string(hashPass)

			log.Println(account.Birthday)
			//dua account vua tao vao db
			account, err := e.accountRepo.CreateUser(ctx, account)
			if err != nil {
				return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed create user to db")
			}

			//account da tao trong db
			return account, nil
		}
		//th 2 la viet ham sai
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user to db")
	}

	//user da ton tai
	return nil, util.NewError(nil, http.StatusInternalServerError, 1010, "user existed in db")

}

//account de hung nguyen cai row ben ham GetByUser

func (e *usecase) SignIn(ctx context.Context, account *model.Account) (*SignInResponse, error) {

	accountinDB, err := e.accountRepo.GetByUserName(ctx, account.Username)
	//err la khong tra ve ket qua mong muon
	if err != nil {
		//th 1 khong tim duoc ket qua mong muon
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//account da tao trong db
			return nil, util.NewError(err, http.StatusInternalServerError, 1010, "user not existed in db")
		}
		//th 2 la viet ham sai
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "cannot get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(accountinDB.Password), []byte(account.Password))
	if err != nil {
		//password khong trung
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "password not match")
	}

	claims := model.MyCustomClaims{
		Account: *accountinDB,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.GetConfig().MySigningKey))
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "can not create token")
	}

	accountinDB.Password = ""

	response := SignInResponse{
		Account: *accountinDB,
		Token:   ss,
	}

	return &response, nil

}

func (e *usecase) EditAccount(ctx context.Context, account *model.Account) (*model.Account, error) {

	//th 1 khong tim duoc ket qua mong muon

	hashPass, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed hash password")
	}

	account.Password = string(hashPass)

	accountUpdated, err := e.accountRepo.UpdateUser(ctx, account)

	if err != nil {

		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "Update failed")
	}

	return accountUpdated, nil

}

func (e *usecase) UserDetail(ctx context.Context, account *model.Account) (*model.Account, error) {

	accountdetail, err := e.accountRepo.GetByUserId(ctx, account.ID)
	accountdetail.Password = ""
	//err la khong tra ve ket qua mong muon
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user by id")
	}

	//account da tao trong db
	return accountdetail, nil
}

func (e *usecase) ShowUserList(ctx context.Context) (*[]model.Account, error) {

	accountlist, err := e.accountRepo.GetUserList(ctx)
	//err la khong tra ve ket qua mong muon
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user list")
	}

	//account da tao trong db
	return accountlist, nil
}

func (e *usecase) DeleteAccount(ctx context.Context, account *model.Account) (*model.Account, error) {

	accountDeleted, err := e.accountRepo.DeleteUser(ctx, account)

	if err != nil {

		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "Delete failed")
	}

	return accountDeleted, nil

}
