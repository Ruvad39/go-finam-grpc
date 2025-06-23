package finam

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"google.golang.org/grpc/metadata"
)

const authKey = "Authorization"             //
const jwtTokenTtl = 12 * time.Minute        // Время жизни токена JWT в минутах (15 минут)
const jwtRefreshInterval = 10 * time.Minute // Интервал обновления токена (в минутах)

// GetJWT Получение JWT токена из API токена
//
// идет  вызов AuthService.Auth
func (c *Client) GetJWT(ctx context.Context) (string, error) {
	if c.token == "" {
		c.accessToken = ""
		return c.accessToken, nil
	}
	req := &pb.AuthRequest{Secret: c.token}
	log.Debug("GetJWT start AuthService.Auth")
	t := time.Now()
	res, err := c.AuthService.Auth(ctx, req)
	if err != nil {
		return c.accessToken, err
	}
	log.Debug("GetJWT end AuthService.Auth", "duration", time.Since(t))
	return res.Token, nil
}

// UpdateJWT
// если jwt токен пустой или вышло его время
// получим JWT
// и запишем его в параметры клиента. Проставим время получения
func (c *Client) UpdateJWT(ctx context.Context) error {
	if c.token == "" {
		c.accessToken = ""
		return fmt.Errorf("UpdateJWT: token пустой")
	}
	// если токен пустой или вышло его время
	if c.accessToken == "" || c.ttlJWT.Before(time.Now()) {
		// получим новый токен
		log.Debug("UpdateJWT. токен пустой или вышло его время = получим новый токен")
		token, err := c.GetJWT(ctx)
		if err != nil {
			return err
		}
		// запишем время окончания токена
		c.ttlJWT = time.Now().Add(jwtTokenTtl)
		c.accessToken = token
	}
	log.Debug("UpdateJWT. токен живой")
	return nil
}

// runJwtRefresher в отдельном потоке периодически обновляем токен
func (c *Client) runJwtRefresher(ctx context.Context) {
	log.Debug("run JwtRefresher")
	if c.token == "" {
		c.accessToken = ""
		log.Warn("JWT refresher token пустой. Выход из метода")
		return
	}
	ticker := time.NewTicker(jwtRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug("start JwtRefresh")
			token, err := c.GetJWT(ctx)
			if err != nil {
				log.Error("failed to refresh JW", "err", err.Error())
			}
			// запишем время окончания токена
			c.ttlJWT = time.Now().Add(jwtTokenTtl)
			c.accessToken = token
		case <-ctx.Done():
			log.Debug("JWT refresher stopped")
			return
		}
	}
}

// GetTokenDetails Получение информации о токене сессии
//
// идет вызов AuthService.TokenDetails
func (c *Client) GetTokenDetails(ctx context.Context) (*pb.TokenDetailsResponse, error) {
	return c.AuthService.TokenDetails(ctx, &pb.TokenDetailsRequest{Token: c.accessToken})
}

// WithAuthToken Создаем новый контекст с заголовком Authorization
// пишем в него jwt токен
func (c *Client) WithAuthToken(ctx context.Context) (context.Context, error) {
	// проверим наличие токена
	err := c.UpdateJWT(ctx)
	if err != nil {
		return ctx, err
	}
	// добавим заголовок
	md := metadata.New(map[string]string{
		authKey: c.accessToken,
	})
	// и добавляем его в ctx
	return metadata.NewOutgoingContext(ctx, md), nil
}
