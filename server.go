package main

import (
	"fmt"
	"github.com/6210530017/assessment/config"
)

func main() {
	config := config.NewConfig()

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", config.Port)
	fmt.Println("Database_URL: ", config.DB_url)
}
