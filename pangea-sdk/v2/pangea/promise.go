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
	running atomic.Bool
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
	p.running.Store(false)
	p.mu.Lock()

	p.execute = func() {
		defer p.finally()
		p.res, p.err = f(ctx, input)
	}

	return p
}

func (p *Promise[T]) Execute() {
	if p.running.Load() {
		return
	}
	p.running.Store(true)
	p.execute()
}

func (p *Promise[T]) finally() {
	if r := recover(); r != nil {
		p.res = nil
		p.err = fmt.Errorf("panic: %v", r)
	}

	p.ready.Store(true)
	p.running.Store(false)
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

func (p *Promise[T]) IsRunning() bool {
	return p.running.Load()
}

func (p *Promise[T]) Get() (*PangeaResponse[T], error) {
	p.Wait()
	return p.res, p.err
}

// We should remember that if we kill context we are going to kill go rutine.
// Do not defer cancel right after calling CallAsync
// func CallAsync[T any, R any](w *Worker, f func(ctx context.Context, input R) (*PangeaResponse[T], error), ctx context.Context, input R) *Promise[T] {
// 	p := newPromise(f, ctx, input)
// 	w.Run(p.Execute)
// 	return p
// }

type Worker struct {
	maxThreads uint
	inprogress chan int
	functions  chan func()
	wg         sync.WaitGroup
	stop       chan bool // To stop runner go rutine
	closed     bool
}

func NewWorker(maxThreads uint) *Worker {
	if maxThreads == 0 {
		maxThreads = 1
	}

	w := Worker{
		maxThreads: maxThreads,
		functions:  make(chan func()),
		inprogress: make(chan int, maxThreads),
		stop:       make(chan bool),
		closed:     false,
	}
	go w.runner()
	return &w
}

func (w *Worker) MaxThreads() uint {
	return w.maxThreads
}

func (w *Worker) ThreadsRunning() uint {
	return uint(len(w.inprogress))
}

func (w *Worker) ThreadsPending() uint {
	return uint(len(w.functions))
}

func (w *Worker) Run(f func()) {
	if w.closed {
		return
	}
	w.functions <- f
}

func (w *Worker) runner() {
	var f func()
	defer func() {
		if r := recover(); r != nil {
			go w.runner()
		} else {
			w.closed = true
		}
	}()

	for {
		select {
		case <-w.stop:
			return // Stop the worker
		case f = <-w.functions:
			// do nothing
		}

		select {
		case <-w.stop:
			return // Stop the worker
		case w.inprogress <- 1:
			// do nothing
		}
		w.run(f)
	}
}

func (w *Worker) finish() {
	if r := recover(); r != nil {
	}
	w.wg.Done()
	<-w.inprogress
}

func (w *Worker) run(f func()) {
	w.wg.Add(1)
	go func() {
		defer w.finish()
		f()
	}()
}

func (w *Worker) Close() {
	w.stop <- true
	w.wg.Wait()
}
