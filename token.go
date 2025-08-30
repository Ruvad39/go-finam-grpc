/*
	Спасибо Артему за идею с TokenAgent
*/

package finam

import (
	"context"
	"sync"
	"time"

	auth_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/auth"
)

const authKey = "Authorization"             //
const jwtRefreshInterval = 12 * time.Minute // Интервал обновления токена (в минутах)

// TokenAgent
// соответствует интерфейсу PerRPCCredentials
type TokenAgent struct {
	jwt string // JWT токен для дальнейшей авторизации
	mu  sync.RWMutex
}

func NewTokenAgent() *TokenAgent {
	return &TokenAgent{}
}

func (a *TokenAgent) setJwt(jwt string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.jwt = jwt
}

func (a *TokenAgent) GetJwt() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.jwt
}

// GetRequestMetadata для соответствия интерфейсу PerRPCCredentials
func (a *TokenAgent) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	meta := make(map[string]string)
	meta[authKey] = a.GetJwt()

	return meta, nil
}

// GetRequestMetadata для соответствия интерфейсу PerRPCCredentials
func (a *TokenAgent) RequireTransportSecurity() bool {
	return true
}

// GetJWT Получение JWT токена из API токена
//
// идет  вызов AuthService.Auth
func (c *Client) GetJWT(ctx context.Context) (string, error) {
	if c.token == "" {
		c.TokenAgent.setJwt("")
		return "", nil
	}
	req := &auth_service.AuthRequest{Secret: c.token}
	log.Debug("GetJWT start AuthService.Auth")
	t := time.Now()
	res, err := c.AuthService.Auth(ctx, req)
	if err != nil {
		return "", err
	}
	log.Debug("GetJWT end AuthService.Auth", "duration", time.Since(t))
	return res.Token, nil
}

// JwtRefresh обновим токен
func (c *Client) JwtRefresh(ctx context.Context) error {
	token, err := c.GetJWT(ctx)
	if err != nil {
		log.Error("failed to refresh JW", "err", err.Error())
		return err
	}
	c.TokenAgent.setJwt(token)
	return nil
}

// runJwtRefresher в отдельном потоке периодически обновляем токен
func (c *Client) runJwtRefresher(ctx context.Context) {
	log.Debug("run JwtRefresher")
	if c.token == "" {
		c.TokenAgent.setJwt("")
		log.Warn("JWT refresher token пустой. Выход из метода")
		return
	}
	ticker := time.NewTicker(c.opts.jwtRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug("start JwtRefresh")
			_ = c.JwtRefresh(ctx)
		case <-ctx.Done():
			log.Debug("JWT refresher stopped")
			return
		}
	}
}

// GetTokenDetails Получение информации о токене сессии
//
// идет вызов AuthService.TokenDetails
func (c *Client) GetTokenDetails(ctx context.Context) (*auth_service.TokenDetailsResponse, error) {
	return c.AuthService.TokenDetails(ctx, &auth_service.TokenDetailsRequest{Token: c.TokenAgent.GetJwt()})
}
