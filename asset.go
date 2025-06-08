package finam

import (
	"context"
	assets_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/assets"
	"time"
)

// GetTime вернем текущее время сервера (в TzMoscow)
func (c *Client) GetTime(ctx context.Context) (time.Time, error) {
	resp, err := c.AssetsService.Clock(ctx, &assets_service.ClockRequest{})
	if err != nil {
		return time.Time{}, err
	}
	return resp.Timestamp.AsTime().In(TzMoscow), err
}

func NewClockRequest() *assets_service.ClockRequest {
	return &assets_service.ClockRequest{}
}

func NewExchangesRequest() *assets_service.ExchangesRequest {
	return &assets_service.ExchangesRequest{}
}

func NewAssetsRequest() *assets_service.AssetsRequest {
	return &assets_service.AssetsRequest{}
}

func NewAssetParamsRequest(symbol, accountId string) *assets_service.GetAssetParamsRequest {
	return &assets_service.GetAssetParamsRequest{Symbol: symbol, AccountId: accountId}
}

func NewAssetRequest(symbol, accountId string) *assets_service.GetAssetRequest {
	return &assets_service.GetAssetRequest{Symbol: symbol, AccountId: accountId}
}

func NewOptionsChainRequest(symbol string) *assets_service.OptionsChainRequest {
	return &assets_service.OptionsChainRequest{UnderlyingSymbol: symbol}
}

func NewScheduleRequest(symbol string) *assets_service.ScheduleRequest {
	return &assets_service.ScheduleRequest{Symbol: symbol}
}

//// OptionsChain Получение цепочки опционов для базового актива
//// symbol Символ базового актива опциона
//func (c *Client) GeOptionsChain(ctx context.Context, symbol string) (*assets_service.OptionsChainResponse, error) {
//	// добавим заголовок с авторизацией
//	ctx, err := c.WithAuthToken(ctx)
//	if err != nil {
//		return nil, err
//	}
//	req := &assets_service.OptionsChainRequest{UnderlyingSymbol: symbol}
//	return c.AssetsService.OptionsChain(ctx, req)
//}

// GeSchedule Получение расписания торгов для инструмента
// symbol Символ инструмента
//func (c *Client) GeSchedule(ctx context.Context, symbol string) (*assets_service.ScheduleResponse, error) {
//	// добавим заголовок с авторизацией
//	ctx, err := c.WithAuthToken(ctx)
//	if err != nil {
//		return nil, err
//	}
//	req := &assets_service.ScheduleRequest{Symbol: symbol}
//	return c.AssetsService.Schedule(ctx, req)
//}
