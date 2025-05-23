/*
TODO quoteChan сжелать буффиризированным. устанивить переменную для значение буыфера
*/

package finam

import (
	"context"
	"crypto/tls"
	accounts_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/accounts"
	assets_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/assets"
	auth_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/auth"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	orders_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
	"time"
)

// Endpoints
const (
	endPoint = "ftrr01.finam.ru:443"
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
	token             string    // Основой токен пользователя
	accessToken       string    // JWT токен для дальнейшей авторизации
	ttlJWT            time.Time // Время завершения действия JWT токена
	conn              *grpc.ClientConn
	AuthService       auth_service.AuthServiceClient
	AccountsService   accounts_service.AccountsServiceClient
	AssetsService     assets_service.AssetsServiceClient
	MarketDataService marketdata_service.MarketDataServiceClient
	OrdersService     orders_service.OrdersServiceClient
	Securities        map[string]Security //  Список инструментов с которыми работаем (или весь список? )
}

func NewClient(ctx context.Context, token string) (*Client, error) {
	tlsConfig := tls.Config{MinVersion: tls.VersionTLS12}
	// TODO выделить в отдельный метод connect()
	log.Debug("NewClient start connect")
	conn, err := grpc.NewClient(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
	if err != nil {
		return nil, err
	}
	client := &Client{
		token:             token,
		conn:              conn,
		AuthService:       auth_service.NewAuthServiceClient(conn),
		AccountsService:   accounts_service.NewAccountsServiceClient(conn),
		AssetsService:     assets_service.NewAssetsServiceClient(conn),
		MarketDataService: marketdata_service.NewMarketDataServiceClient(conn),
		OrdersService:     orders_service.NewOrdersServiceClient(conn),
		Securities:        make(map[string]Security),
	}
	log.Debug("NewClient есть connect")
	err = client.UpdateJWT(ctx) // сразу получим и запишем токен для работы
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() error {
	//close(c.closeChan)
	//close(c.quoteChan)
	//close(c.rawQuoteChan)
	//close(c.errChan)
	return c.conn.Close()

}
