package error

import (
	"context"
	"encoding/json"
	"net/http"
)

// EncodeErrorResponse - encodes error response
func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	//case order.ErrOrderNotFound:
	//	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
