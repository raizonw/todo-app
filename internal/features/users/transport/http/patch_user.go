package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/raizonw/todo-app/internal/core/domain"
	core_logger "github.com/raizonw/todo-app/internal/core/logger"
	core_http_request "github.com/raizonw/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/raizonw/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/raizonw/todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/raizonw/todo-app/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}
type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode HTTP request",
		)
		return
	}

	userPatch := UserPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields:\nFullName: '%s'\nPhoneNumber: '%s'",
			request.FullName,
			request.PhoneNumber,
		),
	)

	rw.WriteHeader(http.StatusOK)
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
