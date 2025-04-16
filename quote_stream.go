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

// startHandleQuoteWorker Воркер для последовательной обработки канала котировок
func (s *Stream) startHandleQuoteWorker(ctx context.Context) {
	log.Debug("startHandleQuoteWorker")
	// "сырые" котировки
	go func() {
		hasHandleRawQuote := s.handleRawQuote != nil
		hasHandleQuote := s.handleQuote != nil
		var err error
		var processedQuote Quote
		var data *marketdata_service.Quote
		var ok bool
		for {
			select {
			case <-ctx.Done():
				log.Debug("HandleRawQuoteWorker ctx.Done() = выход")
				return
			case data, ok = <-s.rawQuoteChan:
				if !ok {
					log.Debug("HandleRawQuoteWorker quoteChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				if hasHandleRawQuote {
					s.handleRawQuote(data)
				}
				// обработаем инкрементальные данные в срез
				// только если указана функция приема
				if hasHandleQuote {
					// Обработка котировки
					processedQuote.Reset()
					processedQuote, err = s.quoteStore.processQuote(data) // Обработаем сырые данные. Вернем Quote
					if err != nil {
						log.Error("ошибка обработки котировки", slog.Any("rawQuote", data), slog.String("err", err.Error()))
						continue
					}
					// пошлем получателю
					s.handleQuote(processedQuote)
				}
				// очистим
				data.Reset()
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
		log.Debug("startQuoteStream: нет подписок на канал: выйдем")
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
	var dataSlice []*marketdata_service.Quote

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
			dataSlice = dataSlice[:0] // предварительно очистим
			dataSlice = msg.GetQuote()
			if len(dataSlice) != 0 {
				for _, row := range dataSlice {
					select {
					case s.rawQuoteChan <- row:
						// успешно отправили
					default:
						log.Error("QuoteStream: канал переполнен, дропаем", slog.Any("rawQuote", row))
					}
				}
			}
			// очистим переменные
			msg.Reset()
		}
	}
}
