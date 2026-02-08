// Package aspect. context_pool provides pooling objects for Context.
package aspect

import "sync"

// -------------------------------------------- Types, Variables & Constants --------------------------------------------

var enableContextPooling = true

var contextPool = sync.Pool{
	New: func() any {
		return &Context{
			Metadata: make(map[string]any, 8),
			Results:  make([]any, 0, 2),
		}
	},
}

// -------------------------------------------- Private Functions --------------------------------------------

// acquireContext gets a Context from the pool and initializes it.
func acquireContext(functionName FuncKey, args ...any) *Context {
	c := contextPool.Get().(*Context)

	c.FunctionName = functionName
	c.Args = args
	c.Results = c.Results[:0]
	c.Error = nil
	c.PanicValue = nil
	c.Skipped = false

	for k := range c.Metadata {
		delete(c.Metadata, k)
	}

	return c
}

// releaseContext resets and returns the Context to the pool.
func releaseContext(c *Context) {
	if cap(c.Results) > 16 {
		c.Results = make([]any, 0, 2)
	}
	if len(c.Metadata) > 32 {
		c.Metadata = make(map[string]any, 8)
	}

	c.Reset()
	contextPool.Put(c)
}
