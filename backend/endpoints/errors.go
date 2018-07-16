package endpoints

import (
	"net/http"

	"github.com/valsgaard/interview-case/backend/common"
)

// ErrBadRequest indicates that the request has an invalid input
var ErrBadRequest = common.PrepareError("EE001", "Bad request, input entries are invalid, malformed or missing").
	SetStatusCode(http.StatusBadRequest)
