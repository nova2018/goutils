package goutils

import "context"

type recoveryHandle func(context.Context, any, []byte)

const (
	ContextRecoveryField = "goutils.recovery"
)

func GetRecoveryHandle(ctx context.Context) recoveryHandle {
	if v := ctx.Value(ContextRecoveryField); v != nil {
		if h, ok := v.(recoveryHandle); ok {
			return h
		}
	}
	return nil
}
