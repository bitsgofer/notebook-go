package webserver

import "net/http"

func NewDevelopmentServer(addr string, fileRoot string) (*http.Server, error) {
	handler := http.FileServer(http.Dir(fileRoot))
	srv := http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: nil,
		// ReadTimeout time.Duration
		// ReadHeaderTimeout time.Duration
		// WriteTimeout time.Duration
		// IdleTimeout time.Duration
		// MaxHeaderBytes int
		// TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
		// ConnState func(net.Conn, ConnState)
		// ErrorLog *log.Logger
		// BaseContext func(net.Listener) context.Context
		// ConnContext func(ctx context.Context, c net.Conn) context.Context
	}

	return &srv, nil
}
