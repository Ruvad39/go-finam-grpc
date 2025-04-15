package finam

import (
	"context"
	marketdata_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/marketdata"
	"time"
)

const (
	reconnectDelay  = 10 * time.Second // Интервал повторной попытки реконнекта
	quoteBufferSize = 100              // Размер буфера канала котировок
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

type Stream struct {
	client         *Client
	subscriptions  map[Subscription]Subscription // Список подписок на поток данных
	closeChan      chan struct{}                 // Сигнальный канал для закрытия коннекта
	reconnectChan  chan Channel                  // Сигнальный канал для необходимости реконекта
	errChan        chan error
	rawQuoteChan   chan *marketdata_service.Quote // Канал с "сырыми" данными по котировкам
	quoteChan      chan Quote                     // Канал с обработанными котировками
	handleQuote    QuoteFunc
	SendRawQuotes  bool                                        // Признак, посылать сырые данные или нет
	quoteStore     QuoteStore                                  // Обработчик данных по котировкам
	streamStarters map[Channel]func(ctx context.Context) error // Список методов для запуска потоков
	workerStarters map[Channel]func(context.Context)           // Список воркеров
}

func (c *Client) NewStream() *Stream {
	s := &Stream{
		client:        c,
		closeChan:     make(chan struct{}),
		errChan:       make(chan error, 1),
		quoteChan:     make(chan Quote, quoteBufferSize),
		rawQuoteChan:  make(chan *marketdata_service.Quote, quoteBufferSize),
		subscriptions: make(map[Subscription]Subscription),
		quoteStore: QuoteStore{
			quoteState: make(map[string]*Quote),
		},
		streamStarters: make(map[Channel]func(context.Context) error),
		workerStarters: make(map[Channel]func(context.Context)),
	}
	// Регистрируем стримы
	s.streamStarters[QuoteChannel] = s.startQuoteStream
	//s.streamStarters[OrderChannel] = s.startOrderStream

	// Регистрируем воркеры
	s.workerStarters[QuoteChannel] = s.startHandleQuoteWorker

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

// StartStream
// собираем данные с подписок
// вызываем методы акивации нужного потока
// func (c *Client) StartStream(ctx context.Context) error {
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
		log.Debug("startStreams: Поток успешно запущен", "channel", ch)
	}
	return nil
}

// Close закроем сигнальный канал, что бы закончить работу
func (s *Stream) Close() {
	close(s.closeChan)
	//close(s.quoteChan)
}

// Reconnect в сигнальный канал рекконета пошлем сообщение
// channel какой канал нужно переконнектить
func (s *Stream) Reconnect(channel Channel) {
	log.Debug("зашли в Reconnect()", "channel", channel)
	select {
	case s.reconnectChan <- channel:
	default:
	}
}

// reconnector ждет команду на переподключение
// работает в отдельно потоке
// вызывает startStreams(ctx)
func (s *Stream) reconnector(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closeChan:
			return
		case ch := <-s.reconnectChan: //  получаем значение типа Channel
			log.Warn("re-connecting stream", "channel", ch)
			reconnectFunc, ok := s.streamStarters[ch]
			if !ok {
				log.Error("Нет обработчика для канала", "channel", ch)
				continue
			}
			log.Warn("reconnector", "период восстановления повторного подключения", reconnectDelay)
			time.Sleep(reconnectDelay)
			if err := reconnectFunc(ctx); err != nil {
				log.Error("Ошибка реконнекта", "channel", ch, "err", err)
				s.Reconnect(ch) // повторный сигнал
			}
		}
	}
}
