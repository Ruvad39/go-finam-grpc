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

type TradesFunc func(data *marketdata_service.SubscribeLatestTradesResponse)

// SetAllTradesHandler
// установим функцию для обработки поступления всех сделок
func (s *Stream) SetAllTradesHandler(handler TradesFunc) {
	s.handleAllTrades = handler
}

// startHandleAllTradesWorker Воркер для последовательной обработки канала всех сделок
func (s *Stream) startHandleAllTradesWorker(ctx context.Context) {
	log.Debug("startHandleAllTradesWorker")
	go func() {
		hasHandleAllTrades := s.handleAllTrades != nil
		//var err error
		//var data *marketdata_service.Trade
		var data *marketdata_service.SubscribeLatestTradesResponse
		var ok bool
		for {
			select {
			case <-ctx.Done():
				log.Debug("HandleAllTradesWorker ctx.Done() = выход")
				return
			case data, ok = <-s.allTradesChan:
				if !ok {
					log.Debug("HandleAllTradesWorker quoteChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				//log.Debug("HandleAllTradesWorker", "data", data)
				if hasHandleAllTrades {
					s.handleAllTrades(data)
				}
				// очистим
				data.Reset()
			}
		}
	}()
}

// startAllTradesStream
// собираем данные с подписки AllTradesChannel
// делаем подписку (SubscribeQuote)
// запускаем в отдельном потоке метод для прослушивания стрима (listenQuoteStream)
func (s *Stream) startAllTradesStream(ctx context.Context) error {
	log.Debug("start startBookStream")
	symbols := s.getSymbolsByChannel(AllTradesChannel)
	// если список пустой = выйдем
	if len(symbols) == 0 {
		log.Debug("startAllTradesStream: нет подписок на канал: выйдем")
		return nil // errors.New("no quote symbols found")
	}
	// есть список инструментов для подписки
	log.Debug("StartStream", "AllTradesSymbols", symbols)

	// добавим заголовок с авторизацией (accessToken)
	var err error
	ctx, err = s.client.WithAuthToken(ctx)
	if err != nil {
		return err
	}
	// (2025-04-17) на текущий момент подписку на стакан можно сделать только на один символ
	// сделаем вызов в цикле по списку символов
	// каждый раз создается новый поток
	for _, symbol := range symbols {
		log.Debug("StartAllTradesStream", "Symbol", symbol)
		var stream grpc.ServerStreamingClient[marketdata_service.SubscribeLatestTradesResponse]
		stream, err = s.client.MarketDataService.SubscribeLatestTrades(ctx, NewSubscribeLatestTradesRequest(symbol))
		if err != nil {
			return err
		}
		// в отдельном потоке запустим чтения данных из стрима
		go s.listenAllTradesStream(ctx, stream)
	}

	return err
}

// listenAllTradesStream чтение данных из стрима всех сделок
func (s *Stream) listenAllTradesStream(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeLatestTradesResponse]) {
	log.Debug("start listenAllTradesStream")
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeLatestTradesResponse
	//var dataSlice []*marketdata_service.Trade

	// читаем поток
	for {
		select {
		case <-ctx.Done():
			log.Debug("listenAllTradesStream: ctx.Done(), выходим")
			return
		case <-s.closeChan:
			log.Debug("listenAllTradesStream: closeChan, выходим")
			return
		default:
			msg, err = stream.Recv()
			if err != nil {
				// Проверка на конкретный код ошибки
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					log.Debug("listenAllTradesStream: Recv: контекст отменён (canceled), выходим")
					return
				}
				if err == io.EOF {
					log.Debug("listenAllTradesStream: Поток завершён")
					return // Поток завершён, выходим из цикла
				} else {
					log.Error("listenAllTradesStream:  Ошибка чтения из потока", "err", err.Error())
					//s.errChan <- err
					s.Reconnect(AllTradesChannel)
					return //  выход
				}
			}
			//log.Debug("listenAllTradesStream", "msg", msg)
			if msg == nil {
				log.Debug("listenAllTradesStream", "msg ISNULL", msg)
				msg.Reset()
				break
			}
			select {
			case s.allTradesChan <- msg:
				// успешно отправили
				//log.Debug("listenAllTradesStream", "msg", msg)
			default:
				log.Error("allTradesChan: канал переполнен, дропаем", slog.Any("Trades", msg))
			}
			// В потоке приходит массив данных
			//dataSlice = dataSlice[:0] // предварительно очистим
			//dataSlice = msg.GetTrades()
			//if len(dataSlice) != 0 {
			//	for _, row := range dataSlice {
			//		select {
			//		case s.allTradesChan <- row:
			//			// успешно отправили
			//		default:
			//			log.Error("allTradesChan: канал переполнен, дропаем", slog.Any("Trades", row))
			//		}
			//	}
			//}

			msg.Reset()

		}
	}
}
