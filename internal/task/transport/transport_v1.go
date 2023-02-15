package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
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

	router.Methods("POST").Path("/v1/tasks").Handler(kithttp.NewServer(
		endpoints.CreateTaskEndpoint,
		decodeCreateTaskRequest,
		helpers.EncodeResponse,
		options...,
	))

	router.Methods("GET").Path("/v1/tasks/{id}").Handler(kithttp.NewServer(
		endpoints.GetTaskEndpoint,
		decodeGetTaskRequest,
		helpers.EncodeResponse,
		options...,
	))

	router.Methods("GET").Path("/v1/tasks").Handler(kithttp.NewServer(
		endpoints.ListTasksEndpoint,
		decodeListTasksRequest,
		helpers.EncodeResponse,
		options...,
	))
}

func decodeCreateTaskRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.TaskWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return createTaskRequest{
		Task: &domain.Task{
			Method: body.Method,
			URL:    body.URL,
		},
	}, nil
}

func decodeGetTaskRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	vars := mux.Vars(request)
	id, err := helpers.ExtractInt64Route(vars, "id")
	if err != nil {
		return nil, err
	}
	return getTaskRequest{
		TaskID: domain.TaskID(id),
	}, nil
}

func decodeListTasksRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	var (
		query = request.URL.Query()
		page  = helpers.ParsePageOrGetDefault(query.Get("page"))
		size  = helpers.ParseSizeOrGetDefault(query.Get("size"))
	)

	return listTaskRequest{
		Criteria: domain.TaskSearchCriteria{
			Page: domain.PageRequest{
				Offset: helpers.PageOffset(page, size),
				Size:   int(size),
			},
		},
	}, nil
}
