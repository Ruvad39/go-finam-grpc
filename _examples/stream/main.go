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

	// подпишемся на котировки (Quote)
	//client.Subscribe(finam.QuoteChannel, "SIM5@RTSX")
	//client.Subscribe(finam.QuoteChannel, "ROSN@MISX")
	client.Subscribe(finam.QuoteChannel, "SBER@MISX")
	// установим метод обработчик данных (раньше StartStream)
	client.SetQuoteHandler(onQuote)
	// запустим потока данных
	err = client.StartStream(ctx)
	if err != nil {
		slog.Error("StartStream", "err", err.Error())
	}

	// пример работы с каналом "сырых" данных
	// получим канал с котировкам
	client.SendRawQuotes = true // проставим признак отправки данных в канал
	rawQuoteChan := client.GetRawQuoteChan()
	// пошлем его на обработку
	go listenRawQuoteChan(ctx, rawQuoteChan)

	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()
	slog.Info("exiting...")
}

// onQuote обработка входящих котировок
func onQuote(quote finam.Quote) {
	//fmt.Printf("onQuote: %v\n", quote)
	slog.Info("onQuote", "time", quote.Time(), "quote", quote)
}

// Читаем канал с сырыми данными котировок
func listenRawQuoteChan(ctx context.Context, quoteChan chan *marketdata_service.Quote) {
	for {
		select {
		case res := <-quoteChan:
			//fmt.Printf("RawQuote: %v\n", res)
			slog.Info("RawQuoteChan", "time", res.Timestamp.AsTime().In(finam.TzMoscow), "rawQuote", res)
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
