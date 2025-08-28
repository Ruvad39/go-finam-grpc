/*
 OrderServiceClient = клиент для работы OrdersService
методы:

// Получение списка заявок по заданному счету
GetOrders(ctx context.Context, accountId string) (*pb.OrdersResponse, error)

// Получение информации о конкретном ордере
GetOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error)

// Отмена биржевой заявки
CancelOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error)

// Выставление биржевой заявки
PlaceOrder(ctx context.Context, order *pb.Order) (*pb.OrderState, error)

// Вспомогательные методы для создания ордера
// создать ордер на покупку по рынку
NewBuyOrder(accountId, symbol string, quantity int) *pb.Order {

// создать ордер на продажу по рынку
NewSellOrder(accountId, symbol string, quantity int) *pb.Order {

// создать ордер на покупку по лимитной цене
NewBuyLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order

// создать ордер на продажу по лимитной цене
NewSellLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order

*/

package finam

import (
	"context"
	v1 "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1"
	pb "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/orders"
)

// OrderServiceClient клиент для работы OrdersService
type OrderServiceClient struct {
	client       *Client
	OrderService pb.OrdersServiceClient
}

// создаем клиент для работы OrdersService
func NewOrderServiceClient(c *Client) *OrderServiceClient {
	return &OrderServiceClient{client: c,
		OrderService: pb.NewOrdersServiceClient(c.conn),
	}
}

// GetOrders Получение списка заявок по заданному счету
func (s *OrderServiceClient) GetOrders(ctx context.Context, accountId string) (*pb.OrdersResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.OrderService.GetOrders(ctx, &pb.OrdersRequest{AccountId: accountId})
}

// GetOrder Получение информации о конкретном ордере
func (s *OrderServiceClient) GetOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.OrderService.GetOrder(ctx, &pb.GetOrderRequest{AccountId: accountId, OrderId: orderId})
}

// CancelOrder Отмена биржевой заявки
func (s *OrderServiceClient) CancelOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.OrderService.CancelOrder(ctx, &pb.CancelOrderRequest{AccountId: accountId, OrderId: orderId})
}

// PlaceOrder Выставление биржевой заявки
//
// на входе
// order *pb.Order = заполненный ордер
func (s *OrderServiceClient) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.OrderState, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.OrderService.PlaceOrder(ctx, order)

}

// NewBuyOrder создать ордер на покупку по рынку
// accountId string
// symbol string
// quantity int = ко-во в штуках
func (s *OrderServiceClient) NewBuyOrder(accountId, symbol string, quantity int) *pb.Order {
	return &pb.Order{
		Type:      pb.OrderType_ORDER_TYPE_MARKET, // по рынку
		Side:      v1.Side_SIDE_BUY,               // покупка
		AccountId: accountId,
		Symbol:    symbol,
		Quantity:  IntToDecimal(quantity),
	}
}

// NewSellOrder создать ордер на продажу по рынку
// accountId string
// symbol string
// quantity int = ко-во в штуках
func (s *OrderServiceClient) NewSellOrder(accountId, symbol string, quantity int) *pb.Order {
	return &pb.Order{
		Type:      pb.OrderType_ORDER_TYPE_MARKET, // по рынку
		Side:      v1.Side_SIDE_SELL,              // продажа
		AccountId: accountId,
		Symbol:    symbol,
		Quantity:  IntToDecimal(quantity),
	}
}

// NewBuyLimitOrder создать ордер на покупку по лимитной цене
// accountId string
// symbol string
// quantity int = ко-во в штуках
// price float6 = по какой цене выставить заявку
func (s *OrderServiceClient) NewBuyLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order {
	return &pb.Order{
		Type:       pb.OrderType_ORDER_TYPE_LIMIT, // лимитная
		Side:       v1.Side_SIDE_BUY,              // покупка
		AccountId:  accountId,
		Symbol:     symbol,
		Quantity:   IntToDecimal(quantity),
		LimitPrice: Float64ToDecimal(price),
	}
}

// NewSellLimitOrder создать ордер на продажу по лимитной цене
// accountId string
// symbol string
// quantity int = ко-во в штуках
// price float6 = по какой цене выставить заявку
func (s *OrderServiceClient) NewSellLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order {
	return &pb.Order{
		Type:       pb.OrderType_ORDER_TYPE_LIMIT, // лимитная
		Side:       v1.Side_SIDE_SELL,             // покупка
		AccountId:  accountId,
		Symbol:     symbol,
		Quantity:   IntToDecimal(quantity),
		LimitPrice: Float64ToDecimal(price),
	}
}
