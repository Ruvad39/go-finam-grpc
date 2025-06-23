package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-grpc"
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
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

	finam.SetLogLevel(slog.LevelDebug)
	// создаем клиент
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создаем клиент для работы MarketDataService
	marketDataService := client.NewMarketDataServiceClient()

	//
	symbol := "SBER@MISX" //"ROSN@MISX"  //"SIU5@RTSX"
	// Получение последней котировки по инструменту
	//getQuote(ctx, marketDataService, symbol)

	// Получение исторических данных по инструменту (агрегированные свечи)
	//getBars(ctx, marketDataService, symbol)

	// Получение текущего стакана по инструменту
	//getOrderBook(ctx, marketDataService, symbol)

	// Получение списка последних сделок по инструменту
	getLatestTrades(ctx, marketDataService, symbol)
}

// Получение последней котировки по инструменту
func getQuote(ctx context.Context, client *finam.MarketDataServiceClient, symbol string) {

	quote, err := client.GetLastQuote(ctx, symbol)
	if err != nil {
		slog.Error("LastQuote", "err", err.Error())
		return
	}
	slog.Info("MarketDataService.LastQuote",
		"Symbol", quote.Symbol,
		"Timestamp", quote.Quote.Timestamp.AsTime().In(finam.TzMoscow),
		"Ask", finam.DecimalToFloat64(quote.Quote.Ask),
		"Bid", finam.DecimalToFloat64(quote.Quote.Bid),
		"Last", finam.DecimalToFloat64(quote.Quote.Last),
		"Additions", quote.Quote.Additions,
	)
	// все данные
	fmt.Printf("LastQuote: %v\n", quote)
}

// Получение исторических данных по инструменту (агрегированные свечи)
func getBars(ctx context.Context, client *finam.MarketDataServiceClient, symbol string) {
	tf := pb.TimeFrame_TIME_FRAME_D
	start, _ := time.Parse("2006-01-02", "2025-05-01")
	end := time.Now()
	// запрос
	bars, err := client.GetBars(ctx, symbol, tf, start, end)
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
func getOrderBook(ctx context.Context, client *finam.MarketDataServiceClient, symbol string) {

	b, err := client.GetOrderBook(ctx, symbol)
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
func getLatestTrades(ctx context.Context, client *finam.MarketDataServiceClient, symbol string) {

	trades, err := client.GetLatestTrades(ctx, symbol)
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
