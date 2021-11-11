package rmq

import (
	"context"
	"sync"
)

// Handler interface for handlers.
type Handler interface {
	Handle(ctx context.Context, body []byte) error
}

// HandlerFunc call anonymous function.
type HandlerFunc func(ctx context.Context, body []byte) error

// Handle for run function.
func (f HandlerFunc) Handle(ctx context.Context, body []byte) error {
	return f(ctx, body)
}

// Router contains handlers.
type Router struct {
	handlers map[interface{}]Handler
	mu       sync.RWMutex
}

// NewRouther return new router.
func NewRouther(keyValues ...interface{}) (*Router, error) {
	var handlers map[interface{}]Handler
	if len(keyValues)%2 == 0 {
		handlers = make(map[interface{}]Handler, len(keyValues)/2)
		for i := 0; i < len(keyValues); i += 2 {
			key := keyValues[i]
			value, ok := keyValues[i+1].(Handler)
			if !ok {
				return nil, ErrWrongValue
			}
			handlers[key] = value
		}
	} else {
		return nil, ErrNotEnoughArguments
	}

	return &Router{
		handlers: handlers,
	}, nil
}

// Add add handler.
func (h *Router) Add(name interface{}, f Handler) *Router {
	h.mu.Lock()
	h.handlers[name] = f
	h.mu.Unlock()
	return h
}

// Remove remove handler.
func (h *Router) Remove(name interface{}) *Router {
	h.mu.Lock()
	delete(h.handlers, name)
	h.mu.Unlock()
	return h
}

// Get get handler.
func (h *Router) Get(name interface{}) Handler {
	h.mu.RLock()
	f := h.handlers[name]
	h.mu.RUnlock()
	return f
}

// Call call handler.
func (h *Router) Call(ctx context.Context, name interface{}, body []byte) error {
	c := h.Get(name)
	if c == nil {
		return ErrNotFoundMethod
	}
	return c.Handle(ctx, body)
}
