/*
TODO quoteChan сжелать буффиризированным. устанивить переменную для значение буыфера
*/

package finam

import (
	"context"
	"crypto/tls"
	"time"

	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	Name        = "FINAM-API-gRPC"
	Version     = "0.2.0"
	VersionDate = "2025-06-23"
)

// Endpoints
const (
	endPoint = "api.finam.ru:443" //"ftrr01.finam.ru:443"
)

const (
	initialDelay = 2 * time.Second
	maxDelay     = 50 * time.Second
)

// Client
type Client struct {
	opts        options   // Параметры клиента
	token       string    // Основой токен пользователя
	accessToken string    // JWT токен для дальнейшей авторизации
	ttlJWT      time.Time // Время завершения действия JWT токена
	conn        *grpc.ClientConn
	AuthService pb.AuthServiceClient
}

func NewClient(ctx context.Context, token string, opts ...Option) (*Client, error) {
	// Устанавливаем значения по умолчанию
	o := &options{
		EndPoint: endPoint,
	}
	// Применяем переданные опции
	for _, opt := range opts {
		opt(o)
	}
	//
	conn, err := grpc.NewClient(o.EndPoint,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                4 * time.Minute,  // отправлять ping каждые 4 минут
			Timeout:             30 * time.Second, // ждать ответа не дольше 30 сек
			PermitWithoutStream: true,             // пинговать даже без активных RPC
		}),
	)
	if err != nil {
		return nil, err
	}
	//
	client := &Client{
		opts:  *o,
		token: token,
		//accessToken: accountId,
		conn:        conn,
		AuthService: pb.NewAuthServiceClient(conn),
	}
	// сразу получим и запишем accessToken для работы
	err = client.UpdateJWT(ctx)
	if err != nil {
		return nil, err
	}
	// в отдельном потоке периодически обновляем accessToken
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
