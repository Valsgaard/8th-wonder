package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

func (suite *EndpointsTestSuite) TestGameStateUpdate() {
	tests := []struct {
		Name               string
		UserID             string
		GamesPlayed        int
		Score              int
		ExpectedSuccess    bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Update",
			UserID:             suite.Users[0],
			GamesPlayed:        2,
			Score:              220,
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
		}, {
			Name:               "MissingUserID",
			UserID:             "",
			GamesPlayed:        2,
			Score:              220,
			ExpectedSuccess:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		}, {
			Name:               "ZeroGamesPlayed",
			UserID:             suite.Users[0],
			GamesPlayed:        0,
			Score:              1000,
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			b, err := json.Marshal(&GameStateUpdateInput{
				GamesPlayed: test.GamesPlayed,
				Score:       test.Score,
			})
			if err != nil {
				t.Fatal(err)
			}

			path := "/user/" + test.UserID + "/state"
			req, err := http.NewRequest("PUT", path, bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": test.UserID})

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewGameStateUpdate)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			if assert.Equal(t, test.ExpectedStatusCode, rr.Code) && test.ExpectedSuccess {
				// Check if state was properly updated
				state, err := port.GetDatastore(suite.ParentCtx).GetGameState(test.UserID)
				require.Nil(t, err)
				assert.Equal(t, test.GamesPlayed, state.GamesPlayed)
				assert.Equal(t, test.Score, state.Score)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}
