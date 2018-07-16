package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
)

func (suite *EndpointsTestSuite) TestFriendsUpdate() {
	tests := []struct {
		Name               string
		UserID             string
		Friends            []string
		ExpectedSuccess    bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Update",
			UserID:             suite.Users[0],
			Friends:            []string{suite.Users[1], suite.Users[2], suite.Users[3]},
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
		}, {
			Name:               "MissingUserID",
			UserID:             "",
			Friends:            []string{suite.Users[0], suite.Users[2]},
			ExpectedSuccess:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		}, {
			Name:               "EmptyInput",
			UserID:             suite.Users[0],
			Friends:            []string{},
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			b, err := json.Marshal(&FriendsUpdateInput{
				Friends: test.Friends,
			})
			if err != nil {
				t.Fatal(err)
			}

			path := "/user/" + test.UserID + "/friends"
			req, err := http.NewRequest("PUT", path, bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": test.UserID})

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewFriendsUpdate)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, test.ExpectedStatusCode, rr.Code)
		}

		suite.T().Run(test.Name, fn)
	}
}
