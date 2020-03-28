package clientip

import (
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestFromRequest(t *testing.T) {
	t.Parallel()

	createRequest := func(remoteAddr string, headers ...string) *http.Request {
		h := make(http.Header, 0)
		if len(headers) == 2 {
			h.Set(headers[0], headers[1])
		}

		return &http.Request{
			RemoteAddr: remoteAddr,
			Header:     h,
		}
	}

	tests := []struct {
		name       string
		req        *http.Request
		expectedIP net.IP
	}{
		{
			name:       "returns the value of x-client-ip",
			req:        createRequest("45.0.0.40", "x-client-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the value of x-client-ip",
			req:        createRequest("45.0.0.40", "x-client-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the first value of x-forwarded-for",
			req:        createRequest("45.0.0.40", "x-forwarded-for", "129.78.138.66, 129.78.64.103, 129.78.64.105"),
			expectedIP: net.ParseIP("129.78.138.66"),
		},
		{
			name:       "returns the first value of x-forwarded-for with ipv6",
			req:        createRequest("45.0.0.40", "x-forwarded-for", "2001:0db8:0123:4567:89ab:cdef:1234:5678, 129.78.64.103, 129.78.64.105"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the first valid IP value of x-forwarded-for",
			req:        createRequest("45.0.0.40", "x-forwarded-for", "unknown, 129.78.64.103, 129.78.64.105"),
			expectedIP: net.ParseIP("129.78.64.103"),
		},
		{
			name:       "returns the correct IP value of x-forwarded-for with port",
			req:        createRequest("45.0.0.40", "x-forwarded-for", "129.78.138.66:12345, 129.78.64.103, 129.78.64.105"),
			expectedIP: net.ParseIP("129.78.138.66"),
		},
		{
			name:       "returns the correct IP value of x-forwarded-for with port",
			req:        createRequest("45.0.0.40", "x-forwarded-for", "[2001:0db8:0123:4567:89ab:cdef:1234:5678]:12345, 129.78.64.103, 129.78.64.105"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of cf-connecting-ip",
			req:        createRequest("45.0.0.40", "cf-connecting-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of cf-connecting-ip",
			req:        createRequest("45.0.0.40", "cf-connecting-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of fastly-client-ip",
			req:        createRequest("45.0.0.40", "fastly-client-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of fastly-client-ip",
			req:        createRequest("45.0.0.40", "fastly-client-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of true-client-ip",
			req:        createRequest("45.0.0.40", "true-client-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of true-client-ip",
			req:        createRequest("45.0.0.40", "true-client-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of x-real-ip",
			req:        createRequest("45.0.0.40", "x-real-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of x-real-ip",
			req:        createRequest("45.0.0.40", "x-real-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of x-cluster-client-ip",
			req:        createRequest("45.0.0.40", "x-cluster-client-ip", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of x-cluster-client-ip",
			req:        createRequest("45.0.0.40", "x-cluster-client-ip", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of x-forwarded",
			req:        createRequest("45.0.0.40", "x-forwarded", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of x-forwarded",
			req:        createRequest("45.0.0.40", "x-forwarded", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the value of forwarded-for",
			req:        createRequest("45.0.0.40", "forwarded-for", "45.9.248.40"),
			expectedIP: net.ParseIP("45.9.248.40"),
		},
		{
			name:       "returns the ipv6 value of forwarded-for",
			req:        createRequest("45.0.0.40", "forwarded-for", "2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the correct value of request.RemoteAddr when it contains the port",
			req:        createRequest("45.0.0.40:8080"),
			expectedIP: net.ParseIP("45.0.0.40"),
		},
		{
			name:       "returns the correct value of request.RemoteAddr when it doesn't contain the port",
			req:        createRequest("45.0.0.40"),
			expectedIP: net.ParseIP("45.0.0.40"),
		},
		{
			name:       "returns the correct ipv6 value of request.RemoteAddr when it contains the port",
			req:        createRequest("[2001:0db8:0123:4567:89ab:cdef:1234:5678]:8080"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns the correct ipv6 value of request.RemoteAddr when it doesn't contain the port",
			req:        createRequest("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
			expectedIP: net.ParseIP("2001:0db8:0123:4567:89ab:cdef:1234:5678"),
		},
		{
			name:       "returns nil when no valid IP was found",
			req:        createRequest(""),
			expectedIP: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ip := FromRequest(tt.req)
			if !reflect.DeepEqual(tt.expectedIP, ip) {
				t.Errorf("expected %s to equal %s", ip, tt.expectedIP)
			}
		})
	}
}
