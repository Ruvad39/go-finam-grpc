package main

import (
	"context"
	"github.com/Ruvad39/go-finam-grpc"
	side "github.com/Ruvad39/go-finam-grpc/trade_api/v1"
	orders_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/orders"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	token, _ := os.LookupEnv("FINAM_TOKEN")
	accountId, _ := os.LookupEnv("FINAM_ACCOUNT_ID")
	_ = accountId
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

	// добавим заголовок с авторизацией (accessToken)
	ctx, err = client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return
	}

	// получим список всех ордеров по заданному счету
	getOrders(ctx, client, accountId)

	// получим информацию по заданному ордеру
	//orderId := "1892950515808633522"
	//getOrder(ctx, client, accountId, orderId)

	// Отмена биржевой заявки
	cancelOrder(ctx, client, accountId, "40045338")
	//getOrders(ctx, client, accountId)
	// еще раз запросим данные по заданному ордеру
	// уже вернется ошибка ("rpc error: code = NotFound desc = Order with id 2033125054207932011 is not found")
	//getOrder(ctx, client, accountId, orderId)

	// пример выставления ордера на покупку\продажу
	//placeOrder(ctx, client, accountId)
	//buyLimit(ctx, client, accountId)
}

// getOrders получим список всех ордеров по заданному счету
func getOrders(ctx context.Context, client *finam.Client, accountId string) {
	orders, err := client.OrdersService.GetOrders(ctx, finam.NewOrdersRequest(accountId))
	if err != nil {
		slog.Error("OrdersService", "GetOrders", err.Error())
	}
	slog.Info("OrdersService", "orders", orders)
}

// getOrder получим информацию по заданному ордеру
func getOrder(ctx context.Context, client *finam.Client, accountId string, orderId string) {
	orderState, err := client.OrdersService.GetOrder(ctx, finam.NewGetOrderRequest(accountId, orderId))
	if err != nil {
		slog.Error("OrdersService", "GetOrder", err.Error())
	}
	slog.Info("OrdersService", "orderState", orderState)
}

// cancelOrder Запрос отмены торговой заявки
func cancelOrder(ctx context.Context, client *finam.Client, accountId string, orderId string) {
	orderState, err := client.OrdersService.CancelOrder(ctx, finam.NewCancelOrderRequest(accountId, orderId))
	if err != nil {
		slog.Error("OrdersService", "CancelOrder", err.Error())
	}
	slog.Info("OrdersService", "orderState", orderState)
}

func placeOrder(ctx context.Context, client *finam.Client, accountId string) {
	symbol := "SiM5@RTSX"
	//symbol := "SBER@MISX"
	// пример покупки
	//newOrder := finam.NewOrderBuy(accountId, symbol, 1) // по рынку
	//newOrder := finam.NewOrderBuyLimit(accountId, symbol, 1, 2840.0) // лимиткой
	// пример продажи по рынку
	//newOrder := finam.NewOrderSell(accountId, symbol, 1)
	//newOrder := finam.NewOrderSellLimit(accountId, symbol, 1, 2851.0) // лимиткой

	//newOrder := finam.NewOrderBuyStopLimit(accountId, symbol, 1, 86994, 86994) // стоп-лимиткой
	newOrder := finam.NewOrderSellStopLimit(accountId, symbol, 1, 96168, 86170) // стоп-лимиткой
	orderState, err := client.OrdersService.PlaceOrder(ctx, newOrder)
	if err != nil {
		slog.Error("OrdersService", "PlaceOrder", err.Error())
	}
	slog.Info("OrdersService.PlaceOrder", "orderState", orderState)
}

// buyLimit пример лимитной покупки
func buyLimit(ctx context.Context, client *finam.Client, accountId string) {
	//symbol := "IMOEXF@RTSX"
	symbol := "SBER@MISX"
	newOrder := &orders_service.Order{
		AccountId:   accountId,
		Symbol:      symbol,
		Side:        side.Side_SIDE_BUY,
		Type:        orders_service.OrderType_ORDER_TYPE_LIMIT,
		Quantity:    finam.IntToDecimal(10),
		LimitPrice:  finam.Float64ToDecimal(294.65),
		TimeInForce: orders_service.TimeInForce_TIME_IN_FORCE_DAY,
		//TimeInForce: orders_service.TimeInForce_TIME_IN_FORCE_GOOD_TILL_CANCEL,
		//ClientOrderId: ""
	}

	orderState, err := client.OrdersService.PlaceOrder(ctx, newOrder)
	if err != nil {
		slog.Error("OrdersService", "PlaceOrder", err.Error())
	}
	slog.Info("OrdersService.PlaceOrder", "orderState", orderState)
}
