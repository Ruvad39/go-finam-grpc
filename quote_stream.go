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

// SetQuoteHandler
// установим функцию для обработки поступления котировок
// должен устанавливать раньше Start
func (s *Stream) SetQuoteHandler(handler QuoteFunc) {
	s.handleQuote = handler
}

// GetRawQuoteChan вернем канал в котором сырые данные по котировкам
func (s *Stream) GetRawQuoteChan() chan *marketdata_service.Quote {
	return s.rawQuoteChan
}

// startHandleQuoteWorker Воркер для последовательной обработки канала котировок
func (s *Stream) startHandleQuoteWorker(ctx context.Context) {
	log.Debug("startHandleQuoteWorker")
	//hasHandleQuote := s.handleQuote != nil
	go func() {
		hasHandleQuote := s.handleQuote != nil
		for {
			select {
			case <-ctx.Done():
				log.Debug("startHandleQuoteWorker ctx.Done() = выход")
				return
			case quote, ok := <-s.quoteChan:
				if !ok {
					log.Debug("startHandleQuoteWorker quoteChan закрыт = выход")
					return // канал закрыт
				}
				if hasHandleQuote {
					s.handleQuote(quote)
				}
			}
		}
	}()
}

// startQuoteStream
// собираем данные с подписки QuoteChannel
// делаем подписку на SubscribeQuote
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
	var processedQuote Quote

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
					s.errChan <- err
					s.Reconnect(QuoteChannel)
					return //  выход
				}
			}
			// TODO возможно выделить в отдельный метод (слишком много логики в одном месте)?
			// В потоке приходит массив данных
			quoteSlice = quoteSlice[:0]
			quoteSlice = msg.GetQuote()
			// Обработаем его
			if len(quoteSlice) != 0 {
				for _, rawQuote := range quoteSlice {
					// Отправляем в канал с сырыми данными, если включено
					if s.SendRawQuotes {
						select {
						case s.rawQuoteChan <- rawQuote:
							// успешно отправили
						default:
							log.Error("listenQuoteStream: канал rawQuoteChan переполнен, дропаем", slog.Any("rawQuote", rawQuote))
						}
					}
					// Обработка котировки
					processedQuote.Reset()
					processedQuote, err = s.quoteStore.processQuote(rawQuote) // Обработаем сырые данные. Вернем срез Quote
					if err != nil {
						log.Error("ошибка обработки котировки", slog.Any("rawQuote", rawQuote), slog.String("err", err.Error()))
						continue
					}
					// Пошлем обработанные данные в канал
					select {
					case s.quoteChan <- processedQuote:
						// запись в канал прошла
					default:
						// канал переполнен, дропаем
						log.Error("listenQuoteStream: канал quoteChan переполнен, дропаем", slog.Any("quote", rawQuote))
					}

				}

			}
			// очистим переменные
			msg.Reset()

		}
	}
}
