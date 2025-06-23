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
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// MarketDataServiceClient клиент для работы MarketDataService
type MarketDataServiceClient struct {
	client            *Client
	MarketDataService pb.MarketDataServiceClient
}

func NewMarketDataServiceClient(c *Client) *MarketDataServiceClient {
	return &MarketDataServiceClient{client: c,
		MarketDataService: pb.NewMarketDataServiceClient(c.conn),
	}
}

// GetLastQuote Получение последней котировки по инструменту
func (s *MarketDataServiceClient) GetLastQuote(ctx context.Context, symbol string) (*pb.QuoteResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.MarketDataService.LastQuote(ctx, &pb.QuoteRequest{Symbol: symbol})
}

// GetBars Получение исторических данных по инструменту (агрегированные свечи)
// symbol string
// tf pb.TimeFrame
// start, end time.Time
func (s *MarketDataServiceClient) GetBars(ctx context.Context, symbol string, tf pb.TimeFrame, start, end time.Time) (*pb.BarsResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return s.MarketDataService.Bars(ctx, &pb.BarsRequest{Symbol: symbol, Timeframe: tf, Interval: i})
}

// GetOrderBook Получение текущего стакана по инструменту
func (s *MarketDataServiceClient) GetOrderBook(ctx context.Context, symbol string) (*pb.OrderBookResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.MarketDataService.OrderBook(ctx, &pb.OrderBookRequest{Symbol: symbol})
}

// GetLatestTradesПолучение списка последних сделок по инструменту
func (s *MarketDataServiceClient) GetLatestTrades(ctx context.Context, symbol string) (*pb.LatestTradesResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.MarketDataService.LatestTrades(ctx, &pb.LatestTradesRequest{Symbol: symbol})
}
