package goutils

import (
	"context"
	"runtime"
)

func GoWithContext(ctx context.Context, fn func(context.Context)) <-chan struct{} {
	return GoWithContextHandler(ctx, fn, nil)
}

func GoWithContextHandler(ctx context.Context, fn func(context.Context), h recoveryHandle) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer recoveryWithHandler(ctx, h)
		defer close(done)
		fn(ctx)
	}()
	return done
}

func recoveryWithHandler(ctx context.Context, h recoveryHandle) {
	if h == nil {
		h = GetRecoveryHandle(ctx)
	}
	if h == nil {
		return
	}

	if err := recover(); err != nil {
		// panic
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)

		// filter
		{
			i := len(buf) - 1
			for ; i >= 0; i-- {
				if buf[i] != 0 {
					break
				}
			}
			buf = buf[0 : i+1]
		}

		if h != nil {
			h(ctx, err, buf)
		}
	}
}
