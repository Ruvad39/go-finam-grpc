package finam

import (
	"context"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"io"
)

// Channel канал для подписки потока данных
type Channel string

const (
	QuoteChannel = Channel("quote") // Подписка на информацию о котировках
	BookChannel  = Channel("book")  // Подписка на биржевой стакан
	// TODO SubscribeLatestTrades Подписка на все сделки
)

// Subscription
type Subscription struct {
	Symbol  string  `json:"symbol"`
	Channel Channel `json:"channel"`
}

func (c *Client) GetQuoteChan() chan *marketdata_service.Quote {
	return c.quoteChan
}

// Subscribe подписка на поток информации
func (c *Client) Subscribe(channel Channel, symbol string) {
	log.Debug("Subscribe", "channel", channel, "symbol", symbol)
	sub := Subscription{
		Channel: channel,
		Symbol:  symbol,
	}
	c.subscriptions[sub] = sub
}

// getSymbolsByChannel
// вернем список инструментов по заданному типу подписки
func (c *Client) getSymbolsByChannel(channel Channel) []string {
	symbols := []string{}
	for sub := range c.subscriptions {
		if sub.Channel == channel {
			symbols = append(symbols, sub.Symbol)
		}
	}
	return symbols
}

func (c *Client) groupSymbolsByChannel() map[Channel][]string {
	result := make(map[Channel][]string)
	for sub := range c.subscriptions {
		result[sub.Channel] = append(result[sub.Channel], sub.Symbol)
	}
	return result
}

// StartStream
// собираем данные с подписок
// вызываем методы акивации нужного потока
func (c *Client) StartStream(ctx context.Context) error {
	// (1) QuoteChannel
	symbols := c.getSymbolsByChannel(QuoteChannel)
	if len(symbols) != 0 {
		log.Debug("StartStream", "QuoteSymbols", symbols)
		go c.startStreamQuote(ctx, symbols)
	}
	return nil
}

// startStreamQuote
// активация подписки
// чтения потока и отправка в канал
func (c *Client) startStreamQuote(ctx context.Context, symbols []string) {
	var err error
	// сразу создадим переменные, что бы их переиспользовать
	var msg *marketdata_service.SubscribeQuoteResponse
	var quoteSlice []*marketdata_service.Quote

	// добавим заголовок с авторизацией (accessToken)
	ctx, err = c.WithAuthToken(ctx)
	if err != nil {
		log.Error("startStreamQuote", "WithAuthToken err", err.Error())
		return
	}
	stream, err := c.MarketDataService.SubscribeQuote(ctx, NewSubscribeQuoteRequest(symbols))
	if err != nil {
		log.Error("StreamQuote", "SubscribeQuote err", err.Error())
		return
	}
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
					// TODO решить что делать с ошибкой. Как ее обрабатывать. Есть ли метод переподключенния?
					log.Error("StreamQuote Ошибка чтения из потока", "err", err.Error())
					return //  выход
				}
			}

			//fmt.Printf("StreamQuote msg: %v\n", msg)
			// Приходит массив данных.
			quoteSlice = msg.GetQuote()
			// Обработаем его и пошлем ссылку дальше в канал
			if len(quoteSlice) != 0 {
				for _, quote := range quoteSlice {
					c.quoteChan <- quote
				}
			}
			//log.Debug("StreamQuote", "msg", msg.GetQuote())
			// очистим переменные
			msg.Reset()
			quoteSlice = quoteSlice[:0]

		}
	}
}
