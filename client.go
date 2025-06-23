/*
TODO quoteChan сжелать буффиризированным. устанивить переменную для значение буыфера
*/

package finam

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"time"

	//accounts_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/accounts"
	//assets_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/assets"
	//auth_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/auth"
	//marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	//orders_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/orders"
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

var logLevel = &slog.LevelVar{} // INFO
var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: logLevel,
})).With(slog.String("package", "go-finam-grpc"))

func SetLogger(logger *slog.Logger) {
	log = logger
}

// SetLogDebug установим уровень логгирования Debug
func SetLogDebug(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
}

// Client
type Client struct {
	token       string    // Основой токен пользователя
	accessToken string    // JWT токен для дальнейшей авторизации
	ttlJWT      time.Time // Время завершения действия JWT токена
	conn        *grpc.ClientConn
	//AuthService       auth_service.AuthServiceClient

}

func NewClient(ctx context.Context, token string) (*Client, error) {
	// TODO выделить в отдельный метод connect()
	log.Debug("NewClient start connect")
	conn, err := grpc.NewClient(endPoint,
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
	client := &Client{
		token: token,
		conn:  conn,
		//AuthService:       auth_service.NewAuthServiceClient(conn),

	}
	log.Debug("NewClient есть connect")
	err = client.UpdateJWT(ctx) // сразу получим и запишем accessToken для работы
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
