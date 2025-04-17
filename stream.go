/*
TODO добавить события Connect DisConnect
*/

package finam

import (
	"context"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"time"
)

const (
	reconnectDelay   = 10 * time.Second // Интервал повторной попытки реконнекта
	quoteBufferSize  = 100              // Размер буфера канала котировок
	bookBufferSize   = 100              // Размер буфера канала стакана
	tradesBufferSize = 100              // Размер буфера канала всех сделок
)

// Channel канал для подписки потока данных
type Channel string

const (
	QuoteChannel     = Channel("quote")      // Подписка на информацию о котировках
	BookChannel      = Channel("book")       // Подписка на биржевой стакан
	AllTradesChannel = Channel("all_trades") // подписка на все сделки
)

// Subscription
type Subscription struct {
	Symbol  string  `json:"symbol"`
	Channel Channel `json:"channel"`
}

type Stream struct {
	client         *Client
	subscriptions  map[Subscription]Subscription               // Список подписок на поток данных
	closeChan      chan struct{}                               // Сигнальный канал для закрытия коннекта
	reconnectChan  chan Channel                                // Канал для необходимости реконекта
	errChan        chan error                                  // Сигнальный канал ошибок
	streamStarters map[Channel]func(ctx context.Context) error // Список методов для запуска потоков
	workerStarters map[Channel]func(context.Context)           // Список воркеров
	// для потока котировок
	rawQuoteChan   chan *marketdata_service.Quote // Канал с "сырыми" данными по котировкам
	handleQuote    QuoteFunc                      // Функция обработчик котировок
	handleRawQuote RawQuoteFunc                   // Функция обработчик "сырых" котировок
	quoteStore     QuoteStore                     // Обработчик данных по котировкам
	// для OrderBook
	rawOrderBookChan   chan *marketdata_service.StreamOrderBook // Канал с "сырыми" данными по стакану
	handleRawOrderBook StreamOrderBookFunc                      // Функция обработчик "сырого" стакана
	// LatestTrades
	handleAllTrades TradesFunc                                             // Функция обработчик всех сделок
	allTradesChan   chan *marketdata_service.SubscribeLatestTradesResponse // Канал всех сделки
}

func (c *Client) NewStream() *Stream {
	s := &Stream{
		client:           c,
		closeChan:        make(chan struct{}),
		reconnectChan:    make(chan Channel),
		errChan:          make(chan error, 1),
		rawQuoteChan:     make(chan *marketdata_service.Quote, quoteBufferSize),
		rawOrderBookChan: make(chan *marketdata_service.StreamOrderBook, bookBufferSize),
		allTradesChan:    make(chan *marketdata_service.SubscribeLatestTradesResponse),
		subscriptions:    make(map[Subscription]Subscription),
		quoteStore: QuoteStore{
			quoteState: make(map[string]*Quote),
		},
		streamStarters: make(map[Channel]func(context.Context) error),
		workerStarters: make(map[Channel]func(context.Context)),
	}
	// Регистрируем стримы
	s.streamStarters[QuoteChannel] = s.startQuoteStream
	s.streamStarters[BookChannel] = s.startBookStream
	s.streamStarters[AllTradesChannel] = s.startAllTradesStream

	// Регистрируем воркеры
	s.workerStarters[QuoteChannel] = s.startHandleQuoteWorker
	s.workerStarters[BookChannel] = s.startHandleOrderBookWorker
	s.workerStarters[AllTradesChannel] = s.startHandleAllTradesWorker

	return s
}

func (s *Stream) GetErrorChan() chan error {
	return s.errChan
}

// Subscribe подписка на поток информации
func (s *Stream) Subscribe(channel Channel, symbol string) {
	//func (c *Client) Subscribe(channel Channel, symbol string) {
	log.Debug("Subscribe", "channel", channel, "symbol", symbol)
	sub := Subscription{
		Channel: channel,
		Symbol:  symbol,
	}
	s.subscriptions[sub] = sub
}

// getSymbolsByChannel
// вернем список инструментов по заданному типу подписки
func (s *Stream) getSymbolsByChannel(channel Channel) []string {
	symbols := []string{}
	for sub := range s.subscriptions {
		if sub.Channel == channel {
			symbols = append(symbols, sub.Symbol)
		}
	}
	return symbols
}

func (s *Stream) groupSymbolsByChannel() map[Channel][]string {
	result := make(map[Channel][]string)
	for sub := range s.subscriptions {
		result[sub.Channel] = append(result[sub.Channel], sub.Symbol)
	}
	return result
}

// Connect
// запускаем воркеры (startHandleWorkers(ctx))
// активируем и запускаем потоки (startStreams(ctx))
// запускаем реконектор (reconnector(ctx))
func (s *Stream) Connect(ctx context.Context) error {
	// запустим воркеры для чтения каналов
	s.startHandleWorkers(ctx)

	// вызовем подписку и запуск потоков
	err := s.startStreams(ctx)
	if err != nil {
		return err
	}

	go s.reconnector(ctx)
	return nil
}

// startHandleWorker запустим воркеры для чтения каналов
func (s *Stream) startHandleWorkers(ctx context.Context) {
	for ch, startWorker := range s.workerStarters {
		log.Debug("startHandleWorkers: Запускаем воркер", "channel", ch)
		go startWorker(ctx)
	}
}

// startStreams запуск потоков
func (s *Stream) startStreams(ctx context.Context) error {
	for ch, starter := range s.streamStarters {
		if err := starter(ctx); err != nil {
			log.Error("startStreams: Не удалось запустить поток", "channel", ch, "err", err)
			return err
		}
		//log.Debug("startStreams: Поток успешно запущен", "channel", ch)
	}
	return nil
}

// Close закроем сигнальный канал, что бы закончить работу
func (s *Stream) Close() {
	close(s.closeChan)
}

// Reconnect в канал рекконета пошлем сообщение
// channel какой канал нужно переконнектить
func (s *Stream) Reconnect(ch Channel) {
	log.Debug("зашли в Reconnect()", "channel", ch)
	select {
	case s.reconnectChan <- ch:
		log.Debug("Reconnect: Успешно отправили в канал", "channel", ch)
		//default:
		//	log.Warn("Reconnect: Канал переполнен, сигнал пропущен", "channel", ch)
	}
}

// reconnector ждет команду на переподключение
// работает в отдельно потоке
// вызывает streamStarters заданного канала
func (s *Stream) reconnector(ctx context.Context) {
	log.Debug("start reconnector")
	for {
		select {
		case <-ctx.Done():
			log.Debug("reconnector: ctx.Done() = выход")
			return
		case <-s.closeChan:
			log.Debug("reconnector: closeChan = выход")
			return
		case ch := <-s.reconnectChan: //  получаем значение типа Channel
			log.Warn("re-connecting stream", "channel", ch)
			reconnectFunc, ok := s.streamStarters[ch]
			if !ok {
				log.Error("Нет обработчика для канала", "channel", ch)
				continue
			}
			log.Warn("reconnector: попытка восстановить поток", "channel", ch)
			if err := reconnectFunc(ctx); err != nil {
				log.Error("Ошибка реконнекта", "channel", ch, "err", err.Error())
				log.Warn("reconnector: новая попытка восстановить поток", "channel", ch, "через", reconnectDelay)
				time.Sleep(reconnectDelay)
				s.Reconnect(ch) // повторный сигнал
			}
		}
	}
}
