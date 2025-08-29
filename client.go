/*
TODO quoteChan сжелать буффиризированным. устанивить переменную для значение буыфера
*/

package finam

import (
	"context"
	"crypto/tls"
	"time"

	auth_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	Name        = "FINAM-API-gRPC"
	Version     = "0.4.0"
	VersionDate = "2025-08-29"
)

// Endpoints
const (
	endPoint = "api.finam.ru:443" // "ftrr01.finam.ru:443"
)

const (
	initialDelay     = 2 * time.Second
	maxDelay         = 50 * time.Second
	keepaliveTime    = 4 * time.Minute  // Отправлять ping каждые 4 минут
	keepaliveTimeout = 30 * time.Second // Ждать ответа не дольше 30 сек
)

// Client
type Client struct {
	opts        clientOptions // Параметры клиента
	TokenAgent  *TokenAgent
	token       string // Основой токен пользователя
	conn        *grpc.ClientConn
	AuthService auth_service.AuthServiceClient
}

func NewClient(ctx context.Context, token string, opts ...ClientOption) (*Client, error) {
	// Устанавливаем значения по умолчанию
	o := ClientOptionsDefault()
	// Применяем переданные опции
	for _, opt := range opts {
		opt(o)
	}

	tokenAgent := NewTokenAgent()
	conn, err := grpc.NewClient(o.EndPoint,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                o.keepaliveTime,    // отправлять ping каждые 4 минут
			Timeout:             o.keepaliveTimeout, // ждать ответа не дольше 30 сек
			PermitWithoutStream: true,               // пинговать даже без активных RPC
		}),
		grpc.WithPerRPCCredentials(tokenAgent),
	)
	if err != nil {
		return nil, err
	}

	// создание клиента
	client := &Client{
		opts:        *o,
		TokenAgent:  tokenAgent,
		token:       token,
		conn:        conn,
		AuthService: auth_service.NewAuthServiceClient(conn),
	}

	// сразу получим и запишем jwt для работы
	err = client.JwtRefresh(ctx)
	if err != nil {
		return nil, err
	}

	// в отдельном потоке периодически обновляем Jwt
	go client.runJwtRefresher(ctx)
	return client, nil
}

func (c *Client) Close() error {
	return c.conn.Close()

}

// NewAccountServiceClien создаем клиент для доступа к AccountService
func (c *Client) NewAccountServiceClient() *AccountServiceClient {
	return NewAccountServiceClient(c)
}

// NewAssetServiceClient создаем клиент для доступа к AssetService
func (c *Client) NewAssetServiceClient() *AssetServiceClient {
	return NewAssetServiceClient(c)
}

// NewMarketDataServiceClient создаем клиент для работы MarketDataService
func (c *Client) NewMarketDataServiceClient() *MarketDataServiceClient {
	return NewMarketDataServiceClient(c)
}

// NewOrderServiceClient создаем клиент для работы OrdersService
func (c *Client) NewOrderServiceClient() *OrderServiceClient {
	return NewOrderServiceClient(c)
}
