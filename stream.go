package finam

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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
		case <-c.CloseC:
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Debug("StreamQuote Поток завершён")
					break // Поток завершён, выходим из цикла
				} else {
					// TODO решить что делать с ошибкой. Как ее обрабатывать. Есть ли метод переподключенния?
					slog.Error("StreamQuote Ошибка чтения из потока", "err", err.Error())
					return //  выход
				}
			}
			// TODO пошлем в канал
			fmt.Printf("StreamQuote msg: %v\n", msg)
			//log.Info("StreamQuote", "msg", msg)

		}
	}
}
