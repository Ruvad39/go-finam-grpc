package main

import (
	"context"
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
	//
	streamQuote(ctx, client)

	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()

	slog.Info("exiting...")
}

func streamQuote(ctx context.Context, client *finam.Client) {
	// пидпишемся
	//client.Subscribe(finam.QuoteChannel, "SIM5@RTSX")
	//client.Subscribe(finam.QuoteChannel, "ROSN@MISX")
	client.Subscribe(finam.QuoteChannel, "MOEX@MISX")
	_ = client.StartStream(ctx)

	// получим канал с котировкам
	quoteChan := client.GetQuoteChan()
	// пошлем его на обработку
	go listenQuoteChan(ctx, quoteChan)
}

// читаем канал с котировками
func listenQuoteChan(ctx context.Context, quoteChan chan *marketdata_service.Quote) {
	for {
		select {
		case res := <-quoteChan:
			//fmt.Printf("quoteСhan: %v\n", res)
			slog.Info("listenQuoteChan", "quote", res)
		case <-ctx.Done():
			return
		}
	}
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
