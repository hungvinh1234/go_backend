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
	"github.com/yudai/pp"
	"golang.org/x/crypto/bcrypt"
)

//Viet noi dung ham va lay cac ham tuong tac ben db

type usecase struct {
	// cardRepo card.Repository
	accountRepo account.Repository
}

type SignInResponse struct {
	model.Account
	Token string `json:"token,omitempty"`
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
		//th 1 khong tim duoc user
		if errors.Is(err, gorm.ErrRecordNotFound) {

			_, err := e.accountRepo.GetByUserEmail(ctx, account.Email)
			if err != nil {
				//th 1 khong tim duoc ket qua mong muon
				if errors.Is(err, gorm.ErrRecordNotFound) {

					//Neu RefUser ma khac null thi moi kiem tra
					if account.RefUser != "" {

						_, err := e.accountRepo.GetByUserName(ctx, account.RefUser)
						if err != nil {
							//th 1 khong tim duoc ket qua mong muon
							if errors.Is(err, gorm.ErrRecordNotFound) {
								return nil, util.NewError(nil, http.StatusInternalServerError, 1011, "Reference User not exist !")
							}
						}

					}

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

				return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user by email")

			}

			//email da ton tai
			return nil, util.NewError(nil, http.StatusInternalServerError, 1011, "Email already taken !")

		}
		//th 2 la viet ham sai
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user by username")
	}

	//user da ton tai
	return nil, util.NewError(nil, http.StatusInternalServerError, 1011, "Username already taken !")

}

//account de hung nguyen cai row ben ham GetByUser

func (e *usecase) SignIn(ctx context.Context, account *model.Account) (*SignInResponse, error) {

	accountinDB, err := e.accountRepo.GetByUserName(ctx, account.Username)
	//err la khong tra ve ket qua mong muon
	if err != nil {
		//th 1 khong tim duoc ket qua mong muon
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//account da tao trong db
			return nil, util.NewError(err, http.StatusInternalServerError, 1010, "User not existed !")
		}
		//th 2 la viet ham sai
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "cannot get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(accountinDB.Password), []byte(account.Password))
	if err != nil {
		//password khong trung
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "Password not match !")
	}

	accountinDB.Password = ""

	claims := model.MyCustomClaims{
		Account: *accountinDB,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.GetConfig().MySigningKey))
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "can not create token")
	}

	response := SignInResponse{
		Account: *accountinDB,
		Token:   ss,
	}

	return &response, nil

}

func (e *usecase) EditAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	currentUser := ctx.Value("user").(*jwt.Token).Claims.(*model.MyCustomClaims).Account
	pp.Println(currentUser)
	//th 1 khong tim duoc ket qua mong muon

	if currentUser.ID != account.ID {

		return nil, util.NewError(nil, http.StatusInternalServerError, 1010, "You don't have permission")

	}

	if account.Password != "" {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed hash password")
		}

		account.Password = string(hashPass)
	}

	accountUpdated, err := e.accountRepo.UpdateUser(ctx, account)

	if err != nil {

		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "Update failed")
	}

	return accountUpdated, nil

}

func (e *usecase) UserDetail(ctx context.Context, account *model.Account) (*model.Account, error) {

	currentUser := ctx.Value("user").(*jwt.Token).Claims.(*model.MyCustomClaims).Account
	pp.Println(currentUser)
	//th 1 khong tim duoc ket qua mong muon

	if currentUser.ID != account.ID {

		return nil, util.NewError(nil, http.StatusInternalServerError, 1010, "You don't have permission")

	}

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

	currentUser := ctx.Value("user").(*jwt.Token).Claims.(*model.MyCustomClaims).Account
	pp.Println(currentUser)
	//th 1 khong tim duoc ket qua mong muon

	if !currentUser.IsAdmin {

		return nil, util.NewError(nil, http.StatusInternalServerError, 1010, "Only Admin can access this page !")

	}

	accountlist, err := e.accountRepo.GetUserList(ctx)
	//err la khong tra ve ket qua mong muon
	if err != nil {
		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "failed get user list")
	}

	//account da tao trong db
	return accountlist, nil
}

func (e *usecase) DeleteAccount(ctx context.Context, account *model.Account) (*model.Account, error) {

	currentUser := ctx.Value("user").(*jwt.Token).Claims.(*model.MyCustomClaims).Account
	pp.Println(currentUser)
	//th 1 khong tim duoc ket qua mong muon

	if !currentUser.IsAdmin {

		return nil, util.NewError(nil, http.StatusInternalServerError, 1010, "Only Admin can access this page !")

	}

	accountDeleted, err := e.accountRepo.DeleteUser(ctx, account)

	if err != nil {

		return nil, util.NewError(err, http.StatusInternalServerError, 1010, "Delete failed")
	}

	return accountDeleted, nil

}
