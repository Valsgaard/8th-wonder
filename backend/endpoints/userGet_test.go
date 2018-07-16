package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
)

func (suite *EndpointsTestSuite) TestUserGet() {
	tests := []struct {
		Name               string
		ExpectedSuccess    bool
		ExpectedStatusCode int
		ExpectedUsers      []string
		ExpectedNames      []string
	}{
		{
			Name:               "Get",
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
			ExpectedUsers:      suite.Users,
			ExpectedNames:      []string{"bot0", "bot1", "bot2", "bot3"},
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest("GET", "/user", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewUserGet)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, test.ExpectedStatusCode, rr.Code)

			// Check the response body is what we expect.
			if test.ExpectedSuccess {
				v := new(UserGetOutput)
				if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
					t.Fatal(err)
				}

				if assert.Equal(t, len(test.ExpectedUsers), len(v.Users)) {
					for i := range v.Users {
						assert.Equal(t, test.ExpectedUsers[i], v.Users[i].UserID)
						assert.Equal(t, test.ExpectedNames[i], v.Users[i].Name)
					}
				}
			}
		}

		suite.T().Run(test.Name, fn)
	}
}
