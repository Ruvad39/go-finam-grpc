/*
разделил на 2 канала с данынми

	rawQuoteChan      chan *marketdata_service.Quote // Канал с "сырыми" данными по котировкам
	quoteChan         chan Quote // Канал с обработанными котировками
*/
package finam

import (
	"context"
	"fmt"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"google.golang.org/grpc"
	"io"
	"log/slog"
	"time"
)

// Информация о котировке
type Quote struct {
	Symbol    string // Символ инструмента
	Timestamp int64  // Метка времени
	Ask       float64
	Bid       float64
	Last      float64 // Цена последней сделки
}

func (q *Quote) Time() time.Time {
	return time.Unix(0, q.Timestamp).In(TzMoscow)
}

type QuoteStore struct {
	quoteState map[string]*Quote // Последнее состояние по символу
	//mu         sync.Mutex        // защита для конкурентного доступа
}

// processQuote обработаем сырые
func (qs *QuoteStore) processQuote(rq *marketdata_service.Quote) (Quote, error) {
	if rq == nil || rq.Symbol == "" {
		return Quote{}, fmt.Errorf("некорректная котировка: отсутствует символ")
	}
	// пока уберу. все делается в одном потоке (в listenQuoteStream)
	//qs.mu.Lock()
	//defer qs.mu.Unlock()

	// Получаем текущее состояние, если есть
	q, ok := qs.quoteState[rq.Symbol]
	if !ok {
		q = &Quote{
			Symbol: rq.Symbol,
		}
		qs.quoteState[rq.Symbol] = q
	}

	// Обновляем только непустые поля
	if rq.Timestamp != nil {
		q.Timestamp = rq.Timestamp.AsTime().UnixNano()
	}
	if rq.Ask != nil {
		q.Ask, _ = DecimalToFloat64E(rq.Ask)
	}
	if rq.Bid != nil {
		q.Bid, _ = DecimalToFloat64E(rq.Bid)
	}
	if rq.Last != nil {
		q.Last, _ = DecimalToFloat64E(rq.Last)
	}

	// Возвращаем копию
	return Quote{
		Symbol:    q.Symbol,
		Timestamp: q.Timestamp,
		Ask:       q.Ask,
		Bid:       q.Bid,
		Last:      q.Last,
	}, nil
}

type QuoteFunc func(quote Quote)

// SetQuoteHandler
// установим функцию для обработки поступления котировок
// должен устанавливать раньше Start
func (c *Client) SetQuoteHandler(handler QuoteFunc) {
	c.handleQuote = handler
}

// GetRawQuoteChan вернем канал в котором сырые данные по котировкам
func (c *Client) GetRawQuoteChan() chan *marketdata_service.Quote {
	return c.rawQuoteChan
}

// startQuoteStream
// собираем данные с подписки QuoteChannel
// делаем подписку на SubscribeQuote
// запускаем в отдельном потоке (Воркер) для последовательной обработки (handleQuote)
// запускаем в отдельном потоке метод для прослушивания стрима (listenQuoteStream)
func (c *Client) startQuoteStream(ctx context.Context) error {
	symbols := c.getSymbolsByChannel(QuoteChannel)
	// если список пустой = выйдем
	if len(symbols) == 0 {
		return nil // errors.New("no quote symbols found")
	}
	// есть список инструментов для подписки
	var err error
	var stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]
	log.Debug("StartStream", "QuoteSymbols", symbols)
	// добавим заголовок с авторизацией (accessToken)
	ctx, err = c.WithAuthToken(ctx)
	if err != nil {
		return err
	}

	stream, err = c.MarketDataService.SubscribeQuote(ctx, NewSubscribeQuoteRequest(symbols))
	if err != nil {
		return err
	}

	// Воркер для последовательной обработки канала котировок
	hasHandleQuote := c.handleQuote != nil
	go func() {
		for quote := range c.quoteChan {
			// только если установили функцию для обработки
			if hasHandleQuote {
				c.handleQuote(quote)
			}
		}
	}()

	// в отдельном потоке запустим чтения данных из стрима
	go c.listenQuoteStream(ctx, stream)

	return err
}

// listenQuoteStream чтение данных из стрима котировок
func (c *Client) listenQuoteStream(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]) {
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeQuoteResponse
	var quoteSlice []*marketdata_service.Quote

	// TODO разобраться с ошибкой отправки в "закрытый" канал ошибок
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
		case <-c.closeChan:
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
						c.errChan <- err
					}
					return //  выход
				}
			}
			// Приходит массив данных.
			quoteSlice = msg.GetQuote()
			// Обработаем его и пошлем ссылку дальше в канал
			if len(quoteSlice) != 0 {
				for _, rawQuote := range quoteSlice {
					// Отправляем в канал с сырыми данными, если включено
					if c.SendRawQuotes {
						select {
						case c.rawQuoteChan <- rawQuote:
							// успешно отправили
						default:
							log.Error("listenQuoteStream: канал rawQuoteChan переполнен, дропаем", slog.Any("rawQuote", rawQuote))
						}
					}
					// Обработка котировки
					processedQuote, err := c.quoteStore.processQuote(rawQuote)
					if err != nil {
						log.Error("ошибка обработки котировки", slog.Any("rawQuote", rawQuote), slog.String("err", err.Error()))
						continue
					}

					select {
					case c.quoteChan <- processedQuote:
						// запись в канал прошла
					default:
						// канал переполнен, дропаем
						log.Error("listenQuoteStream: канал quoteChan переполнен, дропаем", slog.Any("quote", rawQuote))
					}

				}

			}
			// очистим переменные
			msg.Reset()
			quoteSlice = quoteSlice[:0]

		}
	}
}

//func (c *Consumer) processQuote(q *marketdata_service.Quote) (*marketdata_service.Quote, error) {
//	// Пример: фильтрация, нормализация, добавление полей и т.п.
//	if q.Price <= 0 {
//		return nil, fmt.Errorf("некорректная цена")
//	}
//
//	// Допустим, мы клонируем quote и что-то меняем:
//	newQuote := *q
//	newQuote.Price = normalizePrice(q.Price)
//	return &newQuote, nil
//}
