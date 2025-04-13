package finam

import (
	"context"
)

// Channel канал для подписки потока данных
type Channel string

const (
	QuoteChannel = Channel("quote") // Подписка на информацию о котировках
	// TODO BookChannel  = Channel("book")  // Подписка на биржевой стакан
	// TODO SubscribeLatestTrades Подписка на все сделки
)

// Subscription
type Subscription struct {
	Symbol  string  `json:"symbol"`
	Channel Channel `json:"channel"`
}

func (c *Client) GetErrorChan() chan error {
	return c.errChan
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
	err := c.startQuoteStream(ctx)
	if err != nil {
		return err
	}
	return nil
}
