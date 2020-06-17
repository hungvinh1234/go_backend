package http

import (
	"go-template/delivery/http/account"
	"go-template/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewHTTPHandler .
func NewHTTPHandler(ucase *usecase.Usecase) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	account.Init(e.Group("/account"), ucase.Account)

	return e
}
