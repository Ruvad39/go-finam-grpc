package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Ruvad39/go-finam-grpc"
	pb "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/assets"
	"github.com/joho/godotenv"
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

	finam.SetLogLevel(slog.LevelDebug)
	// создаем клиент
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создаем клиент для доступа к AssetService
	assetService := client.NewAssetServiceClient()

	// получим текущее время сервера
	currTime, err := assetService.GetTime(ctx)
	if err != nil {
		slog.Error("AssetsService.Clock", "err", err.Error())
	}
	slog.Info("main", "current time server", currTime)

	// Получение списка доступных бирж, названия и mic коды
	//getExchanges(ctx, assetService)

	// Получение списка доступных инструментов, их описание
	//getAssets(ctx, assetService)

	// возьмем номер счета из .env
	accountId, _ := os.LookupEnv("FINAM_ACCOUNT")
	_ = accountId
	symbol := "IMOEXF@RTSX" //"SBER@MISX" //"FEES@MISX" //"SBER@MISX" // "SIU5@RTSX" // "SBER@MISX" // "RU000A106L18@MISX"
	_ = symbol
	// Получение информации по конкретному инструменту
	getAsset(ctx, assetService, accountId, symbol)

	// Получение торговых параметров по инструменту
	//getAssetParams(ctx, assetService, accountId, symbol)

	// Получение расписания торгов для инструмента
	//getSchedule(ctx, assetService, symbol)
}

// Получение списка доступных бирж, названия и mic коды
func getExchanges(ctx context.Context, client *finam.AssetServiceClient) {
	// Получение списка доступных бирж, названия и mic коды
	exchanges, err := client.GetExchanges(ctx)
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
func getAssets(ctx context.Context, client *finam.AssetServiceClient) {
	assets, err := client.GetAssets(ctx)
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

// Получение информации по конкретному инструменту
func getAsset(ctx context.Context, client *finam.AssetServiceClient, accountId, symbol string) {
	assetInfo, err := client.GetAsset(ctx, accountId, symbol)
	if err != nil {
		slog.Error("AssetsInfo", "err", err.Error())
	}
	slog.Info("AssetsService.GetAsset", slog.Any("assetInfo", assetInfo))

}

// Получение торговых параметров по инструменту
func getAssetParams(ctx context.Context, client *finam.AssetServiceClient, accountId, symbol string) {
	assetParams, err := client.GetAssetParams(ctx, accountId, symbol)
	if err != nil {
		slog.Error("GetAssetParams", "err", err.Error())
	}
	slog.Info("AssetsService.GetAssetParams", slog.Any("assetParams", assetParams))
}

// Получение расписания торгов для инструмента
func getSchedule(ctx context.Context, client *finam.AssetServiceClient, symbol string) {
	res, err := client.GetSchedule(ctx, symbol)
	if err != nil {
		slog.Error("AssetsService.Schedule", "err", err.Error())
	}

	// текущее время (компа)
	currTime, _ := client.GetTime(ctx)
	currSession := &pb.ScheduleResponse_Sessions{}
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
