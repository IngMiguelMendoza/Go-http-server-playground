package main

import (
	"log"
	"net/http"
)

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		next.ServeHTTP(resp, req)
	})
}
