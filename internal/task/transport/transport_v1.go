package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/helpers"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func RegisterRoutersV1(router *mux.Router, endpoints Endpoints, logger log.Logger) {

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	router.Methods("POST").Path("/tasks").Handler(kithttp.NewServer(
		endpoints.CreateTaskEndpoint,
		decodeCreateTaskRequest,
		helpers.EncodeResponse,
		options...,
	))

	router.Methods("GET").Path("/tasks/{id}").Handler(kithttp.NewServer(
		endpoints.GetTaskEndpoint,
		decodeGetTaskRequest,
		helpers.EncodeResponse,
		options...,
	))

	router.Methods("GET").Path("/tasks").Handler(kithttp.NewServer(
		endpoints.ListTasksEndpoint,
		decodeListTasksRequest,
		helpers.EncodeResponse,
		options...,
	))
}

func decodeCreateTaskRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Task); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetTaskRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeListTasksRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}
