package endpoints

import (
	"net/http"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type UserGetOutput struct {
	Users []*User `json:"users"`
}

// User is used as part of UserGetOutput and contains the details of a user
type User struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
}

// NewUserGet is a HandlerFunc processing the request to retrieve users.
func NewUserGet(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()

	// Process data storage
	users, err := port.GetDatastore(ctx).GetUsers()
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

	output := new(UserGetOutput)
	output.Users = make([]*User, 0, len(users))
	for _, u := range users {
		output.Users = append(output.Users, &User{UserID: u.UserID, Name: u.Name})
	}

	// Response
	return common.SuccessResponseJSON(rw, output)
}
