/*
Стрим для ордеров и сделок

опытным путем определил, что на один стрим OrdersService.SubscribeOrderTrade можно подписаться на разные счета одновременно.

пример работы со стримом:

	// создадим поток
	stream := client.NewOrderTradeStream()
	// подпишемся на ордера и сделки по счету
	stream.Subscribe("номер счета", finam.OrderTradeChannel)

	// или работаем с каналом
	// получим канал
	orderChan := stream.OrderChan()
	// .. обработка канала

	// или  установим функцию обработчик
	// в этом случае автоматом запуститься startOrderWorker (чтение канала о отправка данных в onOrder)
	stream.SetOnOrder(OnOrder)

	// запустим поток в работу
	err = stream.Start(ctx)
	if err != nil {
		slog.Error("stream.Start", "err", err.Error())
	}
*/
package finam

import (
	"context"
	"fmt"
	pb "github.com/Ruvad39/go-finam-grpc/tradeapi/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
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
	BaseStream
	client       *Client
	OrderService pb.OrdersServiceClient
	stream       grpc.BidiStreamingClient[pb.OrderTradeRequest, pb.OrderTradeResponse]
	// Каналы
	orderChan chan pb.OrderState   // Канал ордеров
	tradeChan chan pb.AccountTrade // Канал сделок
	// Колбэки для обработки данных
	onOrder func(pb.OrderState)
	onTrade func(pb.AccountTrade)
	//onError func(error)
	// На что подписываемся
	accountId string                        // Номер счета
	dateType  pb.OrderTradeRequest_DataType // Тип данных подписки
	// todo subscriptions map[string]pb.OrderTradeRequest_DataType // Список счетов в подписке
}

func (c *Client) NewOrderTradeStream() *OrderTradeStream {
	s := &OrderTradeStream{
		BaseStream: BaseStream{
			reconnectCh: make(chan struct{}, 1),
			closeCh:     make(chan struct{}),
		},
		client:       c,
		OrderService: pb.NewOrdersServiceClient(c.conn),
		orderChan:    make(chan pb.OrderState, orderBufferSize),
		tradeChan:    make(chan pb.AccountTrade, tradeBufferSize),
		//subscriptions: make(map[string]pb.OrderTradeRequest_DataType),
	}
	// для отладки
	s.Name = "OrderTradeStream"
	// привяжем функцию для "работы"
	s.ConnectFunc = s.subscribeAndListen
	return s
}

// OrderChan вернем канал ордеров
func (s *OrderTradeStream) OrderChan() chan pb.OrderState {
	return s.orderChan
}

// Subscribe
// пока подписка только на один счет
// если вызовем несколько раз = запишется последний
func (s *OrderTradeStream) Subscribe(accountId string, typeChannel pb.OrderTradeRequest_DataType) {
	// пока подписка будет только на один счет
	s.accountId = accountId
	s.dateType = typeChannel
	// в дальнейшем возможно будет на список счетов
	//s.subscriptions[accountId] = typeChannel
}

// SetOnOrder установить callback функцию на ордер
func (s *OrderTradeStream) SetOnOrder(callback func(pb.OrderState)) {
	s.onOrder = callback
	// добавим воркер в список для автоматического запуска при Start()
	s.AddWorker(s.startOrderWorker)
}

// SetOnTrade установить callback функцию на сделку
func (s *OrderTradeStream) SetOnTrade(callback func(pb.AccountTrade)) {
	s.onTrade = callback
	// добавим воркер в список для автоматического запуска при Start()
	s.AddWorker(s.startTradeWorker)
}

// SetOnError установить callback функцию на ошибку
//func (s *OrderTradeStream) SetOnError(callback func(error)) {
//	s.onError = callback
//}

// subscribeAndListen
// делаем подписку (stream.Send)
// запускаем в отдельном потоке метод для прослушивания стрима (listen)
func (s *OrderTradeStream) subscribeAndListen(ctx context.Context) error {
	log.Info("OrderTradeStream.subscribeAndListen")
	var err error
	// добавим заголовок с авторизацией (accessToken)
	ctx, err = s.client.WithAuthToken(ctx)
	if err != nil {
		return err
	}
	// создаем стрим
	stream, err := s.OrderService.SubscribeOrderTrade(ctx)
	if err != nil {
		return err
	}
	//// Отправляем запрос подписки
	if err = stream.Send(&pb.OrderTradeRequest{
		Action:    pb.OrderTradeRequest_ACTION_SUBSCRIBE,
		DataType:  s.dateType, //pb.OrderTradeRequest_DATA_TYPE_ALL,
		AccountId: s.accountId,
	}); err != nil {
		return fmt.Errorf("subscribe request failed: %w", err)
	}

	// TODO подписка на список счетов
	//for accountID, dataType := range s.subscriptions {
	//	log.Debug("OrderTradeStream.subscribeAndListen subscriptions:", "accountID", accountID, "dataType", dataType)
	//	// Отправляем запрос подписки
	//	if err = stream.Send(&pb.OrderTradeRequest{
	//		Action:    pb.OrderTradeRequest_ACTION_SUBSCRIBE,
	//		DataType:  dataType,
	//		AccountId: accountID,
	//	}); err != nil {
	//		return fmt.Errorf("subscribe request failed: %w", err)
	//	}
	//
	//}

	// Запускаем обработчик сообщений в отдельной горутине
	// в отдельном потоке запустим чтения данных из стрима
	go s.listenMessage(ctx, stream)

	return err
}

func (s *OrderTradeStream) listenMessage(ctx context.Context, stream grpc.BidiStreamingClient[pb.OrderTradeRequest, pb.OrderTradeResponse]) {
	log.Info("OrderTradeStream.listenMessage()")
	// сразу создадим переменные, что бы их переиспользовать
	var err error
	var msg *pb.OrderTradeResponse
	for {
		select {
		case <-ctx.Done():
			log.Debug("listenOrderStream: ctx.Done(), выходим")
			return
		case <-s.closeCh:
			log.Debug("listenOrderStream c.closeCh")
			return
		default:
			msg, err = stream.Recv()
			if err != nil {
				// тест проверка статуса
				//st := status.Convert(err)
				//log.Debug("listenOrderStream", "status", st)
				// Проверка на конкретный код ошибки
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					log.Debug("listenOrderStream: Recv: контекст отменён (canceled), выходим")
					return
				}
				if err == io.EOF {
					log.Debug("listenOrderStream: Поток завершён")
					return
				} else {
					log.Error("listenOrderStream:  Ошибка чтения из потока", "err", err.Error())
					// вызовем реконнект
					s.Reconnect()
					return
				}
			}
			s.handleMessage(msg)

			msg.Reset()
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
			select {
			case s.orderChan <- *order:
				// успешно отправили
			default:
				log.Error("orderChan: канал переполнен, дропаем", slog.Any("order", order))
			}
		}
	}
}

// handleTrades обработка сделок
func (s *OrderTradeStream) handleTrades(trades []*pb.AccountTrade) {
	if trades != nil {
		// обработка
		for _, trade := range trades {
			select {
			case s.tradeChan <- *trade:
				// успешно отправили
			default:
				log.Error("tradeChan: канал переполнен, дропаем", slog.Any("trade", trade))
			}
		}
	}
}

// startOrderWorker Воркер для чтения канала ордеров
func (s *OrderTradeStream) startOrderWorker(ctx context.Context) {
	log.Debug("startOrderWorker")
	go func() {
		hasHandleOrder := s.onOrder != nil
		for {
			select {
			case <-ctx.Done():
				log.Debug("OrderWorker ctx.Done() = выход")
				return
			case data, ok := <-s.orderChan:
				if !ok {
					log.Debug("OrderWorker orderChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				if hasHandleOrder {
					s.onOrder(data)
				}
			}
		}
	}()

}

// startTradeWorker Воркер для чтения канала сделок
func (s *OrderTradeStream) startTradeWorker(ctx context.Context) {
	log.Debug("startTradeWorker")
	go func() {
		hasHandle := s.onTrade != nil
		for {
			select {
			case <-ctx.Done():
				log.Debug("tradeWorker ctx.Done() = выход")
				return
			case data, ok := <-s.tradeChan:
				if !ok {
					log.Debug("tradeWorker tradeChan закрыт = выход")
					return // канал закрыт
				}
				// пошлем получателю
				if hasHandle {
					s.onTrade(data)
				}
			}
		}
	}()

}
