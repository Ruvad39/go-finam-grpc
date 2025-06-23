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
	Name        = "FINAM-API-gRPC GO"
	Version     = "0.1.1"
	VersionDate = "2025-04-22"
)

// Endpoints
const (
	endPoint = "api.finam.ru:443" //"ftrr01.finam.ru:443"
)

// Client
type Client struct {
	opts        options   // Параметры клиента
	token       string    // Основой токен пользователя
	accessToken string    // JWT токен для дальнейшей авторизации
	ttlJWT      time.Time // Время завершения действия JWT токена
	AccountId   string    // Код счета по умолчанию
	conn        *grpc.ClientConn
	AuthService pb.AuthServiceClient
}

func NewClient(ctx context.Context, token, accountId string, opts ...Option) (*Client, error) {
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
			Time:                15 * time.Minute, // отправлять ping каждые 15 минут
			Timeout:             30 * time.Second, // ждать ответа не дольше 10 сек
			PermitWithoutStream: true,             // пинговать даже без активных RPC
		}),
	)
	if err != nil {
		return nil, err
	}
	//
	client := &Client{
		opts:        *o,
		token:       token,
		accessToken: accountId,
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
