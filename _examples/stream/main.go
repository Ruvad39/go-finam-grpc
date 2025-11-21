package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ruvad39/go-finam-grpc"
	v1 "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1"
	market_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	order_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/orders"
	"github.com/joho/godotenv"
	slogw "github.com/yougg/slog-writer"
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

	// init logger
	log_app := InitLogger("logs/stream_app.log")

	// создаем клиент
	client, err := finam.NewClient(ctx, token, finam.WithLogger(log_app))
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	//------------------------------------------
	// создадим поток ордеров и сделок
	// order_log := InitLogger("logs/order.log")
	// slog.SetDefault(order_log)
	// newOrderTradeStream(ctx, client)

	//------------------------------------------
	// bar stream
	// логер
	bar_log := InitLogger("logs/bar.log")
	slog.SetDefault(bar_log)

	// пример создание стрима с callback функций
	// NewBarStreamWithCallback(ctx, client)

	// пример создание стрима с возвратом канала
	NewBarStreamWithChannel(ctx, client)

	//--------------------------------------------
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()
	slog.Info("exiting...")

}

// создадим поток ордеров и сделок
func newOrderTradeStream(ctx context.Context, client *finam.Client) {
	slog.Info("newOrderTradeStream", "accountID", accountID)
	stream := client.NewOrderTradeStream(ctx, accountID, onOrder, onTrade)
	_ = stream
	// stream.Close()

}

// NewBarStreamWithCallback пример создание стрима с callback функцией
func NewBarStreamWithCallback(ctx context.Context, client *finam.Client) {
	symbol := "SBER@MISX"
	tf := market_service.TimeFrame_TIME_FRAME_M1
	slog.Info("NewBarStreamWithCallback", "symbol", symbol, "tf", tf.String())

	stream := client.NewBarStreamWithCallback(ctx, symbol, tf, onBar)
	// запуск стрима
	stream.Start()

}

// NewBarStreamWithChannel пример создание стрима с возвратом канала
func NewBarStreamWithChannel(ctx context.Context, client *finam.Client) {
	symbol := "SBER@MISX"
	tf := market_service.TimeFrame_TIME_FRAME_M1
	slog.Info("NewBarStreamWithChannel", "symbol", symbol, "tf", tf.String())

	stream, barChan := client.NewBarStreamWithChannel(ctx, symbol, tf)
	// запустим чтение канала
	go func() {
		for bar := range barChan {
			// обработка
			onBar(bar)
		}
	}()

	// запуск стрима
	stream.Start()

}

// callback метод для обработки ордеров
func onOrder(order *order_service.OrderState) {
	slog.Info("OnOrder", slog.Any("AccountOrder", order))
	//fmt.Printf("onOrder: %v\n", order)
}

// callback  метод для обработки сделок
func onTrade(trade *v1.AccountTrade) {
	slog.Info("onTrade", slog.Any("AccountTrade", trade))
}

// callback метод для обработки баров
func onBar(bar *finam.Bar) {
	slog.Info("onBar", slog.Any("bar", bar.String()))
}

// создать файл логгера
func InitLogger(fileName string) *slog.Logger {
	fw := &slogw.FileWriter{
		Filename:     fileName,
		EnsureFolder: true,
		MaxBackups:   1,
		MaxSize:      1 * 1024 * 1024,
		FileMode:     0644,
		TimeFormat:   slogw.TimeFormatUnix,
		LocalTime:    true,
		ProcessID:    false,
	}
	writer := io.MultiWriter(os.Stdout, fw)
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
