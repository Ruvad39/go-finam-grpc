package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	v1 "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1"
	pb "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/orders"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
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
	accountId, _ := os.LookupEnv("FINAM_ACCOUNT")
	_ = accountId

	finam.SetLogLevel(slog.LevelDebug)
	// создаем клиент
	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	defer client.Close()

	// создаем клиент для работы c OrdersService
	orderService := client.NewOrderServiceClient()

	// получим список всех ордеров по заданному счету
	//getOrders(ctx, orderService, accountId)

	// получим информацию по заданному ордеру
	//getOrder(ctx, orderService, accountId, "69970064220")

	// Отмена биржевой заявки
	//cancelOrder(ctx, orderService, accountId, "69970064220")

	// Выставление биржевой заявки
	placeOrder(ctx, orderService, accountId)

}

// Запрос списка всех ордеров по заданному счету
func getOrders(ctx context.Context, client *finam.OrderServiceClient, accountId string) {

	orders, err := client.GetOrders(ctx, accountId)
	if err != nil {
		slog.Error("OrdersService.GetOrders", "err", err.Error())
		return
	}

	for n, row := range orders.Orders {
		slog.Info("OrdersService.GetOrders", "n", n,
			"order", row)
	}

}

// Получим информацию по заданному ордеру
func getOrder(ctx context.Context, client *finam.OrderServiceClient, accountId string, orderId string) {
	orderState, err := client.GetOrder(ctx, accountId, orderId)
	if err != nil {
		slog.Error("OrdersService", "err", err.Error())
	}
	slog.Info("OrdersService.GetOrder", "orderState", orderState)
}

// Запрос отмены торговой заявки
func cancelOrder(ctx context.Context, client *finam.OrderServiceClient, accountId string, orderId string) {

	orderState, err := client.CancelOrder(ctx, accountId, orderId)
	if err != nil {
		slog.Error("OrdersService.CancelOrder", "err", err.Error())
	}
	slog.Info("OrdersService.CancelOrder", "orderState", orderState)
}

// Выставление биржевой заявки
func placeOrder(ctx context.Context, client *finam.OrderServiceClient, accountId string) {
	symbol := "GLDRUBF@RTSX" // инструмент SPBFUT.GLDRUBF
	quantity := 1            // кол-во в штуках
	//price := 0                                  // по какой цене выставить лимитку
	orderType := pb.OrderType_ORDER_TYPE_MARKET //pb.OrderType_ORDER_TYPE_LIMIT // лимитная
	side := v1.Side_SIDE_BUY                    // покупка

	newOrder := &pb.Order{
		AccountId: accountId,
		Symbol:    symbol,
		Quantity:  finam.IntToDecimal(quantity),
		//LimitPrice:  finam.Float64ToDecimal(price),
		Type:        orderType,
		Side:        side,
		TimeInForce: pb.TimeInForce_TIME_IN_FORCE_DAY,
	}
	_ = newOrder

	// пошлем ордер на рынок
	orderState, err := client.PlaceOrder(ctx, newOrder)

	// Есть вспомогательные методы для создания ордера
	// покупка по рынку
	// orderState, err := client.PlaceOrder(ctx, client.NewBuyOrder(accountId, symbol, quantity))
	// продажа по рынку
	// orderState, err := client.PlaceOrder(ctx, client.NewSellOrder(accountId, symbol, quantity))
	// лимитная покупка по заданной цене
	//orderState, err := client.PlaceOrder(ctx, client.NewBuyLimitOrder(accountId, symbol, quantity, price))
	// лимитная продажа по заданной цене
	//orderState, err := client.PlaceOrder(ctx, client.NewSellLimitOrder(accountId, symbol, quantity, 311.05))

	if err != nil {
		slog.Error("OrdersService.PlaceOrder", "err", err.Error())
	}
	slog.Info("OrdersService.PlaceOrder", "orderState", orderState)
}
