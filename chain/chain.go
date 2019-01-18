package main

type Context struct {
	Nums []int
}

type Middleware interface {
	Do(ctx *Context, chain Chain) error
}

type Chain interface {
	Next(ctx *Context) error
}

type RealMiddlewareChain struct {
	middlewares []Middleware
	index       int
}

func (c *RealMiddlewareChain) Next(ctx *Context) error {
	if c.index >= len(c.middlewares) {
		return nil
	}
	next := &RealMiddlewareChain{middlewares: c.middlewares, index: c.index + 1}
	middleware := c.middlewares[c.index]
	err := middleware.Do(ctx, next)
	return err
}

func (c *RealMiddlewareChain) add(middleware Middleware) {
	c.middlewares = append(c.middlewares, middleware)
}

type FunctionMiddleware struct {
	Fn func(ctx *Context, chain Chain) error
}

func (m *FunctionMiddleware) Do(ctx *Context, chain Chain) error {
	return m.Fn(ctx, chain)
}

func main() {
	chain := &RealMiddlewareChain{}
	chain.add(&FunctionMiddleware{Fn: func(ctx *Context, chain Chain) error {
		ctx.Nums = append(ctx.Nums, 1)
		chain.Next(ctx)
		return nil
	}})
	chain.add(&FunctionMiddleware{Fn: func(ctx *Context, chain Chain) error {
		ctx.Nums = append(ctx.Nums, 2)
		chain.Next(ctx)

		return nil
	}})
	chain.add(&FunctionMiddleware{Fn: func(ctx *Context, chain Chain) error {
		ctx.Nums = append(ctx.Nums, 3)
		chain.Next(ctx)

		return nil
	}})
	chain.add(&FunctionMiddleware{Fn: func(ctx *Context, chain Chain) error {
		ctx.Nums = append(ctx.Nums, 4)
		chain.Next(ctx)
		return nil
	}})
	chain.Next(&Context{})
}
