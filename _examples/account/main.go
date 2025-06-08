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

	finam.SetLogDebug(true)
	slog.Info("start")
	// создание клиента
	ctx := context.Background()
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	//accountId, _ := os.LookupEnv("FINAM_ACCOUNT_ID")
	//getTrades_new(ctx, client, accountId)
	//return

	// Получение информации о токене сессии. Возьмем список счетов
	res, err := client.GetTokenDetails(ctx)
	if err != nil {
		slog.Error("main", "AuthService.TokenDetails", err.Error())
	}
	for row, accountId := range res.AccountIds {
		slog.Info("TokenDetails.AccountIds", "row", row, "accountId", accountId)
		// получим информацию по конкретному счету
		//getAccount_new(ctx, client, accountId)
		//getPositions(ctx, client, accountId)
		//getAccount(ctx, client, accountId)
		getTrades_new(ctx, client, accountId)
		//getTransactions(ctx, client, accountId)
	}

}

// getAccount_new получим информацию по конкретному счету
// используя AccountRequest = обертка для более "простого" получение данных
func getAccount_new(ctx context.Context, client *finam.Client, accountId string) {
	res, err := client.NewAccountRequest(accountId).Do(ctx)
	if err != nil {
		slog.Error("AccountRequest.Do", "GetAccount", err.Error())
		return
	}
	slog.Info("AccountRequest.Do",
		"AccountId", res.AccountId,
		"Type", res.Type,
		"Status", res.Status,
		"Equity", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.Equity)),
		"UnrealizedProfit", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.UnrealizedProfit)),
		"Cash", res.Cash,
	)
}

func getPositions(ctx context.Context, client *finam.Client, accountId string) {
	res, err := client.NewAccountRequest(accountId).Do(ctx)
	if err != nil {
		slog.Error("AccountRequest.Do", "GetAccount", err.Error())
		return
	}
	slog.Info("getPositions", "len(Positions)", len(res.Positions))
	// список позиций
	for row, pos := range res.Positions {
		slog.Info("AccountsService.GetAccount.Positions",
			"row", row,
			"Symbol", pos.Symbol,
			"Quantity", finam.DecimalToInt(pos.Quantity),
			"AveragePrice", finam.DecimalToFloat64(pos.AveragePrice),
			"CurrentPrice", finam.DecimalToFloat64(pos.CurrentPrice),
		)

	}
}

// getAccount получим информацию по конкретному счету
func getAccount(ctx context.Context, client *finam.Client, accountId string) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}
	res, err := client.AccountsService.GetAccount(ctx, &accounts_service.GetAccountRequest{AccountId: accountId})
	if err != nil {
		slog.Error("AccountsService.GetAccount", "GetAccount", err.Error())
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

// Запрос получения истории по сделкам
func getTrades_new(ctx context.Context, client *finam.Client, accountId string) {
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start_time := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end_time := time.Now()

	resp, err := client.NewAccountTradesRequest(accountId).StartTime(start_time).EndTime(end_time).Limit(limit).Do(ctx)
	if err != nil {
		slog.Error("accountService", "Trades", err.Error())
	}
	slog.Info("accountService.Trades", "кол-во сделок", len(resp.Trades))
	for row, t := range resp.Trades {
		slog.Info("AccountsService.Trades", "row", row, "trade", t)
	}
}

// Запрос получения истории по сделкам
func getTrades(ctx context.Context, client *finam.Client, accountId string) {
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end := time.Now()

	req := finam.NewTradesRequest(accountId, limit, start, end)
	res, err := client.AccountsService.Trades(ctx, req)
	if err != nil {
		slog.Error("accountService", "Trades", err.Error())
	}
	for row, t := range res.Trades {
		slog.Info("AccountsService.Trades", "row", row, "trade", t)
	}
}

// Запрос Получение списка транзакций аккаунта
func getTransactions(ctx context.Context, client *finam.Client, accountId string) {
	// запросим данные за последние 24 часа
	var limit int32 = 0
	start := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end := time.Now()

	res, err := client.AccountsService.Transactions(ctx, finam.NewTransactionsRequest(accountId, limit, start, end))
	if err != nil {
		slog.Error("accountService", "Transactions", err.Error())
	}
	for row, t := range res.Transactions {
		slog.Info("AccountsService.Transactions", "row", row, "transaction", t)
	}

}
