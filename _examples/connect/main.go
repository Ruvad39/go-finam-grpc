package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	ctx := context.Background()
	token, _ := os.LookupEnv("FINAM_TOKEN")
	//accountId ,_ := os.LookupEnv("FINAM_ACCOUNT_ID")
	//
	slog.Info("start connect")
	finam.SetLogLevel(slog.LevelDebug)

	client, err := finam.NewClient(ctx, token, "")
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// Получение информации о токене сессии
	res, err := client.GetTokenDetails(ctx)
	if err != nil {
		slog.Error("main", "AuthService.TokenDetails", err.Error())
	}
	slog.Info("main", "res", res)
	// список счетов
	slog.Info("main", "res.AccountIds", res.AccountIds)
}
