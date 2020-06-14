package util

import (
	"log"

	"github.com/gin-gonic/gin"
	echo "github.com/labstack/echo/v4"
)

type response struct {
}

var Response response

func (response) Success(c echo.Context, data interface{}) {
	c.JSON(200, data)
}

func (response) Error(c echo.Context, err MyError) {
	log.Println(err.Raw)
	c.JSON(err.HTTPCode, gin.H{
		"code":    err.ErrorCode,
		"message": err.Message,
	})
}
