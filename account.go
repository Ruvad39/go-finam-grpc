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
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// AccountServiceClient клиент для работы с AccountsService
type AccountServiceClient struct {
	client          *Client
	AccountsService pb.AccountsServiceClient
}

func NewAccountServiceClient(c *Client) *AccountServiceClient {
	return &AccountServiceClient{client: c,
		AccountsService: pb.NewAccountsServiceClient(c.conn),
	}
}

// GetAccount Получение информации по конкретному аккаунту
func (s *AccountServiceClient) GetAccount(ctx context.Context, accountId string) (*pb.GetAccountResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.AccountsService.GetAccount(ctx, &pb.GetAccountRequest{AccountId: accountId})
}

// GetTrades Получение истории по сделкам аккаунта
func (s *AccountServiceClient) GetTrades(ctx context.Context, accountId string, start, end time.Time, limit int32) (*pb.TradesResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return s.AccountsService.Trades(ctx, &pb.TradesRequest{AccountId: accountId, Interval: i, Limit: limit})
}

// GetTransactions Получение списка транзакций аккаунта
func (s *AccountServiceClient) GetTransactions(ctx context.Context, accountId string, start, end time.Time, limit int32) (*pb.TransactionsResponse, error) {
	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return s.AccountsService.Transactions(ctx, &pb.TransactionsRequest{AccountId: accountId, Interval: i, Limit: limit})
}
