package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
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

	// Получение списка доступных бирж, названия и mic коды
	exchanges, err := client.AssetsService.Exchanges(ctx, finam.NewExchangesRequest())
	if err != nil {
		slog.Error("main", " GeExchanges", err.Error())
	}
	slog.Info("main", "Exchanges.len", len(exchanges.Exchanges))
	slog.Info("main", "Exchanges", exchanges.Exchanges)

	return
	// Получение списка доступных инструментов, их описание
	assets, err := client.AssetsService.Assets(ctx, finam.NewAssetsRequest())
	if err != nil {
		slog.Error("AssetsService", "Assets", err.Error())
	}
	slog.Info("main", "GeAssets.len", len(assets.Assets))
	//fmt.Printf("всего иструментов = %d \n", len(assets.Assets))
	for n, sec := range assets.Assets {
		//if sec.Type == "FUTURES" && sec.Mic == "RTSX" && 1 == 2 {
		if 1 == 1 {
			slog.Info("assets",
				"row", n,
				"id", sec.Id,
				"Symbol", sec.Symbol,
				"Name", sec.Name,
				"ticker", sec.Ticker,
				"Type", sec.Type,
				"Mic", sec.Mic,
			)
		}
	}
}
