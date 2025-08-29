package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	v1 "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1"
	market_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	order_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/orders"
	"github.com/joho/godotenv"
	"io"
	"time"

	slogw "github.com/yougg/slog-writer"
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

	// init logger
	log_system := InitLogger("logs/stream_test.log")
	finam.SetLogger(log_system)
	log_order := InitLogger("logs/order.log")
	slog.SetDefault(log_order)

	// создаем клиент
	client, err := finam.NewClient(ctx, token, finam.WithJwtRefreshInterval(3*time.Minute))
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создадим поток ордеров и сделок
	newOrderTradeStream(ctx, client)

	symbol := "SBER@MISX"
	//tf1 := market_service.TimeFrame_TIME_FRAME_M1
	tf := market_service.TimeFrame_TIME_FRAME_D
	_ = symbol
	_ = tf
	//newBarStream(ctx, client, symbol, tf)

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

func newBarStream(ctx context.Context, client *finam.Client, symbol string, tf market_service.TimeFrame) {
	slog.Info("newBarStream", "symbol", symbol, "tf", tf.String())
	stream := client.NewBarStream(ctx, symbol, tf, onBar)
	_ = stream
	// stream.Close()

}

// callback
func onOrder(order *order_service.OrderState) {
	slog.Info("OnOrder", slog.Any("AccountOrder", order))
	//fmt.Printf("onOrder: %v\n", order)
}

// callback
func onTrade(trade *v1.AccountTrade) {
	slog.Info("onTrade", slog.Any("AccountTrade", trade))
}

// callback
func onBar(bar *finam.Bar) {
	slog.Info("onBar", slog.Any("bar", bar.String()))
}

func InitLogger(fileName string) *slog.Logger {
	fw := &slogw.FileWriter{
		Filename:     fileName,
		EnsureFolder: false,
		MaxBackups:   1,
		MaxSize:      1024,
		FileMode:     0644,
		TimeFormat:   slogw.TimeFormatUnix,
		LocalTime:    true,
		ProcessID:    false,
	}
	writer := io.MultiWriter(os.Stdout, fw)
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
