package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type GameStateGetInput struct {
	UserID string
}

type GameStateGetOutput struct {
	GamesPlayed int `json:"gamesPlayed"`
	Score       int `json:"score"`
}

// NewGameStateGet is a HandlerFunc processing the request to retrieve a users game state.
func NewGameStateGet(rw http.ResponseWriter, r *http.Request) common.Error {
	ctx := r.Context()
	input := GameStateGetInput{
		UserID: mux.Vars(r)["id"],
	}

	// Process data storage
	gameState, err := port.GetDatastore(ctx).GetGameState(input.UserID)
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
	return common.SuccessResponseJSON(rw, &GameStateGetOutput{
		GamesPlayed: gameState.GamesPlayed,
		Score:       gameState.Score,
	})
}
