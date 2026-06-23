package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/raizonw/todo-app/internal/core/logger"
	core_http_response "github.com/raizonw/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/raizonw/todo-app/internal/core/transport/http/utils"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получить задачу
// @Description Получить подробную информацию о задаче по её ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse "Информация о задаче"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `id` path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
