package finam

import (
	"log/slog"
	"time"
)

// clientOptions - параметры клиента
type clientOptions struct {
	EndPoint           string        // точка доступа к gRPC
	jwtRefreshInterval time.Duration // С какой периодичностью обновлять JWT
	keepaliveTime      time.Duration // С какой периодичностью отправлять ping
	keepaliveTimeout   time.Duration // Сколько ждать ответа
	callTimeout        time.Duration // Timeout при вызове методов grpc
}

// ClientOption - тип функции для настройки
type ClientOption func(*clientOptions)

func ClientOptionsDefault() *clientOptions {
	return &clientOptions{
		EndPoint:           endPoint,
		jwtRefreshInterval: jwtRefreshInterval,
		keepaliveTime:      keepaliveTime,
		keepaliveTimeout:   keepaliveTimeout,
		callTimeout:        callTimeout,
	}
}

// WithEndPoint точка доступа к gRPC
func WithEndPoint(value string) ClientOption {
	return func(o *clientOptions) {
		o.EndPoint = value
	}
}

// WithJwtRefreshInterval установить время обновления JWT
func WithJwtRefreshInterval(value time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.jwtRefreshInterval = value
	}
}

func WithKeepaliveTime(value time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.keepaliveTime = value
	}
}

// WithCallTimeout
// установить Timeout для вызова методов
func WithCallTimeout(value time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.callTimeout = value
	}
}

func WithKeepaliveTimeout(value time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.keepaliveTimeout = value
	}
}

// WithLogger установить logger
func WithLogger(logger *slog.Logger) ClientOption {
	log = logger
	return func(o *clientOptions) {

	}
}
