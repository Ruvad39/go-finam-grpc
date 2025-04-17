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

// type RawOrderBookFunc func(quote *marketdata_service.OrderBook_Row)
type StreamOrderBookFunc func(book *marketdata_service.StreamOrderBook)

// SetRawOrderBookHandler
// установим функцию для обработки поступления "сырых" данных по стакану
func (s *Stream) SetRawOrderBookHandler(handler StreamOrderBookFunc) {
	s.handleRawOrderBook = handler
}

// startHandleQuoteWorker Воркер для последовательной обработки канала котировок
func (s *Stream) startHandleOrderBookWorker(ctx context.Context) {
	log.Debug("startHandleOrderBookWorker")
	go func() {
		hasHandleRawOrderBook := s.handleRawOrderBook != nil
		//hasHandleQuote := s.handleQuote != nil
		//var err error
		var data *marketdata_service.StreamOrderBook
		var ok bool
		for {
			select {
			case <-ctx.Done():
				log.Debug("HandleOrderBookWorker ctx.Done() = выход")
				return
			case data, ok = <-s.rawOrderBookChan:
				if !ok {
					log.Debug("HandleOrderBookWorker quoteChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				if hasHandleRawOrderBook {
					s.handleRawOrderBook(data)
				}
				// TODO возможно тут будет вызов метода "построить стакан"

				// очистим
				data.Reset()
			}
		}
	}()

}

// startBookStream
// собираем данные с подписки BookChannel
// делаем подписку (SubscribeQuote)
// запускаем в отдельном потоке метод для прослушивания стрима (listenQuoteStream)
func (s *Stream) startBookStream(ctx context.Context) error {
	log.Debug("start startBookStream")
	symbols := s.getSymbolsByChannel(BookChannel)
	// если список пустой = выйдем
	if len(symbols) == 0 {
		log.Debug("startBookStream: нет подписок на канал: выйдем")
		return nil // errors.New("no quote symbols found")
	}
	// есть список инструментов для подписки
	log.Debug("StartStream", "BookSymbols", symbols)

	// добавим заголовок с авторизацией (accessToken)
	var err error
	ctx, err = s.client.WithAuthToken(ctx)
	if err != nil {
		return err
	}
	// (2025-04-16) на текущий момент подписку на стакан можно сделать только на один символ
	// сделаем вызов в цикле по списку символов
	// каждый раз создается новый поток
	for _, symbol := range symbols {
		log.Debug("StartBookStream", "Symbol", symbol)
		var stream grpc.ServerStreamingClient[marketdata_service.SubscribeOrderBookResponse]
		stream, err = s.client.MarketDataService.SubscribeOrderBook(ctx, NewSubscribeOrderBookRequest(symbol))
		if err != nil {
			return err
		}
		// в отдельном потоке запустим чтения данных из стрима
		go s.listenBookStream(ctx, stream)
	}

	return err
}

// listenQuoteStream чтение данных из стрима стакана
func (s *Stream) listenBookStream(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeOrderBookResponse]) {
	log.Debug("start listenBookStream")
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeOrderBookResponse
	var dataSlice []*marketdata_service.StreamOrderBook

	// читаем поток
	for {
		select {
		case <-ctx.Done():
			log.Debug("listenBookStream: ctx.Done(), выходим")
			return
		case <-s.closeChan:
			log.Debug("listenBookStream: closeChan, выходим")
			return
		default:
			msg, err = stream.Recv()
			if err != nil {
				// Проверка на конкретный код ошибки
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					log.Debug("listenBookStream: Recv: контекст отменён (canceled), выходим")
					return
				}
				if err == io.EOF {
					log.Debug("listenBookStream: Поток завершён")
					return // Поток завершён, выходим из цикла
				} else {
					log.Error("listenBookStream:  Ошибка чтения из потока", "err", err.Error())
					//s.errChan <- err
					s.Reconnect(QuoteChannel)
					return //  выход
				}
			}
			//StreamOrderBook
			// В потоке приходит массив данных
			dataSlice = dataSlice[:0] // предварительно очистим
			dataSlice = msg.GetOrderBook()
			if len(dataSlice) != 0 {
				for _, row := range dataSlice {
					select {
					case s.rawOrderBookChan <- row:
						// успешно отправили
					default:
						log.Error("OrderBookStream: канал переполнен, дропаем", slog.Any("rawOrderBook", row))
					}
				}
			}

			msg.Reset()

		}
	}
}
