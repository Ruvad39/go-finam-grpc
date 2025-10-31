/*
AccountServiceClient = клиент для работы с AccountsService
методы:

// Получение информации по конкретному счету
GetAccount(ctx context.Context, accountId string)

// Получение истории по сделкам заданного счета
GetTrades(ctx context.Context, accountId string, start, end time.Time, limit int32)

// Получение списка транзакций по счету
GetTransactions(ctx context.Context, accountId string, start, end time.Time, limit int32)

*/

package finam

import (
	"context"
	"time"

	accounts_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/accounts"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AccountServiceClient клиент для работы с AccountsService
type AccountServiceClient struct {
	client          *Client
	AccountsService accounts_service.AccountsServiceClient
}

func NewAccountServiceClient(c *Client) *AccountServiceClient {
	return &AccountServiceClient{client: c,
		AccountsService: accounts_service.NewAccountsServiceClient(c.conn),
	}
}

// GetAccount Получение информации по конкретному аккаунту
func (s *AccountServiceClient) GetAccount(ctx context.Context, accountId string) (*accounts_service.GetAccountResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.AccountsService.GetAccount(ctx, &accounts_service.GetAccountRequest{AccountId: accountId})
}

// GetTrades Получение истории по сделкам аккаунта
func (s *AccountServiceClient) GetTrades(ctx context.Context, accountId string, start, end time.Time, limit int32) (*accounts_service.TradesResponse, error) {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.AccountsService.Trades(ctx, &accounts_service.TradesRequest{AccountId: accountId, Interval: i, Limit: limit})
}

// GetTransactions Получение списка транзакций аккаунта
func (s *AccountServiceClient) GetTransactions(ctx context.Context, accountId string, start, end time.Time, limit int32) (*accounts_service.TransactionsResponse, error) {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.AccountsService.Transactions(ctx, &accounts_service.TransactionsRequest{AccountId: accountId, Interval: i, Limit: limit})
}
