package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

var token string
var accountID string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	// получаем переменные из .env
	token, _ = os.LookupEnv("FINAM_TOKEN")
	accountID, _ = os.LookupEnv("FINAM_ACCOUNT")

	finam.SetLogLevel(slog.LevelDebug)
	// создаем клиент
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создадим поток ордеров и сделок
	newOrderTradeStream(ctx, client)

	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()
	slog.Info("exiting...")

}

// создадим поток ордеров и сделок
func newOrderTradeStream(ctx context.Context, client *finam.Client) {
	stream := client.NewOrderTradeStream()
	// подпишемся на ордера и сделки по счету
	// OrderTradeChannel = Подписка и на Ордера и на Сделки
	// OrderChannel = Подписка только на Ордера
	// TradeChannel = Подписка только на Сделки
	// номер счета + тип подписки
	stream.Subscribe(accountID, finam.OrderTradeChannel)

	// установим функцию обработчик
	// в этом случае автоматом запуститься startOrderWorker (чтение канала и отправка данных в onOrder)
	stream.SetOnOrder(onOrder)
	stream.SetOnTrade(onTrade)

	// запустим поток в работу
	err := stream.Start(ctx)
	if err != nil {
		slog.Error("OrderTradeStream.Start", "err", err.Error())
	}
}

// callback
func onOrder(order pb.OrderState) {
	slog.Info("OnOrder", slog.Any("AccountOrder", order))
	//fmt.Printf("onOrder: %v\n", order)
}

// callback
func onTrade(trade pb.AccountTrade) {
	slog.Info("onTrade", slog.Any("AccountTrade", trade))
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
