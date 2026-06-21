package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_logger "github.com/raizonw/todo-app/internal/core/logger"
	core_http_response "github.com/raizonw/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/raizonw/todo-app/internal/core/transport/http/utils"
)

type GetStatisticsResponse struct {
	TaskCreated               int
	TaskCompleted             int
	TaskCompletedRate         *float64
	TaskAverageCompletionTime *string
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	params, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, params.userID, params.from, params.to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
	}

	response := toDTOFromDomain(statistics)
	responseHandler.JSONResponse(response, http.StatusOK)

}

type queryParams struct {
	userID *int
	from   *time.Time
	to     *time.Time
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TaskAverageCompletionTime != nil {
		duration := statistics.TaskAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TaskCreated:               statistics.TaskCreated,
		TaskCompleted:             statistics.TaskCompleted,
		TaskCompletedRate:         statistics.TaskCompletedRate,
		TaskAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (queryParams, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_utils.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get `user_id` query param: %w", err)
	}

	from, err := core_http_utils.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get `from` query param: %w", err)
	}

	to, err := core_http_utils.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get `to` query param: %w", err)
	}

	return queryParams{
		userID: userID,
		from:   from,
		to:     to,
	}, nil
}
