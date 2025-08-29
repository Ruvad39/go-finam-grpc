package finam

import (
	"context"
	"math/rand"
	"time"

	marketdata_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
)

type BarStream struct {
	ctx               context.Context
	cancel            context.CancelFunc
	done              chan struct{}
	retryDelay        time.Duration
	client            *Client
	MarketDataService marketdata_service.MarketDataServiceClient
	symbol            string                       // На какой инструмент подписка
	timeframe         marketdata_service.TimeFrame // На какой тайм-фрейм подписка
	// Колбэки для обработки данных
	//onBar func(bars *marketdata_service.SubscribeBarsResponse) // пока весь ответ, потом одну свечу???
	onBar BarFunc
}

// NewOrderBookStream
// создадим стрим по заданному символу
// данные будем возвращать в метод callback
func (c *Client) NewBarStream(parent context.Context, symbol string, timeframe marketdata_service.TimeFrame, callback func(bar *Bar)) *BarStream {
	ctx, cancel := context.WithCancel(parent)
	s := &BarStream{
		ctx:               ctx,
		cancel:            cancel,
		client:            c,
		MarketDataService: marketdata_service.NewMarketDataServiceClient(c.conn),
		symbol:            symbol,
		timeframe:         timeframe,
		onBar:             callback,
		done:              make(chan struct{}),
		retryDelay:        initialDelay,
	}

	go s.run()
	return s
}

func (s *BarStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *BarStream) run() {
	defer func() {
		log.Debug("[BarStream] exit run()", "symbol", s.symbol, "timeframe", s.timeframe)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[BarStream]", "symbol", s.symbol, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[BarStream] start reconnect", "symbol", s.symbol, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[BarStream] context cancelled, stopping", "symbol", s.symbol)
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
func (s *BarStream) subscribeAndListen() error {
	log.Debug("[BarStream] subscribeAndListen", "symbol", s.symbol, "timeframe", s.timeframe)

	stream, err := s.MarketDataService.SubscribeBars(s.ctx, &marketdata_service.SubscribeBarsRequest{Symbol: s.symbol, Timeframe: s.timeframe})
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

func (s *BarStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[marketdata_service.SubscribeBarsResponse]) error {
	log.Debug("[BarStream] listen", "symbol", s.symbol, "timeframe", s.timeframe)
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
func (s *BarStream) handleMessage(msg *marketdata_service.SubscribeBarsResponse) {
	if s.onBar == nil {
		return
	}
	// GetSymbol()
	// GetBars() []*Bar

	//log.Info("BarStream", "len(o.Rows)", len(msg.GetBars()))
	for _, bar := range msg.GetBars() {
		//log.Info("OrderBookStream.handleMessage", "n", n, "bar", bar)
		newBar := &Bar{
			Symbol:    s.symbol,
			Timeframe: s.timeframe,
			Timestamp: bar.Timestamp,
			Open:      bar.Open,
			High:      bar.High,
			Low:       bar.Low,
			Close:     bar.Close,
			Volume:    bar.Volume,
		}

		s.onBar(newBar)
	}

}
