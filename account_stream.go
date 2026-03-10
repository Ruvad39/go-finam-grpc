package finam

import (
	"context"
	"math/rand"
	"time"

	accounts_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/accounts"
	"google.golang.org/grpc"
)

// AccountStream
// стрим информации  по заданному счету
type AccountStream struct {
	ctx            context.Context
	cancel         context.CancelFunc
	done           chan struct{}
	retryDelay     time.Duration
	client         *Client
	AccountService accounts_service.AccountsServiceClient
	onAccount      func(*accounts_service.GetAccountResponse) // callback функция
	dataChannel    chan *accounts_service.GetAccountResponse  // канал для данных (альтернатива callback функции)
	accountId      string                                     // Номер счета для подписки
	running        bool                                       // признак что уже запустили в работу
}

// NewOrderStream
// создадим стрим информации по заданному счету
//
// на входе callback функция для обработки данных
func (c *Client) NewAccountStream(parent context.Context, accountId string, callback func(*accounts_service.GetAccountResponse)) *AccountStream {
	ctx, cancel := context.WithCancel(parent)
	s := &AccountStream{
		ctx:            ctx,
		cancel:         cancel,
		client:         c,
		AccountService: accounts_service.NewAccountsServiceClient(c.conn),
		done:           make(chan struct{}),
		retryDelay:     initialDelay,
		accountId:      accountId,
		onAccount:      callback,
	}
	s.running = true
	go s.run()
	return s
}

// NewAccountStreamWithCallback
// создадим стрим информации по заданному счету
//
// на входе callback функция для обработки данных
//
// стрим НЕ запускается по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewAccountStreamWithCallback(parent context.Context, accountId string, callback func(*accounts_service.GetAccountResponse)) *AccountStream {
	ctx, cancel := context.WithCancel(parent)
	s := &AccountStream{
		ctx:            ctx,
		cancel:         cancel,
		client:         c,
		AccountService: accounts_service.NewAccountsServiceClient(c.conn),
		done:           make(chan struct{}),
		retryDelay:     initialDelay,
		accountId:      accountId,
		onAccount:      callback,
	}

	return s
}

// NewAccountStreamStreamWithChannel
// создадим стрим по заданному счяету
// возвращаем канал с данными
// стрим НЕ запускам по умолчанию => Нужно выполнить метод  Start()
func (c *Client) NewAccountStreamStreamWithChannel(parent context.Context, accountId string) (*AccountStream, chan *accounts_service.GetAccountResponse) {
	ctx, cancel := context.WithCancel(parent)
	s := &AccountStream{
		ctx:            ctx,
		cancel:         cancel,
		client:         c,
		AccountService: accounts_service.NewAccountsServiceClient(c.conn),
		accountId:      accountId,
		done:           make(chan struct{}),
		retryDelay:     initialDelay,
		dataChannel:    make(chan *accounts_service.GetAccountResponse, 100),
	}

	return s, s.dataChannel
}

// Start
func (s *AccountStream) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
}
func (s *AccountStream) Close() {
	s.cancel()
	<-s.done // дождаться завершения run()
}

func (s *AccountStream) run() {
	defer func() {
		log.Debug("[AccountStream] exit run()", "accountId", s.accountId)
		close(s.done)
	}()
	for {
		err := s.subscribeAndListen()
		// выход без ошибки
		if err == nil {
			return
		}
		log.Error("[AccountStream]", "accountId", s.accountId, "err", err.Error())
		// Проверка на конкретный код ошибки
		if shouldTerminate(err) {
			return
		}
		log.Warn("[AccountStream] start reconnect", "accountId", s.accountId, "retryDelay", s.retryDelay)
		select {
		case <-s.ctx.Done():
			log.Debug("[AccountStream] context cancelled, stopping", "accountId", s.accountId)
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
func (s *AccountStream) subscribeAndListen() error {
	log.Debug("[AccountStream].subscribeAndListen", "accountId", s.accountId)

	// создаем стрим
	stream, err := s.AccountService.SubscribeAccount(s.ctx, &accounts_service.GetAccountRequest{AccountId: s.accountId})
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

func (s *AccountStream) listen(ctx context.Context, stream grpc.ServerStreamingClient[accounts_service.GetAccountResponse]) error {
	log.Debug("[AccountStream].listen", "accountId", s.accountId)
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
func (s *AccountStream) handleMessage(msg *accounts_service.GetAccountResponse) {
	if msg != nil {
		// обработка
		// отправим в канал
		if s.dataChannel != nil {
			s.dataChannel <- msg
		}
		// отправил в callback
		if s.onAccount != nil {
			s.onAccount(msg)
		}

	}

}
