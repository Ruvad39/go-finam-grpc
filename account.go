package finam

import (
	"context"
	accounts_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/accounts"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

// AccountRequest Получение Информация о конкретном аккаунте
type AccountRequest struct {
	client         *Client
	accountRequest *accounts_service.GetAccountRequest
}

// NewAccountRequest
func (c *Client) NewAccountRequest(accountId string) *AccountRequest {
	return &AccountRequest{
		client:         c,
		accountRequest: &accounts_service.GetAccountRequest{AccountId: accountId},
	}

}

// Do выполним запрос AccountsService.GetAccount()
func (r *AccountRequest) Do(ctx context.Context) (*accounts_service.GetAccountResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := r.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return r.client.AccountsService.GetAccount(ctx, r.accountRequest)
}

func NewGetAccountRequest(accountId string) *accounts_service.GetAccountRequest {
	return &accounts_service.GetAccountRequest{AccountId: accountId}

}

// AccountTradesRequest Получение истории по сделкам аккаунта
type AccountTradesRequest struct {
	client    *Client
	accountId string    // Идентификатор аккаунта
	limit     int32     // Лимит количества сделок
	startTime time.Time // Начало запрашиваемого периода
	endTime   time.Time // Окончание запрашиваемого периода
}

func (c *Client) NewAccountTradesRequest(accountId string) *AccountTradesRequest {
	return &AccountTradesRequest{
		client:    c,
		accountId: accountId,
	}
}

func (r *AccountTradesRequest) Limit(value int32) *AccountTradesRequest {
	r.limit = value
	return r
}

func (r *AccountTradesRequest) StartTime(value time.Time) *AccountTradesRequest {
	r.startTime = value
	return r
}

func (r *AccountTradesRequest) EndTime(value time.Time) *AccountTradesRequest {
	r.endTime = value
	return r
}

// Do Получение истории по сделкам аккаунта
func (r *AccountTradesRequest) Do(ctx context.Context) (*accounts_service.TradesResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := r.client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("main", "WithAuthToken", err.Error())
		// если прошла ошибка, дальше работа бесполезна, не будет авторизации
		return nil, err
	}
	inv := &interval.Interval{
		StartTime: timestamppb.New(r.startTime),
		EndTime:   timestamppb.New(r.endTime),
	}
	req := &accounts_service.TradesRequest{
		AccountId: r.accountId,
		Limit:     r.limit,
		Interval:  inv,
	}
	return r.client.AccountsService.Trades(ctx, req)
}

// NewTradesRequest
func NewTradesRequest(accountId string, limit int32, start, end time.Time) *accounts_service.TradesRequest {
	inv := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	result := &accounts_service.TradesRequest{
		AccountId: accountId,
		Limit:     limit,
		Interval:  inv,
	}
	return result
}

// NewTransactionsRequest
func NewTransactionsRequest(accountId string, limit int32, start, end time.Time) *accounts_service.TransactionsRequest {
	inv := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	result := &accounts_service.TransactionsRequest{
		AccountId: accountId,
		Limit:     limit,
		Interval:  inv,
	}
	return result
}

// GetTransactions Получение списка транзакций аккаунта
// accountId	string	Идентификатор аккаунта
// limit	int32	Лимит количества сделок
// start, end time.Time Начало и окончание запрашиваемого периода
//func (c *Client) GetTransactions(ctx context.Context, accountId string, limit int32, start, end time.Time) (*accounts_service.TransactionsResponse, error) {
//	// добавим заголовок с авторизацией
//	ctx, err := c.WithAuthToken(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return c.AccountsService.Transactions(ctx, NewTransactionsRequest(accountId, limit, start, end))
//}
