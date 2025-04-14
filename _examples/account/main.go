package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-grpc"
	accounts_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/accounts"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

func main() {
	// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	token, _ := os.LookupEnv("FINAM_TOKEN")

	slog.Info("start")
	// создание клиента
	ctx := context.Background()
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// добавим заголовок с авторизацией (accessToken)
	ctx, err = client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}

	// Получение информации по конкретному аккаунту
	accountId, _ := os.LookupEnv("FINAM_ACCOUNT_ID")

	res, err := client.AccountsService.GetAccount(ctx, &accounts_service.GetAccountRequest{AccountId: accountId})
	if err != nil {
		slog.Error("accountService", "GetAccount", err.Error())
	}
	slog.Info("main", "Account", res)
	// список позиций
	//slog.Info("main", "Positions", res.Positions)
	for row, pos := range res.Positions {
		slog.Info("positions",
			"row", row,
			"Symbol", pos.Symbol,
			"Quantity", finam.DecimalToFloat64(pos.Quantity),
			"AveragePrice", finam.DecimalToFloat64(pos.AveragePrice),
			"CurrentPrice", finam.DecimalToFloat64(pos.CurrentPrice),
		)
	}

	// Запрос получения истории по сделкам
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end := time.Now()
	req := finam.NewTradesRequest(accountId, limit, start, end)
	res2, err := client.AccountsService.Trades(ctx, req)
	if err != nil {
		slog.Error("accountService", "Trades", err.Error())
	}
	slog.Info("accountService", "Trades", res2)

	// Запрос Получение списка транзакций аккаунта
	res3, err := client.AccountsService.Transactions(ctx, finam.NewTransactionsRequest(accountId, limit, start, end))
	if err != nil {
		slog.Error("accountService", "Transactions", err.Error())
	}
	//slog.Info("accountService", "Transactions", res3)
	fmt.Printf("Transactions = %s \n", res3)

}
