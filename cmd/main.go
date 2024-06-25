package main

import (
	"fmt"
	"userServerAuth/internal/config"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)
	// инициализировать логгер
	// инициализировать точку входа в приложение
	// запустить grpc сервер
}
