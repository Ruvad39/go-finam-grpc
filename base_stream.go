/*

onOrder func(ctx context.Context)

пример метода NewOrderStream
stream := &OrderStream{
		BaseStream: BaseStream{
			reconnectCh: make(chan struct{}, 1), // Буферизированный канал
			closeCh:     make(chan struct{}),    // Небуферизированный
		},

*/

package finam

import (
	"context"
	"fmt"
	"time"
)

const (
	initialDelay = 2 * time.Second
	maxDelay     = 50 * time.Second
)

// BaseStream содержит общую логику для всех стримов.
// Должен встраиваться в конкретные реализации.
type BaseStream struct {
	Name        string                          // имя стрима (для отладки)
	ConnectFunc func(ctx context.Context) error // функция которая запускается при коннеете-реконнекте
	reconnectCh chan struct{}                   // Сигнальный канал для необходимости реконекта
	closeCh     chan struct{}                   // Сигнальный канал для закрытия коннекта
	workerFunc  []func(ctx context.Context)     // Список функций, которые надо запутсить при старте
}

// Close закрываем канал closeCh
func (s *BaseStream) Close() error {
	close(s.closeCh)
	return nil
}

// Start запуск стрима
func (s *BaseStream) Start(ctx context.Context) error {
	log.Info("BaseStream.Start()", "stream name", s.Name)
	if s.ConnectFunc == nil {
		return fmt.Errorf("ConnectFunc is nil")
	}

	// запустим воркеры для чтения канала данных
	s.startWorkers(ctx)

	err := s.ConnectFunc(ctx)
	if err != nil {
		return err
	}

	// запустим программу реконнекта в отдельной горутине
	go s.reconnector(ctx)

	return nil
}

// AddWorker добавим функцию (worker)
func (s *BaseStream) AddWorker(cb func(ctx context.Context)) {
	s.workerFunc = append(s.workerFunc, cb)
}

// startWorkers запустим методы
func (s *BaseStream) startWorkers(ctx context.Context) {
	log.Debug("BaseStream.startWorkers", "stream name", s.Name)
	for n, cb := range s.workerFunc {
		log.Debug("BaseStream.startWorkers", "n", n, "func", cb)
		cb(ctx)
	}
}

// Reconnect в сигнальный канал рекконета пошлем сообщение
func (s *BaseStream) Reconnect() {
	log.Info("BaseStream.Reconnect()", "stream name", s.Name)
	select {
	case s.reconnectCh <- struct{}{}:
	default:
	}
}

// reconnector ждет команду на переподключение
// работает в отдельно потоке
// запускает tryReconnect
func (s *BaseStream) reconnector(ctx context.Context) {
	log.Info("BaseStream.reconnector()", "stream name", s.Name)
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closeCh:
			return
		case <-s.reconnectCh:
			s.tryReconnect(ctx)
		}
	}
}

// TryReconnect tries to reconnect .
func (s *BaseStream) tryReconnect(ctx context.Context) {
	log.Info("BaseStream.tryReconnect()", "stream name", s.Name)
	delay := initialDelay
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closeCh:
			return
		case <-time.After(delay):
			if err := s.ConnectFunc(ctx); err != nil {
				if delay >= maxDelay {
					delay = maxDelay
				} else {
					delay *= 2
				}
				log.Error("BaseStream: re-connect error, try to reconnect later", "stream name", s.Name, "delay", delay, "err", err.Error())
				continue
			}
			log.Warn("BaseStream: reconnect success", "stream name", s.Name)
			return
		}
	}
}
