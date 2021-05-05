package middleware

import "net/http"

func AccessMiddleWare(next http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(writer, request)
	})
}