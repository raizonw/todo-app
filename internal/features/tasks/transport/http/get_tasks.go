package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/raizonw/todo-app/internal/core/logger"
	core_http_response "github.com/raizonw/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/raizonw/todo-app/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Получить список задач
// @Description Получить список задач с фильтрацией по user_id и поддержкой пагинации (limit/offset)
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_id query int false "ID пользователя-автора задач"
// @Param limit query int false "Лимит количества задач"
// @Param offset query int false "Смещение относительно начала"
// @Success 200 {object} GetTasksResponse "Список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID, limit, offset query params",
		)

		return
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(tasksDTOFromDomain(taskDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	user_id, err := core_http_utils.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `user_id` query param: %w", err)
	}

	limit, err := core_http_utils.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return user_id, limit, offset, nil
}
