package finam

// options - параметры клиента
type options struct {
	EndPoint string // точка доступа к gRPC
}

// Option - тип функции для настройки
type Option func(*options)

// WithEndPoint точка доступа к gRPC
func WithEndPoint(value string) Option {
	return func(o *options) {
		o.EndPoint = value
	}
}
