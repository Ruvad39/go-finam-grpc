package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-grpc"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	token, _ := os.LookupEnv("FINAM_TOKEN")

	slog.Info("start")
	finam.SetLogDebug(true)
	// создание клиента
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	stream := client.NewStream()
	// подпишемся на котировки (Quote)
	stream.Subscribe(finam.QuoteChannel, "SIM5@RTSX")
	//stream.Subscribe(finam.QuoteChannel, "ROSN@MISX")
	stream.Subscribe(finam.QuoteChannel, "SBER@MISX")

	// установим метод обработчик данных (раньше StartStream)
	stream.SetQuoteHandler(onQuote)
	//  установим метод обработчик "сырых" данных (раньше StartStream)
	stream.SetRawQuoteHandler(onRawQuote)
	// запустим поток данных
	err = stream.Connect(ctx)
	if err != nil {
		slog.Error("stream.Connect", "err", err.Error())
		return
	}

	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()
	slog.Info("exiting...")
}

// onQuote обработка входящих котировок
func onQuote(quote finam.Quote) {
	//fmt.Printf("onQuote: %v\n", quote)
	slog.Info("onQuote", "time", quote.Time(), "quote", quote)
	_ = quote
}

// обработаем сырые котировки
func onRawQuote(quote *marketdata_service.Quote) {
	fmt.Printf("onRawQuote: %v\n", quote)
	//slog.Info("RawQuoteChan", "time", quote.Timestamp.AsTime().In(finam.TzMoscow), "rawQuote", quote)
	//_ = quote
}

// waitForSignal Ожидание сигнала о закрытие
func waitForSignal(ctx context.Context, signals ...os.Signal) os.Signal {
	var exit = make(chan os.Signal, 1)
	signal.Notify(exit, signals...)
	defer signal.Stop(exit)

	select {
	case sig := <-exit:
		slog.Info("WaitForSignal", "signals", sig)
		return sig
	case <-ctx.Done():
		return nil
	}
	return nil
}
