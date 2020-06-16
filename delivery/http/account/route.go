package account

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go-template/model"
	"go-template/usecase/account"
	"go-template/util"
)

type Route struct {
	accountUsecase account.IUsecase
}

func Init(group *echo.Group, accountUsecase account.IUsecase) {
	r := &Route{accountUsecase}
	group.POST("/signup", r.SignUp)
	group.POST("/signin", r.SignIn)
	group.PUT("/:id/edit", r.EditUser)
	group.POST("/:id", r.ShowUserDetail)
	group.POST("/userlist", r.ShowUserList)
	group.DELETE("/:id", r.DeleteUser)

}

// Dang ky
func (r *Route) SignUp(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	//account tu front end truyen len
	account := &model.Account{}

	//gan tu front end len bien account
	err := c.Bind(account)
	if err != nil {
		util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
		return nil
	}

	log.Println(account)

	//truyen vao useCase xu ly logic
	account, err = r.accountUsecase.CreateAccount(ctx, account)
	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, account)
	return nil
}

// Thay ten ham thoi la dung duoc, con lai la xu ly loi

//Dang nhap
func (r *Route) SignIn(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	//account tu front end truyen len
	account := &model.Account{}

	//gan tu front end len bien account
	err := c.Bind(account)
	if err != nil {
		util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
		return nil
	}

	log.Println(account)

	//truyen vao useCase xu ly logic
	//account with token
	accountwtk, err := r.accountUsecase.SignIn(ctx, account)
	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, accountwtk)
	return nil

}

func (r *Route) EditUser(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	//account tu front end truyen len
	account := &model.Account{}

	//gan tu front end len bien account
	err := c.Bind(account)
	if err != nil {
		util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
		return nil
	}

	log.Println(account)

	id := c.Param("id")
	idint64, err := strconv.ParseInt(id, 10, 64)

	//truyen vao useCase xu ly logic

	account.ID = idint64

	account, err = r.accountUsecase.EditAccount(ctx, account)

	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, account)
	return nil

}

func (r *Route) ShowUserList(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	// //account tu front end truyen len
	// account := &model.Account{}

	// //gan tu front end len bien account
	// err := c.Bind(account)
	// if err != nil {
	// 	util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
	// 	return nil
	// }

	// id := c.Param("id")
	// idint64, err := strconv.ParseInt(id, 10, 64)

	//truyen vao useCase xu ly logic

	// account.ID = idint64

	accounts, err := r.accountUsecase.ShowUserList(ctx)

	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, accounts)
	return nil

}

func (r *Route) ShowUserDetail(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	//account tu front end truyen len
	account := &model.Account{}

	//gan tu front end len bien account
	err := c.Bind(account)
	if err != nil {
		util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
		return nil
	}

	log.Println(account)

	id := c.Param("id")
	idint64, err := strconv.ParseInt(id, 10, 64)

	//truyen vao useCase xu ly logic

	account.ID = idint64

	account, err = r.accountUsecase.UserDetail(ctx, account)

	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, account)
	return nil

}

func (r *Route) DeleteUser(c echo.Context) error {
	ctx := &util.CustomEchoContext{c}

	//account tu front end truyen len
	account := &model.Account{}

	//gan tu front end len bien account
	err := c.Bind(account)
	if err != nil {
		util.Response.Error(ctx, util.NewError(err, http.StatusNotAcceptable, 1000, "invalid input"))
		return nil
	}

	log.Println(account)

	id := c.Param("id")
	idint64, err := strconv.ParseInt(id, 10, 64)

	//truyen vao useCase xu ly logic

	account.ID = idint64

	account, err = r.accountUsecase.DeleteAccount(ctx, account)

	if err != nil {
		util.Response.Error(c, err.(util.MyError))
		return nil
	}

	util.Response.Success(c, account)
	return nil

}
