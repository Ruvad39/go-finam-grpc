package finam

import (
	side "github.com/Ruvad39/go-finam-grpc/trade_api/v1"
	orders_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/orders"
)

// NewOrdersRequest
func NewOrdersRequest(accountId string) *orders_service.OrdersRequest {
	return &orders_service.OrdersRequest{AccountId: accountId}
}

// NewGetOrderRequest
func NewGetOrderRequest(accountId, orderId string) *orders_service.GetOrderRequest {
	return &orders_service.GetOrderRequest{AccountId: accountId, OrderId: orderId}
}

// NewCancelOrderReques
func NewCancelOrderRequest(accountId, orderId string) *orders_service.CancelOrderRequest {
	return &orders_service.CancelOrderRequest{AccountId: accountId, OrderId: orderId}
}

// NewOrderBuy сформируем ордер на покупку по рынку
// accountId string Идентификатор аккаунта
// symbol string  Символ инструмента
// quantity Количество в шт.
func NewOrderBuy(accountId, symbol string, quantity int) *orders_service.Order {
	return &orders_service.Order{
		Side:      side.Side_SIDE_BUY,
		Type:      orders_service.OrderType_ORDER_TYPE_MARKET,
		AccountId: accountId,
		Symbol:    symbol,
		Quantity:  IntToDecimal(quantity),
	}
}

// NewOrderSell сформируем ордер на продажу по рынку
// accountId string Идентификатор аккаунта
// symbol string  Символ инструмента
// quantity Количество в шт.
func NewOrderSell(accountId, symbol string, quantity int) *orders_service.Order {
	return &orders_service.Order{
		Side:      side.Side_SIDE_SELL,
		Type:      orders_service.OrderType_ORDER_TYPE_MARKET,
		AccountId: accountId,
		Symbol:    symbol,
		Quantity:  IntToDecimal(quantity),
	}
}

// NewOrderBuyLimut сформируем ордер на покупку по лимитной цене
// accountId string Идентификатор аккаунта
// symbol string  Символ инструмента
// quantity Количество в шт.
// price float64 по какой цене заявка
func NewOrderBuyLimit(accountId, symbol string, quantity int, price float64) *orders_service.Order {
	return &orders_service.Order{
		Side:       side.Side_SIDE_BUY,
		Type:       orders_service.OrderType_ORDER_TYPE_LIMIT,
		AccountId:  accountId,
		Symbol:     symbol,
		Quantity:   IntToDecimal(quantity),
		LimitPrice: Float64ToDecimal(price),
	}
}

// NewOrderSellLimit сформируем ордер на продажу по лимитной цене
// accountId string Идентификатор аккаунта
// symbol string  Символ инструмента
// quantity Количество в шт.
// price float64 по какой цене заявка
func NewOrderSellLimit(accountId, symbol string, quantity int, price float64) *orders_service.Order {
	return &orders_service.Order{
		Side:       side.Side_SIDE_SELL,
		Type:       orders_service.OrderType_ORDER_TYPE_LIMIT,
		AccountId:  accountId,
		Symbol:     symbol,
		Quantity:   IntToDecimal(quantity),
		LimitPrice: Float64ToDecimal(price),
	}
}
