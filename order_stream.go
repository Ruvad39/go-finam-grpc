package finam

import (
	"context"
	"math/rand"
	"time"

	orders_service "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"google.golang.org/grpc"
)

// OrderStream
type OrderStream struct {
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	retryDelay   time.Duration
	client       *Client
	OrderService orders_service.OrdersServiceClient
	onOrder      func(*orders_service.OrderState) // callback функция
	accountId    string                           // Номер счета для подписки
	running      bool                             // признак что уже запустили в работу
}

// NewOrderStream создание стрима ордеров
func (c *Client) NewOrderStream(parent context.Context, accountId string, callbackOrder func(*orders_service.OrderState)) *OrderStream {
	ctx, cancel := context.WithCancel(parent)
	s := &OrderStream{
		ctx:          ctx,
		cancel:       cancel,
		client:       c,
		OrderService: orders_service.NewOrdersServiceClient(c.conn),
		done:         make(chan struct{}),
		retryDelay:   initialDelay,
		accountId:    accountId,
		onOrder:      callbackOrder,
	}
	s.running = true
	go s.run()
	return s
}

// NewOrderStreamWithCallback
// создадим стрим ордеров по заданному счету
//
// на входе callback функция для обработки данных
//
// стрим НЕ запускается по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewOrderStreamWithCallback(parent context.Context, accountId string, callbackOrder func(*orders_service.OrderState)) *OrderStream {
	ctx, cancel := context.WithCancel(parent)
	s := &OrderStream{
		ctx:          ctx,
		cancel:       cancel,
		client:       c,
		OrderService: orders_service.NewOrdersServiceClient(c.conn),
		done:         make(chan struct{}),
		retryDelay:   initialDelay,
		accountId:    accountId,
		onOrder:      callbackOrder,
	}

	return s
}

// Start
func (s *OrderStream) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
}
func (s *OrderStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *OrderStream) run() {
	defer func() {
		log.Debug("[OrderStream] exit run()", "accountId", s.accountId)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[OrderStream]", "accountId", s.accountId, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[OrderStream] start reconnect", "accountId", s.accountId, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[OrderStream] context cancelled, stopping", "accountId", s.accountId)
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
func (s *OrderStream) subscribeAndListen() error {
	log.Debug("[OrderStream].subscribeAndListen", "accountId", s.accountId)

	// создаем стрим
	stream, err := s.OrderService.SubscribeOrders(s.ctx, &orders_service.SubscribeOrdersRequest{AccountId: s.accountId})
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

func (s *OrderStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[orders_service.SubscribeOrdersResponse]) error {
	log.Debug("[OrderStream].listen", "accountId", s.accountId)
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
func (s *OrderStream) handleMessage(msg *orders_service.SubscribeOrdersResponse) {
	// s.handleOrders(msg.GetOrders())
	if msg.GetOrders() != nil {
		// обработка в цикле
		for _, order := range msg.GetOrders() {
			// отправим в callback функцию
			if s.onOrder != nil {
				s.onOrder(order)
			}
		}
	}

}
