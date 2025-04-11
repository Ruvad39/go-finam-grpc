/*
TODO создать метод обновляющий токен каждые 20 минут в отдельном потоке
*/

package finam

import (
	"context"
	"fmt"
	auth_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/auth"
	"google.golang.org/grpc/metadata"
	"time"
)

const jwtTokenTtl = 12 // Время жизни токена JWT в минутах (15 минут)

// GetJWT Получение JWT токена из API токена
//
// идет  вызов AuthService.Auth
func (c *Client) GetJWT(ctx context.Context) (string, error) {
	if c.token == "" {
		c.accessToken = ""
		return c.accessToken, nil
	}
	req := &auth_service.AuthRequest{Secret: c.token}
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
		c.ttlJWT = time.Now().Add(jwtTokenTtl * time.Minute)
		c.accessToken = token
	}
	log.Debug("UpdateJWT. токен живой")
	return nil
}

// GetTokenDetails Получение информации о токене сессии
//
// идет вызов AuthService.TokenDetails
func (c *Client) GetTokenDetails(ctx context.Context) (*auth_service.TokenDetailsResponse, error) {
	return c.AuthService.TokenDetails(ctx, &auth_service.TokenDetailsRequest{Token: c.accessToken})
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
		"Authorization": c.accessToken,
	})
	return metadata.NewOutgoingContext(ctx, md), nil
}
