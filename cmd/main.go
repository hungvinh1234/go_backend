package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-template/client/postgres"
	"go-template/config"
	serviceHttp "go-template/delivery/http"
	"go-template/repository"
	"go-template/usecase"
)

func main() {
	cfg := config.GetConfig()

	defer postgres.Disconnect()

	// setup locale
	{
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		time.Local = loc
	}

	repo := repository.New(postgres.GetClient)

	ucase := usecase.New(repo)

	errs := make(chan error)

	// http
	{
		h := serviceHttp.NewHTTPHandler(ucase)
		go func() {
			errs <- h.Start(":" + cfg.Port)
		}()
	}

	// graceful
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("exit", <-errs)
}
