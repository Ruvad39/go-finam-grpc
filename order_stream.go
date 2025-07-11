/*
Стрим для ордеров и сделок

опытным путем определил, что на один стрим OrdersService.SubscribeOrderTrade можно подписаться на разные счета одновременно.

пример работы со стримом:
*/
package finam

import (
	"context"
	"fmt"
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

const (
	orderBufferSize = 100 // Размер буфера канала ордеров
	tradeBufferSize = 100 // Размер буфера канала сделок
)

// типы подписок
const (
	OrderTradeChannel = pb.OrderTradeRequest_DATA_TYPE_ALL    // Подписка на Ордера и на Сделки
	OrderChannel      = pb.OrderTradeRequest_DATA_TYPE_ORDERS // Подписка только на Ордера
	TradeChannel      = pb.OrderTradeRequest_DATA_TYPE_TRADES // Подписка только на Сделки
)

// OrderTradeStream
type OrderTradeStream struct {
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	retryDelay   time.Duration
	client       *Client
	OrderService pb.OrdersServiceClient
	onOrder      func(*pb.OrderState)
	onTrade      func(*pb.AccountTrade)
	accountId    string // Номер счета для подписки
}

// callback func(book []*pb.StreamOrderBook)
func (c *Client) NewOrderTradeStream(parent context.Context,
	accountId string,
	callbackOrder func(*pb.OrderState),
	callbackTrade func(*pb.AccountTrade),
) *OrderTradeStream {
	ctx, cancel := context.WithCancel(parent)
	s := &OrderTradeStream{
		ctx:          ctx,
		cancel:       cancel,
		client:       c,
		OrderService: pb.NewOrdersServiceClient(c.conn),
		done:         make(chan struct{}),
		retryDelay:   initialDelay,
		accountId:    accountId,
		onOrder:      callbackOrder,
		onTrade:      callbackTrade,
	}
	go s.run()
	return s
}

func (s *OrderTradeStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *OrderTradeStream) run() {
	defer func() {
		log.Debug("[OrderTradeStream] exit run()", "accountId", s.accountId)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[OrderTradeStream]", "accountId", s.accountId, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[OrderTradeStream] start reconnect", "accountId", s.accountId, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[OrderTradeStream] context cancelled, stopping", "accountId", s.accountId)
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
func (s *OrderTradeStream) subscribeAndListen() error {
	log.Debug("[OrderTradeStream].subscribeAndListen", "accountId", s.accountId)

	// добавим заголовок с авторизацией (accessToken)
	ctx, err := s.client.WithAuthToken(s.ctx)
	if err != nil {
		return err
	}
	// создаем стрим
	stream, err := s.OrderService.SubscribeOrderTrade(ctx)
	if err != nil {
		return err
	}

	// тип подписки
	dateType := pb.OrderTradeRequest_DATA_TYPE_ALL
	if s.onOrder == nil {
		dateType = pb.OrderTradeRequest_DATA_TYPE_TRADES
	}
	if s.onTrade == nil {
		dateType = pb.OrderTradeRequest_DATA_TYPE_ORDERS
	}
	// Отправляем запрос подписки
	if err = stream.Send(&pb.OrderTradeRequest{
		Action:    pb.OrderTradeRequest_ACTION_SUBSCRIBE,
		DataType:  dateType,
		AccountId: s.accountId,
	}); err != nil {
		//return err
		return fmt.Errorf("stream.Send err: %w", err)
	}
	//stream.CloseSend()

	// чтение потока
	return s.listen(ctx, stream)

}

func (s *OrderTradeStream) listen(ctx context.Context, stream grpc.BidiStreamingClient[pb.OrderTradeRequest, pb.OrderTradeResponse]) error {
	log.Debug("[OrderTradeStream] listenMessage", "accountId", s.accountId)
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
func (s *OrderTradeStream) handleMessage(msg *pb.OrderTradeResponse) {
	//log.Info("OrderTradeStream.handleMessage", "msg", msg)
	// DEBUG
	//log.Info("------------------")
	//log.Info("OrderTradeStream.handleMessag = пришел пакет")
	//log.Info("OrderTradeStream.handleMessage", "len(msg.GetOrders())", len(msg.GetOrders()))
	//for n, order := range msg.GetOrders() {
	//	log.Info("GetOrders()", "n", n, "order", order)
	//}
	//log.Info("OrderTradeStream.handleMessage", "len(msg.GetTrades())", len(msg.GetTrades()))
	//for n, trade := range msg.GetTrades() {
	//	log.Info("GetTrades()", "n", n, "trade", trade)
	//}
	//log.Info("------------------")

	s.handleOrders(msg.GetOrders())
	s.handleTrades(msg.GetTrades())

}

// handleOrders обработка ордеров
func (s *OrderTradeStream) handleOrders(orders []*pb.OrderState) {
	if orders != nil {
		// обработка
		for _, order := range orders {
			if s.onOrder != nil {
				s.onOrder(order)
			}
		}
	}
}

// handleTrades обработка сделок
func (s *OrderTradeStream) handleTrades(trades []*pb.AccountTrade) {
	if trades != nil {
		// обработка
		for _, trade := range trades {
			if s.onTrade != nil {
				s.onTrade(trade)
			}
		}
	}
}
