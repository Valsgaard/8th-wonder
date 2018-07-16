package endpoints

import (
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type UserCreateInput struct {
	Name string `json:"name"`
}

type UserCreateOutput struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
}

// NewUserCreate is a HandlerFunc processing the request to create new users.
func NewUserCreate(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()

	// Parse input
	input := new(UserCreateInput)
	if err := common.ReadJSONRequest(r, input); err != nil {
		return common.ErrorResponseJSON(
			rw,
			http.StatusBadRequest,
			&common.ErrorResponse{
				ErrorCode:    "E000",
				ErrorMessage: err.Error(),
			},
		)
	}

	// Validate name
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return common.ErrorResponseJSON(
			rw,
			http.StatusBadRequest,
			&common.ErrorResponse{
				ErrorCode:    "E000",
				ErrorMessage: "Invalid name",
			},
		)
	}

	// Prepare UserID
	// Note: We're just going to use a V1 UUID for now and cross fingers
	// that there won't be any collisions
	userID := uuid.NewV1().String()

	// Process data storage
	user, err := port.GetDatastore(ctx).NewUser(userID, input.Name)
	if err != nil {
		return common.ErrorResponseJSON(
			rw,
			http.StatusBadRequest,
			&common.ErrorResponse{
				ErrorCode:    "E000",
				ErrorMessage: err.Error(),
			},
		)
	}

	// Response
	return common.SuccessResponseJSON(rw, &UserCreateOutput{
		UserID: user.UserID,
		Name:   user.Name,
	})
}
