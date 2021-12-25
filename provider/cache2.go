package provider

import (
	"errors"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
)

// Cached2 wraps a getter with a cache
type Cached2[T any] struct {
	mux     sync.Mutex
	clock   clock.Clock
	updated time.Time
	cache   time.Duration
	val     T
	err     error
}

// NewCached2 wraps a getter with a cache
func NewCached2[T any](g func() (T, error), cache time.Duration) func() (T, error) {
	c := &Cached2[T]{
		clock:  clock.New(),
		cache:  cache,
	}

	_ = bus.Subscribe(reset, c.reset)

	return func() (T, error) {
		c.mux.Lock()
		defer c.mux.Unlock()

		if c.mustUpdate() {
			c.val, c.err = g()
			c.updated = c.clock.Now()
		}

		return c.val, c.err
	}
}

func (c *Cached2[T]) reset() {
	c.mux.Lock()
	c.updated = time.Time{}
	c.mux.Unlock()
}

func (c *Cached2[T]) mustUpdate() bool {
	return c.clock.Since(c.updated) > c.cache || errors.Is(c.err, api.ErrMustRetry)
}

// // FloatGetter gets float value
// func (c *Cached2[T]) FloatGetter() func() (float64, error) {
// 	g, ok := c.getter.(func() (float64, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (float64, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(float64), c.err
// 	}
// }

// // IntGetter gets int value
// func (c *Cached2[T]) IntGetter() func() (int64, error) {
// 	g, ok := c.getter.(func() (int64, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (int64, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(int64), c.err
// 	}
// }

// // StringGetter gets string value
// func (c *Cached2[T]) StringGetter() func() (string, error) {
// 	g, ok := c.getter.(func() (string, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (string, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(string), c.err
// 	}
// }

// // BoolGetter gets bool value
// func (c *Cached2[T]) BoolGetter() func() (bool, error) {
// 	g, ok := c.getter.(func() (bool, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (bool, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(bool), c.err
// 	}
// }

// // DurationGetter gets time.Duration value
// func (c *Cached2[T]) DurationGetter() func() (time.Duration, error) {
// 	g, ok := c.getter.(func() (time.Duration, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (time.Duration, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(time.Duration), c.err
// 	}
// }

// // TimeGetter gets time.Time value
// func (c *Cached2[T]) TimeGetter() func() (time.Time, error) {
// 	g, ok := c.getter.(func() (time.Time, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (time.Time, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val.(time.Time), c.err
// 	}
// }

// // InterfaceGetter gets interface value
// func (c *Cached2[T]) InterfaceGetter() func() (interface{}, error) {
// 	g, ok := c.getter.(func() (interface{}, error))
// 	if !ok {
// 		log.FATAL.Fatalf("invalid type: %T", c.getter)
// 	}

// 	return func() (interface{}, error) {
// 		c.mux.Lock()
// 		defer c.mux.Unlock()

// 		if c.mustUpdate() {
// 			c.val, c.err = g()
// 			c.updated = c.clock.Now()
// 		}

// 		return c.val, c.err
// 	}
// }
