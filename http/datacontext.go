package http

import (
	"context"
	"net/http"
)

const (
	dataContextKey = "data"
)

// DataContext type of data in request context
type DataContext map[string]interface{}

// DataContextMiddleware add data to request context
func DataContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		dataContext := make(DataContext)
		ctx := context.WithValue(request.Context(), dataContextKey, &dataContext)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

// GetDataContext return data from request context
func GetDataContext(request *http.Request) *DataContext {
	return request.Context().Value(dataContextKey).(*DataContext)
}
