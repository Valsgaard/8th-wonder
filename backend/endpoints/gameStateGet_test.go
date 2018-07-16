package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
)

func (suite *EndpointsTestSuite) TestGameStateGet() {
	tests := []struct {
		Name                string
		UserID              string
		ExpectedSuccess     bool
		ExpectedStatusCode  int
		ExpectedGamesPlayed int
		ExpectedScore       int
	}{
		{
			Name:                "Get",
			UserID:              suite.Users[0],
			ExpectedSuccess:     true,
			ExpectedStatusCode:  http.StatusOK,
			ExpectedGamesPlayed: 10,
			ExpectedScore:       110,
		}, {
			Name:                "GetZero",
			UserID:              suite.Users[1],
			ExpectedSuccess:     true,
			ExpectedStatusCode:  http.StatusOK,
			ExpectedGamesPlayed: 0,
			ExpectedScore:       0,
		}, {
			Name:                "MissingUserID",
			UserID:              "",
			ExpectedSuccess:     false,
			ExpectedStatusCode:  http.StatusBadRequest,
			ExpectedGamesPlayed: 0,
			ExpectedScore:       0,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			path := "/user/" + test.UserID + "/state"
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": test.UserID})

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewGameStateGet)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			require.Equal(t, test.ExpectedStatusCode, rr.Code)

			// Check the response body is what we expect.
			if test.ExpectedSuccess {
				v := new(GameStateGetOutput)
				if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, test.ExpectedGamesPlayed, v.GamesPlayed)
				assert.Equal(t, test.ExpectedScore, v.Score)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}
