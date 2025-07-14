package router

import (
	"net/http"
	"strings"
)

type routeKey struct {
	Method string
	Path   string
}

var routes = make(map[routeKey]http.HandlerFunc)

// Зарегистрировать маршрут
func Handle(method string, path string, handler http.HandlerFunc) {
	routes[routeKey{Method: strings.ToUpper(method), Path: path}] = handler
}

// Роутер как основной обработчик
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := routes[routeKey{Method: r.Method, Path: r.URL.Path}]; ok {
		handler(w, r)
	} else {
		http.NotFound(w, r)
	}
}
