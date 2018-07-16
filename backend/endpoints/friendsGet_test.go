package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
)

func (suite *EndpointsTestSuite) TestFriendsGet() {
	tests := []struct {
		Name               string
		UserID             string
		ExpectedSuccess    bool
		ExpectedStatusCode int
		ExpectedUsers      []string
		ExpectedNames      []string
		ExpectedHighscores []int
	}{
		{
			Name:               "Get",
			UserID:             suite.Users[0],
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
			ExpectedUsers:      []string{suite.Users[1], suite.Users[2], suite.Users[3]},
			ExpectedNames:      []string{"bot1", "bot2", "bot3"},
			ExpectedHighscores: []int{0, 0, 0},
		}, {
			Name:               "MissingUserID",
			UserID:             "",
			ExpectedSuccess:    false,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedUsers:      nil,
			ExpectedNames:      nil,
			ExpectedHighscores: nil,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			path := "/user/" + test.UserID + "/friends"
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": test.UserID})

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewFriendsGet)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, test.ExpectedStatusCode, rr.Code)

			// Check the response body is what we expect.
			if test.ExpectedSuccess {
				v := new(FriendsGetOutput)
				if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
					t.Fatal(err)
				}

				if assert.Equal(t, len(test.ExpectedUsers), len(v.Friends)) {
					for i := range v.Friends {
						assert.Equal(t, test.ExpectedUsers[i], v.Friends[i].ID)
						assert.Equal(t, test.ExpectedNames[i], v.Friends[i].Name)
						assert.Equal(t, test.ExpectedHighscores[i], v.Friends[i].Highscore)
					}
				}
			}

		}

		suite.T().Run(test.Name, fn)
	}
}
