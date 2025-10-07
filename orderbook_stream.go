package finam

import (
	"context"
	"math/rand"
	"time"

	pb "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
)

type OrderBookStream struct {
	ctx               context.Context
	cancel            context.CancelFunc
	done              chan struct{}
	retryDelay        time.Duration
	client            *Client
	MarketDataService pb.MarketDataServiceClient
	symbol            string // На какой инструмент подписка
	// Колбэки для обработки данных
	onUpdate func(book []*pb.StreamOrderBook) // todo пока весь ответ. потом срез стакана???
}

// NewOrderBookStream
// создадим стрим по заданному символу
// данные будем возвращать в метод callback
func (c *Client) NewOrderBookStream(parent context.Context, symbol string, callback func(book []*pb.StreamOrderBook)) *OrderBookStream {
	ctx, cancel := context.WithCancel(parent)
	s := &OrderBookStream{
		ctx:               ctx,
		cancel:            cancel,
		client:            c,
		MarketDataService: pb.NewMarketDataServiceClient(c.conn),
		symbol:            symbol,
		onUpdate:          callback,
		done:              make(chan struct{}),
		retryDelay:        initialDelay,
	}

	go s.run()
	return s
}

func (s *OrderBookStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *OrderBookStream) run() {
	defer func() {
		log.Debug("[OrderBookStream] exit run()", "symbol", s.symbol)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[OrderBookStream]", "symbol", s.symbol, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[OrderBookStream] start reconnect", "symbol", s.symbol, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[OrderBookStream] context cancelled, stopping", "symbol", s.symbol)
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
func (s *OrderBookStream) subscribeAndListen() error {
	log.Debug("[OrderBookStream] subscribeAndListen", "symbol", s.symbol)

	stream, err := s.MarketDataService.SubscribeOrderBook(s.ctx, &pb.SubscribeOrderBookRequest{Symbol: s.symbol})
	if err != nil {
		// критичная ошибка = должен быть полный выход
		s.cancel()
		return err
	}
	// успешный коннект = обнулим время
	s.retryDelay = initialDelay
	// запустим чтения данных из стрима
	return s.listen(s.ctx, stream)

}

func (s *OrderBookStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[pb.SubscribeOrderBookResponse]) error {
	log.Debug("[OrderBookStream] listen", "symbol", s.symbol)
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
			s.handleMessage(msg.GetOrderBook())
		}
	}
}

// handleMessage обработка сообщения
func (s *OrderBookStream) handleMessage(msg []*pb.StreamOrderBook) {
	if s.onUpdate != nil {
		s.onUpdate(msg)
	}
	// TODO OrderBookBuilder

}
