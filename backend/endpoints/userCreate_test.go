package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/valsgaard/interview-case/backend/common"
	. "github.com/valsgaard/interview-case/backend/endpoints"
)

func (suite *EndpointsTestSuite) TestUserCreate() {
	tests := []struct {
		Name               string
		InputName          string
		ExpectedSuccess    bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Post",
			InputName:          "Name1",
			ExpectedSuccess:    true,
			ExpectedStatusCode: http.StatusOK,
		}, {
			Name:               "MissingName",
			InputName:          "",
			ExpectedSuccess:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			b, err := json.Marshal(&UserCreateInput{
				Name: test.InputName,
			})
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			// Prepare recorder
			rr := httptest.NewRecorder()
			handler := common.NewHandlerFunc(suite.ParentCtx, NewUserCreate)

			// Call endpoint
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, test.ExpectedStatusCode, rr.Code)

			// Check the response body is what we expect.
			if test.ExpectedSuccess {
				v := new(UserCreateOutput)
				if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
					t.Fatal(err)
				}

				_, err := uuid.FromString(v.UserID)
				assert.Nil(t, err)
				assert.Equal(t, test.InputName, v.Name)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}
