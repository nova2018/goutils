package goutils

import (
	"context"
	"sync"
)

type Stopped interface {
	FreeAble
	IsStopped() bool
	Stop()
}

type stopped struct {
	isStopped bool
	lock      sync.RWMutex
	p         Pool
}

var (
	_ Stopped = &stopped{}
)

func (c *stopped) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.isStopped = true
}

func (c *stopped) IsStopped() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.isStopped
}

func (c *stopped) Free() {
	c.p.Put(c)
}

func (c *stopped) SetPool(p Pool) {
	c.p = p
}

func (c *stopped) Reset() {
	c.lock = sync.RWMutex{}
	c.isStopped = false
}

var (
	_stoppedPool = NewPool(func() interface{} {
		return &stopped{}
	})
)

func AcquireStopped() Stopped {
	return _stoppedPool.Get().(Stopped)
}

type SimpleProducer[T any] func(ctx context.Context, chanProducer chan<- T, stop Stopped) error

type SimpleConsumer[T1, T2 any] func(ctx context.Context, chanProducer <-chan T1, chanConsumer chan<- T2, stop Stopped) error

type SimpleWait[T1, T2 any] func(ctx context.Context, chanConsumer <-chan T1, stop Stopped) (T2, error)

func SimpleParallel[T1, T2, T3 any](ctx context.Context, concurrent int, producer SimpleProducer[T1], consumer SimpleConsumer[T1, T2], wait SimpleWait[T2, T3]) (T3, []error) {
	listErr := make([]error, 0, concurrent+2)

	stop := AcquireStopped()
	defer stop.Free()

	wg := &sync.WaitGroup{}

	chProducer := make(chan T1, 2*concurrent)
	wg.Add(1)
	GoWithContext(ctx, func(ctx context.Context) {
		defer wg.Done()
		defer close(chProducer)
		err := producer(ctx, chProducer, stop)
		if err != nil {
			listErr = append(listErr, err)
		}
	})

	chConsumer := make(chan T2, 4*concurrent)
	wg.Add(1)
	GoWithContext(ctx, func(ctx context.Context) {
		defer wg.Done()
		defer close(chConsumer)
		wgConsumer := &sync.WaitGroup{}
		for i := 0; i < concurrent; i++ {
			wgConsumer.Add(1)
			GoWithContext(ctx, func(ctx context.Context) {
				defer wgConsumer.Done()
				err := consumer(ctx, chProducer, chConsumer, stop)
				if err != nil {
					listErr = append(listErr, err)
				}
			})
		}
		wgConsumer.Wait()
	})

	wg.Add(1)
	var result T3
	var err error
	GoWithContext(ctx, func(ctx context.Context) {
		defer wg.Done()
		result, err = wait(ctx, chConsumer, stop)
		if err != nil {
			listErr = append(listErr, err)
		}
	})

	wg.Wait()

	return result, listErr
}
