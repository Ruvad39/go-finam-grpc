package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-grpc"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
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
	finam.SetLogDebug(true)
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	symbol := "SBER@MISX" //"ROSN@MISX"  //"SIM5@RTSX"
	// Получение последней котировки по инструменту
	//getQuote(ctx, client, symbol)

	// Получение исторических данных по инструменту (агрегированные свечи)
	//getBars(ctx, client, symbol)

	// Получение текущего стакана по инструменту
	// getOrderBook(ctx, client, symbol)

	// Получение списка последних сделок по инструменту
	getLatestTrades(ctx, client, symbol)

}

// Получение последней котировки по инструменту
func getQuote(ctx context.Context, client *finam.Client, symbol string) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}
	//symbol := "ROSN@MISX" //"SBER@MISX" //"SIM5@RTSX"
	q, err := client.MarketDataService.LastQuote(ctx, finam.NewQuoteRequest(symbol))
	if err != nil {
		slog.Error("LastQuote", "err", err.Error())
		return
	}
	slog.Info("MarketDataService",
		"Symbol", q.Symbol,
		"Timestamp", q.Quote.Timestamp.AsTime().In(finam.TzMoscow),
		"Ask", finam.DecimalToFloat64(q.Quote.Ask),
		"Bid", finam.DecimalToFloat64(q.Quote.Bid),
		"Last", finam.DecimalToFloat64(q.Quote.Last),
		"Additions", q.Quote.Additions,
	)
	// все данные
	fmt.Printf("LastQuote: %v\n", q)

}

// Получение исторических данных по инструменту (агрегированные свечи)
func getBars(ctx context.Context, client *finam.Client, symbol string) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}

	//symbol := "SBER@MISX" //"SIM5@RTSX" MISX
	// получение списка свечей
	tf := marketdata_service.BarsRequest_TIME_FRAME_D
	start, _ := time.Parse("2006-01-02", "2025-01-01")
	end := time.Now()
	req := finam.NewBarsRequest(symbol, tf, start, end)
	bars, err := client.MarketDataService.Bars(ctx, req)
	if err != nil {
		slog.Error("MarketDataService.Bars", "err", err.Error())
		return
	}
	slog.Info("MarketDataService.Bars", "Bars.len", len(bars.Bars))
	for row, bar := range bars.Bars {
		slog.Info("Bar", "row", row,
			"Timestamp", bar.Timestamp.AsTime().In(finam.TzMoscow),
			"Open", finam.DecimalToFloat64(bar.Open),
			"High", finam.DecimalToFloat64(bar.High),
			"Low", finam.DecimalToFloat64(bar.Low),
			"Close", finam.DecimalToFloat64(bar.Close),
			"Volume_int", finam.DecimalToInt(bar.Volume),
		)
	}

}

// Получение текущего стакана по инструменту
func getOrderBook(ctx context.Context, client *finam.Client, symbol string) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}
	b, err := client.MarketDataService.OrderBook(ctx, finam.NewOrderBookRequest(symbol))
	if err != nil {
		slog.Error("MarketDataService.OrderBook", "err", err.Error())
		return
	}
	slog.Info("MarketDataService.OrderBook",
		"Symbol", b.Symbol,
		//"OrderBook", b.Orderbook,
	)
	for n, row := range b.Orderbook.Rows {
		slog.Info("OrderBook.Rows", "n", n,
			"row", row,
		)

	}

}

// Получение списка последних сделок по инструменту
func getLatestTrades(ctx context.Context, client *finam.Client, symbol string) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}
	trades, err := client.MarketDataService.LatestTrades(ctx, finam.NewLatestTradesRequest(symbol))
	if err != nil {
		slog.Error("MarketDataService.LatestTrades", "err", err.Error())
		return
	}
	slog.Info("MarketDataService.LatestTrades", "Symbol", trades.Symbol)
	for n, row := range trades.Trades {
		slog.Info("Trade",
			"n", n,
			"time", row.Timestamp.AsTime().In(finam.TzMoscow),
			"trades", row,
		)
	}

}
