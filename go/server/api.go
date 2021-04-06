/*
 *
 * Copyright 2021 The Vitess Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/exec"
	"log"
	"net/http"
)

type responseError struct {
	Err string `json:"error"`
}

const (
	errorNeedsValue = "needs '%s' value"
)

func newResponseErrorFromString(err string) responseError {
	return responseError{Err: err}
}

func newResponseError(err error) responseError {
	return responseError{Err: err.Error()}
}

func (s *Server) webhookHandler(c *gin.Context) {
	type webhookPayload struct {
		Ref string `json:"ref"`
	}
	var wbhPayload webhookPayload

	err := json.NewDecoder(c.Request.Body).Decode(&wbhPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newResponseError(err))
		return
	} else if wbhPayload.Ref == "" {
		c.JSON(http.StatusBadRequest, newResponseErrorFromString(fmt.Sprintf(errorNeedsValue, "ref")))
		return
	}

	e, err := exec.NewExecWithConfig("")
	if err != nil {
		c.JSON(http.StatusInternalServerError, newResponseError(err))
		return
	}

	log.Println(e)

	// TODO: concurrent call to macro bench here

	type webhookResponse struct {
		Started bool `json:"started"`
	}
	c.JSON(http.StatusOK, webhookResponse{Started: true})
}
