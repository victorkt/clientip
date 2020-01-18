package clientip

import (
	"context"
	"net"
)

type ipCtxKey struct{}

// FromContext returns the client IP address stored in the context.
// The value will only be set in context when using the middleware.
func FromContext(ctx context.Context) net.IP {
	ip, ok := ctx.Value(ipCtxKey{}).(net.IP)
	if !ok {
		return nil
	}
	return ip
}

func toContext(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, ipCtxKey{}, ip)
}
