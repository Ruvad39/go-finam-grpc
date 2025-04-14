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

// startQuoteStream
// собираем данные с подписки QuoteChannel
// делаем подписку на SubscribeQuote
// запускаем в отдельном потоке (Воркер) для последовательной обработки (handleQuote)
// запускаем в отдельном потоке метод для прослушивания стрима (listenQuoteStream)
func (s *Stream) startQuoteStream(ctx context.Context) error {
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

	// Воркер для последовательной обработки канала котировок
	hasHandleQuote := s.handleQuote != nil
	go func() {
		for quote := range s.quoteChan {
			// только если установили функцию для обработки
			if hasHandleQuote {
				s.handleQuote(quote)
			}
		}
	}()

	// в отдельном потоке запустим чтения данных из стрима
	go s.listenQuoteStream(ctx, stream)

	return err
}

// listenQuoteStream чтение данных из стрима котировок
func (s *Stream) listenQuoteStream(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]) {
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeQuoteResponse
	var quoteSlice []*marketdata_service.Quote
	var processedQuote Quote

	// TODO разобраться с ошибкой отправки в "закрытый" канал ошибок. Не закрывать канал?
	defer func() {
		if r := recover(); r != nil {
			log.Debug("предотвращена паника при отправке в errChan:", "recover", r)
		}
	}()

	// читаем поток
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closeChan:
			return
		default:
			msg, err = stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Debug("StreamQuote Поток завершён")
					break // Поток завершён, выходим из цикла
				} else {
					// TODO решить что делать с ошибкой. Как ее обрабатывать.
					// пока пошлем в канал ошибок и выйдем
					log.Error("StreamQuote Ошибка чтения из потока", "err", err.Error())
					if err != nil {
						s.errChan <- err
					}
					break // test
					//return //  выход
				}
			}
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
