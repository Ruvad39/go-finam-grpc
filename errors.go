package finam

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

// shouldTerminate
// по коду ошибки определим, нужно закрывать поток (stream)
func shouldTerminate(err error) bool {
	// закрытие потока = выходим
	if errors.Is(err, io.EOF) {
		return true
	}
	// Отмена через контекст = выходим
	if status.Code(err) == codes.Canceled || errors.Is(err, context.Canceled) {
		return true
	}
	// NotFound = идет когда не найден символ = выходим
	if status.Code(err) == codes.NotFound {
		return true
	}

	// ошибки сервера = переподключаемся
	if status.Code(err) == codes.Unavailable || status.Code(err) == codes.Internal {
		return false
	}

	return false
}
