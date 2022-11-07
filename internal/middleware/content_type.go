package middleware

import "net/http"

type Content struct {
	handler http.Handler
}

func (c *Content)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c.handler.ServeHTTP(w,r)
}

func NewContentMiddleware(handlerToWrap http.Handler) *Content {
    return &Content{handlerToWrap}
}