package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type FriendsUpdateInput struct {
	UserID  string
	Friends []string `json:"friends"`
}

// NewFriendsUpdate is a HandlerFunc processing the request to update a users friend list.
func NewFriendsUpdate(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()

	// Parse input
	input := new(FriendsUpdateInput)
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

	input.UserID = mux.Vars(r)["id"]

	// Process data storage
	err := port.GetDatastore(ctx).UpdateFriends(
		input.UserID,
		input.Friends,
	)

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
	return common.SuccessResponseEmpty(rw)
}
