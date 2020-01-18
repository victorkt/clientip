package clientip

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	var handlerCalled bool
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		ip := FromContext(r.Context())
		expectedIP := net.ParseIP("45.0.0.40")
		if !reflect.DeepEqual(expectedIP, ip) {
			t.Errorf("expected client IP from request to equal %s, got: %s", expectedIP, ip)
		}
		_, _ = fmt.Fprintln(w, "hello world")
	})

	srv := httptest.NewServer(Middleware(fn))
	defer srv.Close()

	req, err := http.NewRequest("POST", srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("x-forwarded-for", "45.0.0.40")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected response http code %d, got: %d", http.StatusOK, res.StatusCode)
	}

	if !handlerCalled {
		t.Error("expected test handler to have been called but didn't")
	}
}
