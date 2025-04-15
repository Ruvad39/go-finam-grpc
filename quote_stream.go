/*
разделил на 2 канала с данынми

	rawQuoteChan      chan *marketdata_service.Quote // Канал с "сырыми" данными по котировкам
	quoteChan         chan Quote // Канал с обработанными котировками
*/
package finam

import (
	"context"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
)

type QuoteFunc func(quote Quote)
type RawQuoteFunc func(quote *marketdata_service.Quote)

// SetQuoteHandler
// установим функцию для обработки поступления котировок
// должен устанавливать раньше Start
func (s *Stream) SetQuoteHandler(handler QuoteFunc) {
	s.handleQuote = handler
}

// SetRawQuoteHandler
// установим функцию для обработки поступления "сырых" котировок
// должен устанавливать раньше Start
func (s *Stream) SetRawQuoteHandler(handler RawQuoteFunc) {
	s.handleRawQuote = handler
}

// GetRawQuoteChan вернем канал в котором сырые данные по котировкам
//func (s *Stream) GetRawQuoteChan() chan *marketdata_service.Quote {
//	return s.rawQuoteChan
//}

// startHandleQuoteWorker Воркер для последовательной обработки канала котировок
func (s *Stream) startHandleQuoteWorker(ctx context.Context) {
	log.Debug("startHandleQuoteWorker")
	// "сырые" котировки
	go func() {
		hasHandleRawQuote := s.handleRawQuote != nil
		hasHandleQuote := s.handleQuote != nil
		var err error
		var processedQuote Quote
		var rawQuote *marketdata_service.Quote
		var ok bool
		for {
			select {
			case <-ctx.Done():
				log.Debug("HandleRawQuoteWorker ctx.Done() = выход")
				return
			case rawQuote, ok = <-s.rawQuoteChan:
				if !ok {
					log.Debug("HandleRawQuoteWorker quoteChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				if hasHandleRawQuote {
					s.handleRawQuote(rawQuote)
				}
				// обработаем инкрементальные данные в срез
				// только если указана функция приема
				if hasHandleQuote {
					// Обработка котировки
					processedQuote.Reset()
					processedQuote, err = s.quoteStore.processQuote(rawQuote) // Обработаем сырые данные. Вернем Quote
					if err != nil {
						log.Error("ошибка обработки котировки", slog.Any("rawQuote", rawQuote), slog.String("err", err.Error()))
						continue
					}
					// пошлем получателю
					s.handleQuote(processedQuote)
				}
				// очистим
				rawQuote.Reset()
			}
		}
	}()

}

// startQuoteStream
// собираем данные с подписки QuoteChannel
// делаем подписку (SubscribeQuote)
// запускаем в отдельном потоке метод для прослушивания стрима (listenQuoteStream)
func (s *Stream) startQuoteStream(ctx context.Context) error {
	log.Debug("start startQuoteStream")
	symbols := s.getSymbolsByChannel(QuoteChannel)
	// если список пустой = выйдем
	if len(symbols) == 0 {
		return nil // errors.New("no quote symbols found")
	}
	// есть список инструментов для подписки
	var err error
	var stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]
	log.Debug("StartStream", "QuoteSymbols", symbols)
	// добавим заголовок с авторизацией (accessToken)
	ctx, err = s.client.WithAuthToken(ctx)
	if err != nil {
		return err
	}
	// команда на подписку данных
	stream, err = s.client.MarketDataService.SubscribeQuote(ctx, NewSubscribeQuoteRequest(symbols))
	if err != nil {
		return err
	}

	// в отдельном потоке запустим чтения данных из стрима
	go s.listenQuoteStream(ctx, stream)

	return err
}

// listenQuoteStream чтение данных из стрима котировок
func (s *Stream) listenQuoteStream(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]) {
	log.Debug("start listenQuoteStream")
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeQuoteResponse
	var quoteSlice []*marketdata_service.Quote

	// читаем поток
	for {
		select {
		case <-ctx.Done():
			log.Debug("listenQuoteStream: ctx.Done(), выходим")
			return
		case <-s.closeChan:
			log.Debug("listenQuoteStream: closeChan, выходим")
			return
		default:
			msg, err = stream.Recv()
			if err != nil {
				// Проверка на конкретный код ошибки
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					log.Debug("listenQuoteStream: Recv: контекст отменён (canceled), выходим")
					return
				}
				if err == io.EOF {
					log.Debug("listenQuoteStream: Поток завершён")
					return // Поток завершён, выходим из цикла
				} else {
					log.Error("listenQuoteStream:  Ошибка чтения из потока", "err", err.Error())
					//s.errChan <- err
					s.Reconnect(QuoteChannel)
					return //  выход
				}
			}
			// В потоке приходит массив данных
			quoteSlice = quoteSlice[:0] // предварительно очистим
			quoteSlice = msg.GetQuote()
			if len(quoteSlice) != 0 {
				for _, rawQuote := range quoteSlice {
					select {
					case s.rawQuoteChan <- rawQuote:
						// успешно отправили
					default:
						log.Error("QuoteStream: канал переполнен, дропаем", slog.Any("rawQuote", rawQuote))
					}
				}
			}
			// очистим переменные
			msg.Reset()
		}
	}
}
