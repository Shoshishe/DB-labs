package ioc

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

var UseHttpMux = provider(func() *http.ServeMux {
	return http.NewServeMux()
})

func useCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, Cookie, withCredentials")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var UseHttpServer = func(ctx context.Context) *http.Server {
	return provider(func() *http.Server {
		slog.Info("Starting server on :8080 port")
		server := &http.Server{
			Addr:    ":8080",
			Handler: useCors(UseHttpMux()),
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
