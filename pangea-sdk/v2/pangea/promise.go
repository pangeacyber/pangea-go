package pangea

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Promise[T any] struct {
	res     *PangeaResponse[T]
	err     error
	ready   atomic.Bool
	mu      sync.Mutex
	cancel  context.CancelCauseFunc // To close promise if desired
	execute func()
}

func NewPromise[T any, R any](f func(ctx context.Context, input R) (*PangeaResponse[T], error), ctx context.Context, input R) *Promise[T] {
	ctx, cancel := context.WithCancelCause(ctx)

	p := &Promise[T]{
		ready:  atomic.Bool{},
		err:    nil,
		res:    nil,
		cancel: cancel,
	}
	p.ready.Store(false)
	p.mu.Lock()

	go func() {
		defer p.finally()
		p.res, p.err = f(ctx, input)
	}()

	return p
}

func (p *Promise[T]) finally() {
	if r := recover(); r != nil {
		p.res = nil
		p.err = fmt.Errorf("panic: %v", r)
	}

	p.ready.Store(true)
	p.mu.Unlock()
}

func (p *Promise[T]) Wait() {
	if p.ready.Load() {
		return
	}

	p.mu.Lock()
	p.mu.Unlock()
	return
}

func (p *Promise[T]) Cancel() {
	p.cancel(errors.New("Canceled from promise"))
}

func (p *Promise[T]) IsReady() bool {
	return p.ready.Load()
}

func (p *Promise[T]) Get() (*PangeaResponse[T], error) {
	p.Wait()
	return p.res, p.err
}
