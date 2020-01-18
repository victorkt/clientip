package clientip

import (
	"context"
	"net"
	"reflect"
	"testing"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ctx        context.Context
		expectedIP net.IP
	}{
		{
			name:       "returns the correct ip when set",
			ctx:        toContext(context.Background(), net.ParseIP("45.0.0.40")),
			expectedIP: net.ParseIP("45.0.0.40"),
		},
		{
			name:       "returns nil ip when value is not net.IP",
			ctx:        context.WithValue(context.Background(), ipCtxKey{}, "45.0.0.40"),
			expectedIP: nil,
		},
		{
			name:       "returns nil ip when not set",
			ctx:        context.Background(),
			expectedIP: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ip := FromContext(tt.ctx)
			if !reflect.DeepEqual(tt.expectedIP, ip) {
				t.Errorf("expected %s to equal %s", tt.expectedIP, ip)
			}
		})
	}
}
