package enum

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
	"wp-enum/pkg/data"
)

func TestEnumerateReturnsErrorWithInvalidEnumType(t *testing.T) {
	for i := 10; i < 15; i++ {
		res, err := Enumerate(data.EnumerationType(i), "test.com")
		if res != nil {
			t.Fatalf("expected nil result for invalid enumeration type: %d", i)
		}
		if err == nil {
			t.Fatalf("expected error for invalid enumeration type")
		}
	}
}

func jsonSuccess() http.Handler {
	resp := []data.ApiResponse{
		data.ApiResponse{Name: "admin", Id: 1},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json, _ := json.Marshal(resp)
		w.Write(json)
	})
}

func jsonFailure() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "whatever")
	})
}

func restSuccess() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		author, ok := r.URL.Query()["author"]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if "9" == author[0] {
			w.Header().Set("Location", "admin")
			w.WriteHeader(http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
}

func fakeJsonApiSuccessServer(addr string, handler http.Handler) io.Closer {
	closer, _ := listenAndServeWithClose(addr, handler)
	return closer
}

func listenAndServeWithClose(addr string, handler http.Handler) (io.Closer, error) {
	var (
		listener  net.Listener
		srvCloser io.Closer
		err       error
	)
	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
	}()

	srvCloser = listener
	return srvCloser, nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}