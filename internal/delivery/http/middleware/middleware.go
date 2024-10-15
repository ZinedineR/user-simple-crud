package api

import "user-simple-crud/internal/delivery/http"

type Middleware struct {
	http.Handler
}
