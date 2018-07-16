package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

/**************************************************************************
***************************************************************************
**                                                                       **
**    Transport - HTTP                                                   **
**                                                                       **
***************************************************************************
**************************************************************************/

type handlerFunc func(rw http.ResponseWriter, r *http.Request) Error

// NewHandlerFunc adds our own context to the HTTP request and handles error response
func NewHandlerFunc(ctx context.Context, fn handlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Get mux Vars and add them to our own Context
		vars := mux.Vars(r)
		r = mux.SetURLVars(r.WithContext(ctx), vars)

		if err := fn(rw, r); err != nil {
			e := ErrorResponseJSON(rw, err.StatusCode(), &ErrorResponse{
				ErrorCode:    err.Code(),
				ErrorMessage: err.Message(),
			})

			if e != nil {
				Log(ctx).Warn(e)
				return
			}

			Log(ctx).Error(err)
		}
	}
}

// ListenAndServe start the http listener, and returns a graceful shudown function
func ListenAndServe(ctx context.Context, server *http.Server, t time.Duration) func() {
	go func() {
		if err := server.ListenAndServe(); err != nil {
			Log(ctx).Error(err)
		}
	}()

	return func() {
		c, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		if err := server.Shutdown(c); err != nil {
			Log(ctx).Error(err)
		}
	}
}

/**************************************************************************
***************************************************************************
**                                                                       **
**    Transport - HTTP - Request                                         **
**                                                                       **
***************************************************************************
**************************************************************************/

// ReadJSONRequest reads and unmarshals the requests json body into the
// given struct
func ReadJSONRequest(r *http.Request, v interface{}) Error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewError(err, "")
	}

	if err := json.Unmarshal(b, v); err != nil {
		return NewError(err, "")
	}

	return nil
}

/**************************************************************************
***************************************************************************
**                                                                       **
**    Transport - HTTP - Response                                        **
**                                                                       **
***************************************************************************
**************************************************************************/

// ErrorResponse is the default error response json format
type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// ErrorResponseJSON sends a response, intended for error handling.
// If no status code is given, it defaults to 500.
func ErrorResponseJSON(rw http.ResponseWriter, code int, v interface{}) Error {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("1")
		return NewError(err, "")
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	rw.WriteHeader(code)
	rw.Header().Set("content-type", "application/json")
	_, err = rw.Write(b)
	if err != nil {
		return NewError(err, "")
	}

	return nil
}

// SuccessResponseJSON marshals and writes the given
func SuccessResponseJSON(rw http.ResponseWriter, v interface{}) Error {
	b, err := json.Marshal(v)
	if err != nil {
		return NewError(err, "")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("content-type", "application/json")
	_, err = rw.Write(b)
	if err != nil {
		return NewError(err, "")
	}

	return nil
}

// SuccessResponseEmpty sends a response with an empty body
func SuccessResponseEmpty(rw http.ResponseWriter) Error {
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write(nil)
	if err != nil {
		return NewError(err, "")
	}

	return nil
}

/**************************************************************************
***************************************************************************
**                                                                       **
**    Flow control                                                       **
**                                                                       **
***************************************************************************
**************************************************************************/

// SystemTermination returns a channel waiting for termination from the OS
func SystemTermination() chan os.Signal {
	sigchan := make(chan os.Signal, 0)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	return sigchan
}
