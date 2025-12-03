package finam

import (
	"context"
	"math/rand"
	"time"

	v1 "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1"
	orders_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/orders"
	"google.golang.org/grpc"
)

// TradeStream
type TradeStream struct {
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	retryDelay   time.Duration
	client       *Client
	OrderService orders_service.OrdersServiceClient
	onTrade      func(*v1.AccountTrade) // callback функция
	accountId    string                 // Номер счета для подписки
	running      bool                   // признак что уже запустили в работу
}

// NewOrderStream
// создадим стрим сделок по заданному счету
//
// на входе callback функция для обработки данных
func (c *Client) NewTradeStream(parent context.Context, accountId string, callback func(*v1.AccountTrade)) *TradeStream {
	ctx, cancel := context.WithCancel(parent)
	s := &TradeStream{
		ctx:          ctx,
		cancel:       cancel,
		client:       c,
		OrderService: orders_service.NewOrdersServiceClient(c.conn),
		done:         make(chan struct{}),
		retryDelay:   initialDelay,
		accountId:    accountId,
		onTrade:      callback,
	}
	s.running = true
	go s.run()
	return s
}

// NewTradeStreamWithCallback
// создадим стрим сделок по заданному счету
//
// на входе callback функция для обработки данных
//
// стрим НЕ запускается по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewTradeStreamWithCallback(parent context.Context, accountId string, callback func(*v1.AccountTrade)) *TradeStream {
	ctx, cancel := context.WithCancel(parent)
	s := &TradeStream{
		ctx:          ctx,
		cancel:       cancel,
		client:       c,
		OrderService: orders_service.NewOrdersServiceClient(c.conn),
		done:         make(chan struct{}),
		retryDelay:   initialDelay,
		accountId:    accountId,
		onTrade:      callback,
	}

	return s
}

// Start
func (s *TradeStream) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
}
func (s *TradeStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *TradeStream) run() {
	defer func() {
		log.Debug("[TradeStream] exit run()", "accountId", s.accountId)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[TradeStream]", "accountId", s.accountId, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[TradeStream] start reconnect", "accountId", s.accountId, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[TradeStream] context cancelled, stopping", "accountId", s.accountId)
			return
		case <-time.After(s.retryDelay):
			jitter := time.Duration(rand.Int63n(int64(s.retryDelay / 2)))
			s.retryDelay = min(s.retryDelay*2+jitter, maxDelay) // Макс. 50 сек
		}

	}
}

// subscribeAndListen
// делаем подписку (stream.Send)
// запускаем в отдельном потоке метод для прослушивания стрима (listen)
func (s *TradeStream) subscribeAndListen() error {
	log.Debug("[OrderStream].subscribeAndListen", "accountId", s.accountId)

	// создаем стрим
	stream, err := s.OrderService.SubscribeTrades(s.ctx, &orders_service.SubscribeTradesRequest{AccountId: s.accountId})
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

func (s *TradeStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[orders_service.SubscribeTradesResponse]) error {
	log.Debug("[TradeStream].listen", "accountId", s.accountId)
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
func (s *TradeStream) handleMessage(msg *orders_service.SubscribeTradesResponse) {
	if msg.GetTrades() != nil {
		// обработка
		for _, trade := range msg.GetTrades() {
			if s.onTrade != nil {
				s.onTrade(trade)
			}
		}
	}

}
