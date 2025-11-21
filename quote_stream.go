/*

 */

package finam

import (
	"context"
	"math/rand"
	"time"

	marketdata_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
)

type QuoteStream struct {
	ctx               context.Context
	cancel            context.CancelFunc
	done              chan struct{}
	retryDelay        time.Duration
	client            *Client
	MarketDataService marketdata_service.MarketDataServiceClient
	symbols           []string                          // список инструментов для подписки
	onQuote           func(q *marketdata_service.Quote) // callback функция для обработки данных
	quoteChannel      chan *marketdata_service.Quote    // канал для данных (альтернатива callback функции)
	running           bool                              // признак что уже запустили в работу
}

// NewQuoteStreamWithCallback
// создадим стрим по заданному символу
// на входе callback функция для обработки данных
// стрим НЕ запускам по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewQuoteStreamWithCallback(parent context.Context, symbols []string, callback func(q *marketdata_service.Quote)) *QuoteStream {
	ctx, cancel := context.WithCancel(parent)
	s := &QuoteStream{
		ctx:               ctx,
		cancel:            cancel,
		client:            c,
		MarketDataService: marketdata_service.NewMarketDataServiceClient(c.conn),
		symbols:           symbols,
		onQuote:           callback,
		done:              make(chan struct{}),
		retryDelay:        initialDelay,
	}

	return s
}

// NewQuoteStreamWithChannel
// создадим стрим по заданному символу
// возвращаем канал с данными
// стрим НЕ запускам по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewQuoteStreamWithChannel(parent context.Context, symbols []string) (*QuoteStream, chan *marketdata_service.Quote) {
	ctx, cancel := context.WithCancel(parent)
	s := &QuoteStream{
		ctx:               ctx,
		cancel:            cancel,
		client:            c,
		MarketDataService: marketdata_service.NewMarketDataServiceClient(c.conn),
		symbols:           symbols,
		done:              make(chan struct{}),
		retryDelay:        initialDelay,
		quoteChannel:      make(chan *marketdata_service.Quote, 100),
	}

	return s, s.quoteChannel
}

// Start
func (s *QuoteStream) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
}

// Close
func (s *QuoteStream) Close() {
	log.Debug("[QuoteStream] Close()", "symbols", s.symbols)
	s.cancel()
	// закроем канал
	if s.quoteChannel != nil {
		close(s.quoteChannel)
	}
	<-s.done // дождаться завершения run()
}

// run
func (s *QuoteStream) run() {
	defer func() {
		log.Debug("[QuoteStream] exit run()", "symbols", s.symbols)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[QuoteStream]", "symbols", s.symbols, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[QuoteStream] start reconnect", "symbols", s.symbols, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[QuoteStream] context cancelled, stopping", "symbols", s.symbols)
			return
		case <-time.After(s.retryDelay):
			jitter := time.Duration(rand.Int63n(int64(s.retryDelay / 2)))
			s.retryDelay = min(s.retryDelay*2+jitter, maxDelay) // Макс. 50 сек
		}

	}
}

// subscribeAndListen
// создаем стрим
// запускаем в отдельном потоке метод для прослушивания стрима (listen)
func (s *QuoteStream) subscribeAndListen() error {
	log.Debug("[QuoteStream] subscribeAndListen", "symbol", s.symbols)

	stream, err := s.MarketDataService.SubscribeQuote(s.ctx, &marketdata_service.SubscribeQuoteRequest{Symbols: s.symbols})
	if err != nil {
		// критичная ошибка = должен быть полный выход
		s.Close() //s.cancel()
		return err
	}
	// успешный коннект = обнулим время
	s.retryDelay = initialDelay
	// запустим чтения данных из стрима
	return s.listen(s.ctx, stream)

}

func (s *QuoteStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeQuoteResponse]) error {
	log.Debug("[QuoteStream] listen", "symbol", s.symbols)
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			msg, err := stream.Recv()
			if err != nil {
				// Проверка на конкретный код ошибки в run()
				return err
			}
			s.handleMessage(msg)
		}
	}
}

// handleMessage обработка сообщения
func (s *QuoteStream) handleMessage(msg *marketdata_service.SubscribeQuoteResponse) {

	//log.Info("QuoteStream", "len(o.Rows)", len(msg.GetQuote()))
	for _, q := range msg.GetQuote() {
		//log.Info("QuoteStream.handleMessage", "i", i, "Quote", q)

		// отправим в канал
		if s.quoteChannel != nil {
			s.quoteChannel <- q
		}
		// отправим в callback функцию
		if s.onQuote != nil {
			s.onQuote(q)
		}

	}

}
