package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type FriendsGetInput struct {
	UserID string
}

type FriendsGetOutput struct {
	Friends []*Friend `json:"friends"`
}

// Friend is a part of FriendsGetOutput and contains details about a 'friend'
type Friend struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Highscore int    `json:"highscore"`
}

// NewFriendsGet is a HandlerFunc processing the request to retrieve a users friend list.
func NewFriendsGet(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()
	input := FriendsGetInput{
		UserID: mux.Vars(r)["id"],
	}

	if _, stderr := uuid.FromString(input.UserID); stderr != nil {
		return common.NewError(ErrBadRequest, "Invalid JSON format").
			SetInternal(stderr)
	}

	// Process data storage
	friends, err := port.GetDatastore(ctx).GetFriends(input.UserID)
	if err != nil {
		return err.SetStatusCode(http.StatusBadRequest)
	}

	// Prepare output
	output := new(FriendsGetOutput)
	output.Friends = make([]*Friend, 0, len(friends))
	for _, f := range friends {
		output.Friends = append(output.Friends, &Friend{
			ID:        f.UserID,
			Name:      f.Name,
			Highscore: f.HighScore,
		})
	}

	// Response
	return common.SuccessResponseJSON(rw, output)
}
