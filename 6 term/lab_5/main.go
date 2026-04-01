package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed script.lua
var script string

const (
	limit    = 100
	duration = time.Minute
)

type Limiter interface {
	Limit(ip net.IP) bool
}

type ListLimiter struct{ rdb *redis.Client }

func NewListLimiter(rdb *redis.Client) *ListLimiter {
	return &ListLimiter{rdb: rdb}
}
func (l *ListLimiter) Limit(ip net.IP) bool {
	key := fmt.Sprintf("rate:%s", ip.String())
	val, err := redis.NewScript(script).Run(context.Background(), l.rdb, []string{key}).Result()
	must(err)
	return val.(int64) == 1
}

type CounterLimiter struct{ rdb *redis.Client }

func NewCounterLimiter(rdb *redis.Client) *CounterLimiter { return &CounterLimiter{rdb: rdb} }
func (l *CounterLimiter) Limit(ip net.IP) bool {
	key := fmt.Sprintf("rate:%s", ip.String())
	count, err := l.rdb.Incr(context.Background(), key).Result()
	must(err)
	if count == 1 {
		l.rdb.Expire(context.Background(), key, duration)
	} else if count > limit {
		return false
	}
	return true
}

type RestHandler struct{ limiter Limiter }

func NewRestHandler(limiter Limiter) *RestHandler {
	return &RestHandler{limiter: limiter}
}
func (h *RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	must(err)
	ip := net.ParseIP(host)
	if h.limiter.Limit(ip) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode("Rate limit exceeded")
		return
	}
	json.NewEncoder(w).Encode("Get yourself a cookie!")
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	mux := http.NewServeMux()
	handler := &RestHandler{NewListLimiter(rdb)}
	mux.Handle("GET /cookie", handler)

	go func() {
		for {
			values, err := rdb.BRPop(context.Background(), 0, "queue:email").Result()
			if err != nil {
				if err != redis.Nil {
					must(err)
				}
			}
			email := values[1]
			slog.Info("Let's just imagine I sent message to email", "email", email)
		}
	}()

	go func() {
		for {
			values, err := rdb.BRPop(context.Background(), 0, "queue:logs").Result()
			if err != nil {
				if err != redis.Nil {
					must(err)
				}
			}
			log := values[1]
			slog.Info(log)
		}
	}()

	logsTicker := time.NewTicker(time.Millisecond * time.Duration(rand.Int31n(1000)))
	emailsTicker := time.NewTicker(time.Millisecond * time.Duration(rand.Int31n(1000)))
	go func() {
		for {
			select {
			case <-logsTicker.C:
				rdb.RPush(context.Background(), "queue:logs", randLog())
			case <-emailsTicker.C:
				rdb.RPush(context.Background(), "queue:email", randEmail())
			}
		}
	}()

	must(http.ListenAndServe(":8080", mux))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randEmail() string {
	n := rand.Intn(1<<5 - 1)
	if n < 10 {
		n += 10
	}
	b := make([]byte, 2*n+2)
	i := 0
	for i < n {
		b[i] = letters[rand.Intn(len(letters))]
		i++
	}
	b[i] = '@'
	i++
	for i < 2*n-2 {
		b[i] = letters[rand.Intn(len(letters))]
		i++
	}
	b[i], b[i+1], b[i+2], b[i+3] = '.', 'c', 'o', 'm'
	return string(b)
}

func randLog() string {
	n := rand.Intn(1<<6 - 1)
	if n < 10 {
		n += 10
	}
	b := make([]byte, n)
	for i := range n {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
