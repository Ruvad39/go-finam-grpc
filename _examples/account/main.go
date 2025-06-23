package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-grpc"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

// предполагаем что есть файл .env
// в котором записан secret-Token в переменной FINAM_TOKEN
// и номер счета в FINAM_ACCOUNT_ID
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	ctx := context.Background()
	// получаем переменные из .env
	token, _ := os.LookupEnv("FINAM_TOKEN")

	//finam.SetLogLevel(slog.LevelDebug)
	// создаем клиент
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создаем клиент для доступа к AccountService
	accountService := client.NewAccountServiceClient()

	// Получение информации о токене сессии. Возьмем список счетов
	res, err := client.GetTokenDetails(ctx)
	if err != nil {
		slog.Error("main", "AuthService.TokenDetails", err.Error())
		return
	}
	for row, accountId := range res.AccountIds {
		slog.Info("TokenDetails.AccountIds", "row", row, "accountId", accountId)

		// получим информацию по конкретному счету
		//getAccount(ctx, accountService, account_Id)

		// получим список сделок
		//getTrades(ctx, accountService, accountId)

		getTransactions(ctx, accountService, accountId)
	}
}

func getAccount(ctx context.Context, client *finam.AccountServiceClient, accountId string) {
	res, err := client.GetAccount(ctx, accountId)
	if err != nil {
		slog.Error("AccountsService.GetAccount", "GetAccount", err.Error())
		return
	}
	slog.Info("AccountsService.GetAccount",
		"AccountId", res.AccountId,
		"Type", res.Type,
		"Status", res.Status,
		"Equity", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.Equity)),
		"UnrealizedProfit", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.UnrealizedProfit)),
		"Cash", res.Cash,
	)
	// список позиций
	for row, pos := range res.Positions {
		slog.Info("AccountsService.GetAccount.Positions",
			"row", row,
			"Symbol", pos.Symbol,
			"Quantity", finam.DecimalToFloat64(pos.Quantity),
			"AveragePrice", finam.DecimalToFloat64(pos.AveragePrice),
			"CurrentPrice", finam.DecimalToFloat64(pos.CurrentPrice),
		)

	}
}

// Запрос Получение списка сделок
func getTrades(ctx context.Context, client *finam.AccountServiceClient, accountId string) {
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end := time.Now()

	res, err := client.GetTrades(ctx, accountId, start, end, limit)
	if err != nil {
		slog.Error("AccountsService.Trades", "err", err.Error())
		return
	}
	slog.Info("AccountsService.Trades", "кол-во сделок", len(res.Trades))
	for row, t := range res.Trades {
		slog.Info("AccountsService.Trades", "row", row, "trade", t)
	}
}

// Запрос Получение списка транзакций аккаунта
func getTransactions(ctx context.Context, client *finam.AccountServiceClient, accountId string) {
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end := time.Now()

	res, err := client.GetTransactions(ctx, accountId, start, end, limit)
	if err != nil {
		slog.Error("AccountsService.Transactions", "err", err.Error())
		return
	}
	slog.Info("AccountsService.Transactions", "кол-во операций", len(res.Transactions))
	for row, t := range res.Transactions {
		slog.Info("AccountsService.Transactions", "row", row, "transaction", t)
	}
}
