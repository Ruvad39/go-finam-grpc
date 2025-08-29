/*
MarketDataServiceClient = клиент для работы с MarketDataService
методы:

//Получение исторических данных по инструменту (агрегированные свечи)
GetBars(ctx context.Context, symbol string, tf pb.TimeFrame, start, end time.Time)

// Получение последней котировки по инструменту
GetLastQuote(ctx context.Context, symbol string)

// получение списка последних сделок по инструменту
GetLatestTrades(ctx context.Context, symbol string)

// Получение текущего стакана по инструменту
 GetOrderBook(ctx context.Context, symbol string)

*/

package finam

import (
	"context"
	"fmt"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Bar аналог marketdata_service.Bar
// плюс доп. поля: Symbol + Timeframe
type Bar struct {
	// Инструмент
	Symbol string
	// Таймфрейм
	Timeframe marketdata_service.TimeFrame
	// Метка времени
	Timestamp *timestamppb.Timestamp
	// Цена открытия свечи
	Open *decimal.Decimal
	// Максимальная цена свечи
	High *decimal.Decimal
	// Минимальная цена свечи
	Low *decimal.Decimal
	// Цена закрытия свечи
	Close *decimal.Decimal
	// Объём торгов за свечу в шт.
	Volume *decimal.Decimal `protobuf:"bytes,6,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (b *Bar) GetOpen() float64 {
	return DecimalToFloat64(b.Open)
}
func (b *Bar) GetHigh() float64 {
	return DecimalToFloat64(b.High)
}
func (b *Bar) GetLow() float64 {
	return DecimalToFloat64(b.Low)
}
func (b *Bar) GetClose() float64 {
	return DecimalToFloat64(b.Close)
}
func (b *Bar) GetVolume() int {
	return DecimalToInt(b.Volume)
}
func (b *Bar) GetTime() time.Time {
	return b.Timestamp.AsTime().In(TzMoscow)
}

func (b *Bar) String() string {
	str := fmt.Sprintf("%v,%v,%v, Open:%v, High:%v, Low:%v, CLose:%v, Volume:%v",
		b.GetTime().String(), b.Symbol, b.Timeframe.String(), b.GetOpen(), b.GetHigh(), b.GetLow(), b.GetClose(), b.GetVolume())
	return str
}

type BarFunc func(bar *Bar)

// MarketDataServiceClient клиент для работы MarketDataService
type MarketDataServiceClient struct {
	client            *Client
	MarketDataService marketdata_service.MarketDataServiceClient
}

func NewMarketDataServiceClient(c *Client) *MarketDataServiceClient {
	return &MarketDataServiceClient{client: c,
		MarketDataService: marketdata_service.NewMarketDataServiceClient(c.conn),
	}
}

// GetLastQuote Получение последней котировки по инструменту
func (s *MarketDataServiceClient) GetLastQuote(ctx context.Context, symbol string) (*marketdata_service.QuoteResponse, error) {
	return s.MarketDataService.LastQuote(ctx, &marketdata_service.QuoteRequest{Symbol: symbol})
}

// GetBars Получение исторических данных по инструменту (агрегированные свечи)
// symbol string
// tf pb.TimeFrame
// start, end time.Time
func (s *MarketDataServiceClient) GetBars(ctx context.Context, symbol string, tf marketdata_service.TimeFrame, start, end time.Time) (*marketdata_service.BarsResponse, error) {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return s.MarketDataService.Bars(ctx, &marketdata_service.BarsRequest{Symbol: symbol, Timeframe: tf, Interval: i})
}

// GetOrderBook Получение текущего стакана по инструменту
func (s *MarketDataServiceClient) GetOrderBook(ctx context.Context, symbol string) (*marketdata_service.OrderBookResponse, error) {
	return s.MarketDataService.OrderBook(ctx, &marketdata_service.OrderBookRequest{Symbol: symbol})
}

// GetLatestTradesПолучение списка последних сделок по инструменту
func (s *MarketDataServiceClient) GetLatestTrades(ctx context.Context, symbol string) (*marketdata_service.LatestTradesResponse, error) {
	return s.MarketDataService.LatestTrades(ctx, &marketdata_service.LatestTradesRequest{Symbol: symbol})
}
