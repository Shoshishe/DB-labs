package ioc

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

var UseHttpMux = func() *http.ServeMux {
	return http.NewServeMux()
}

var UseHttpServer = func(ctx context.Context) *http.Server {
	return provider(func() *http.Server {
		slog.Info("Starting server on :8080 port")
		server := &http.Server{
			Addr:    ":8080",
			Handler: UseHttpMux(),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			}}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
			slog.Info("Server starting on :8080.")
		}()
		return server
	})()
}

func UseRoutes() {}
