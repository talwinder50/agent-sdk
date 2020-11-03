/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package rest provides rest operation.
package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/common/log"

	"github.com/trustbloc/agent-sdk/pkg/controller/command"
)

var logger = log.New("aries-framework/rest")

// Handler http handler for each controller API endpoint.
type Handler interface {
	Path() string
	Method() string
	Handle() http.HandlerFunc
}

// Execute executes given command with args provided and writes error to
// response writer.
func Execute(exec command.Exec, rw http.ResponseWriter, req io.Reader) {
	rw.Header().Set("Content-Type", "application/json")

	err := exec(rw, req)
	if err != nil {
		SendError(rw, err)
	}
}

type genericErrorBody struct {
	Code    command.Code `json:"code"`
	Message string       `json:"message"`
}

// SendError sends command error as http response in generic error format.
func SendError(rw http.ResponseWriter, err command.Error) {
	var status int

	status = http.StatusInternalServerError
	if err.Type() == command.ValidationError {
		status = http.StatusBadRequest
	}

	SendHTTPStatusError(rw, status, err.Code(), err)
}

// SendHTTPStatusError sends given http status code to response with error body.
func SendHTTPStatusError(rw http.ResponseWriter, httpStatus int, code command.Code, err error) {
	rw.WriteHeader(httpStatus)

	e := json.NewEncoder(rw).Encode(genericErrorBody{
		Code:    code,
		Message: err.Error(),
	})
	if e != nil {
		logger.Errorf("Unable to send error response, %s", e)
	}
}
