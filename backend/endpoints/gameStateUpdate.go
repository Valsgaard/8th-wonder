package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type GameStateUpdateInput struct {
	UserID      string
	GamesPlayed int `json:"gamesPlayed"`
	Score       int `json:"score"`
}

// NewGameStateUpdate is a HandlerFunc processing the request to update a users game state.
func NewGameStateUpdate(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()

	// Parse input
	input := new(GameStateUpdateInput)
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
	err := port.GetDatastore(ctx).UpdateGameState(
		input.UserID,
		input.GamesPlayed,
		input.Score,
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
