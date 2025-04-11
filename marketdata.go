package finam

import (
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// SubscribeQuote

// NewBarsRequest
func NewBarsRequest(symbol string, tf marketdata_service.BarsRequest_TimeFrame, start, end time.Time) *marketdata_service.BarsRequest {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return &marketdata_service.BarsRequest{
		Symbol:    symbol,
		Timeframe: tf,
		Interval:  i,
	}
}

// NewQuoteRequest
func NewQuoteRequest(symbol string) *marketdata_service.QuoteRequest {
	return &marketdata_service.QuoteRequest{
		Symbol: symbol,
	}
}

// NewOrderBookRequest
func NewOrderBookRequest(symbol string) *marketdata_service.OrderBookRequest {
	return &marketdata_service.OrderBookRequest{
		Symbol: symbol,
	}
}

// NewLatestTradesRequest
func NewLatestTradesRequest(symbol string) *marketdata_service.LatestTradesRequest {
	return &marketdata_service.LatestTradesRequest{
		Symbol: symbol,
	}
}

// NewSubscribeQuoteRequest
func NewSubscribeQuoteRequest(symbols []string) *marketdata_service.SubscribeQuoteRequest {
	return &marketdata_service.SubscribeQuoteRequest{Symbols: symbols}
}
