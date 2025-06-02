package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Ruvad39/go-finam-grpc"
	assets_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/assets"
	"github.com/joho/godotenv"
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
	finam.SetLogDebug(true) // проставим признак отладки
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
	//getExchanges(ctx, client)
	// test Получение списка доступных бирж =  подсчет кол-ва дублей
	// getExchangesCount(ctx, client)

	// Получение списка доступных инструментов, их описание
	//getAssets(ctx, client)
	//
	accountId, _ := os.LookupEnv("FINAM_ACCOUNT_ID")
	getGetAssetParams(ctx, client, "SBER@MISX", accountId)

	// Получение расписания торгов для инструмента
	//getSchedule(ctx, client, "SBER@MISX")

}

// Получение списка доступных бирж, названия и mic коды
func getExchanges(ctx context.Context, client *finam.Client) {
	// Получение списка доступных бирж, названия и mic коды
	exchanges, err := client.AssetsService.Exchanges(ctx, finam.NewExchangesRequest())
	if err != nil {
		slog.Error("AssetsService.Exchanges", " GeExchanges", err.Error())
	}
	slog.Info("main", "Exchanges.len", len(exchanges.Exchanges))
	for row, ex := range exchanges.Exchanges {
		slog.Info("exchanges", "row", row,
			"Exchanges", ex)
	}

}

// Получение списка доступных инструментов, их описание
func getAssets(ctx context.Context, client *finam.Client) {

	assets, err := client.AssetsService.Assets(ctx, finam.NewAssetsRequest())
	if err != nil {
		slog.Error("AssetsService.Assets", "err", err.Error())
	}
	slog.Info("AssetsService.Assets", "assets.len", len(assets.Assets))

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

func getGetAssetParams(ctx context.Context, client *finam.Client, symbol, accountId string) {
	assetParams, err := client.AssetsService.GetAssetParams(ctx, finam.NewAssetParamsRequest(symbol, accountId))
	if err != nil {
		slog.Error("AssetsParams", "err", err.Error())
	}
	slog.Info("AssetsService.AssetsParams", slog.Any("assetParams", assetParams))
}

// Получение расписания торгов для инструмента
func getSchedule(ctx context.Context, client *finam.Client, symbol string) {
	res, err := client.AssetsService.Schedule(ctx, finam.NewScheduleRequest(symbol))
	if err != nil {
		slog.Error("AssetsService.Schedule", "err", err.Error())
	}

	// текущее время (компа)
	currTime := time.Now().In(finam.TzMoscow)
	currSession := &assets_service.ScheduleResponse_Sessions{}
	// список всех сессий
	for n, session := range res.Sessions {
		slog.Info("Sessions", "n", n,
			"Type", session.Type,
			"start", session.Interval.StartTime.AsTime().In(finam.TzMoscow),
			"end", session.Interval.EndTime.AsTime().In(finam.TzMoscow),
		)
		// определим текущую сессиию
		if finam.IsWithinInterval(currTime, session.Interval) {
			currSession = session
		}
	}

	slog.Info("текущая сессия",
		"currTime", currTime,
		"Type", currSession.Type,
		"start", currSession.Interval.StartTime.AsTime().In(finam.TzMoscow),
		"end", currSession.Interval.EndTime.AsTime().In(finam.TzMoscow),
	)

}

// TEST Получение списка доступных бирж, подсчет кол-ва дублей
func getExchangesCount(ctx context.Context, client *finam.Client) {
	res, err := client.AssetsService.Exchanges(ctx, finam.NewExchangesRequest())
	if err != nil {
		slog.Error("AssetsService.Exchanges", " GeExchanges", err.Error())
	}
	// посчитаем кол-во дублей
	counts := make(map[string]int)
	for _, item := range res.Exchanges {
		counts[item.Mic]++
	}
	// Вывод повторяющихся
	fmt.Println("список записей, которые повторяются:")
	for mic, count := range counts {
		if count > 1 {
			fmt.Printf("%s = %d \n", mic, count)
		}
	}
}
