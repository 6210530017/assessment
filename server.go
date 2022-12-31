package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/6210530017/assessment/config"
	"github.com/6210530017/assessment/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

func Setup(url string) (*sql.DB, func()) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	teardown := func() {
		db.Close()
	}

	return db, teardown
}

func main() {
	config := config.NewConfig()

	db, teardown := Setup(config.DB_url)
	defer teardown()

	h := handler.NewHandler(db)

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	e.POST("/expenses", h.CreateExpense)

	go func() {
		if err := e.Start(config.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
