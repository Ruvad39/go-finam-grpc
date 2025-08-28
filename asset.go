/*
AssetServiceClient = клиент для работы с AssetsService
методы:

//  вернем текущее время сервера (в TzMoscow)
GetTime(ctx context.Context)

//Получение списка доступных бирж, названия и mic коды
GetExchanges(ctx context.Context)

// Получение списка доступных инструментов, их описание
GetAssets(ctx context.Context)

// Получение информации по конкретному инструменту
GetAsset(ctx context.Context, accountId, symbol string)

// Получение торговых параметров по инструменту
GetAssetParams(ctx context.Context, accountId, symbol string)

// Получение расписания торгов для инструмента
GetSchedule(ctx context.Context, symbol string)

TODO OptionsChain


*/

package finam

import (
	"context"
	assets_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/assets"
	"time"
)

// AssetServiceClient клиент для работы AssetsService
type AssetServiceClient struct {
	client       *Client
	AssetService assets_service.AssetsServiceClient
}

func NewAssetServiceClient(c *Client) *AssetServiceClient {
	return &AssetServiceClient{client: c,
		AssetService: assets_service.NewAssetsServiceClient(c.conn),
	}
}

// GetTime вернем текущее время сервера (в TzMoscow)
func (s *AssetServiceClient) GetTime(ctx context.Context) (time.Time, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return time.Time{}, err
	}
	resp, err := s.AssetService.Clock(ctx, &assets_service.ClockRequest{})
	if err != nil {
		return time.Time{}, err
	}
	return resp.Timestamp.AsTime().In(TzMoscow), err
}

// GetExchanges Получение списка доступных бирж, названия и mic коды
func (s *AssetServiceClient) GetExchanges(ctx context.Context) (*assets_service.ExchangesResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AssetService.Exchanges(ctx, &assets_service.ExchangesRequest{})
}

// GetAssets Получение списка доступных инструментов, их описание
func (s *AssetServiceClient) GetAssets(ctx context.Context) (*assets_service.AssetsResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AssetService.Assets(ctx, &assets_service.AssetsRequest{})
}

// GetAsset Получение информации по конкретному инструменту
func (s *AssetServiceClient) GetAsset(ctx context.Context, accountId, symbol string) (*assets_service.GetAssetResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AssetService.GetAsset(ctx, &assets_service.GetAssetRequest{AccountId: accountId, Symbol: symbol})
}

// GetAssetParams Получение торговых параметров по инструменту
func (s *AssetServiceClient) GetAssetParams(ctx context.Context, accountId, symbol string) (*assets_service.GetAssetParamsResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AssetService.GetAssetParams(ctx, &assets_service.GetAssetParamsRequest{AccountId: accountId, Symbol: symbol})
}

// GetSchedule Получение расписания торгов для инструмента
func (s *AssetServiceClient) GetSchedule(ctx context.Context, symbol string) (*assets_service.ScheduleResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AssetService.Schedule(ctx, &assets_service.ScheduleRequest{Symbol: symbol})
}
