package main

import (
	"net/http"
)

func main() {
	for range 101 {
		resp, _ := http.Get("http://localhost:8080/cookie")
		print(resp.StatusCode)
	}
}
