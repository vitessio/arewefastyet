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
	"github.com/vitessio/arewefastyet/go/tools/git"
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
		Ref        string `json:"ref"`
		PathConfig string `json:"path_config"`
	}
	var wbhPayload webhookPayload

	err := json.NewDecoder(c.Request.Body).Decode(&wbhPayload)
	if err != nil {
		slog.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, newResponseError(err))
		return
	} else if wbhPayload.Ref == "" {
		errorStr := fmt.Sprintf(errorNeedsValue, "ref")
		slog.Warn(errorStr)
		c.JSON(http.StatusBadRequest, newResponseErrorFromString(errorStr))
		return
	} else if wbhPayload.Ref != "refs/heads/master" {
		slog.Warn("not a runnable reference")
		c.JSON(http.StatusBadRequest, newResponseErrorFromString("not a runnable reference"))
		return
	}

	pathConfig := wbhPayload.PathConfig
	if pathConfig == "" {
		pathConfig = s.defaultExecConfigFile
	}

	// Will load any configuration (microbench, OLTP, TPCC, OLTP+TPCC, etc).
	e, err := exec.NewExecWithConfig(pathConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newResponseError(err))
		return
	}

	ref, err := git.GetCommitHashFromClonedRef(wbhPayload.Ref, "https://github.com/vitessio/vitess")
	if err != nil {
		slog.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, newResponseError(err))
		return
	}

	e.Source = "webhook"
	e.GitRef = ref

	// Start a goroutine with the running execution
	go func() {
		err = e.Prepare()
		if err != nil {
			slog.Error("Prepare step", err.Error())
			return
		}

		err = e.SetOutputToDefaultPath()
		if err != nil {
			slog.Error("Prepare step", err.Error())
			return
		}

		err = e.Execute()
		if err != nil {
			slog.Error("Execution step", err.Error())
			return
		}

		err = e.CleanUp()
		if err != nil {
			slog.Error("Clean step", err.Error())
			return
		}
	}()

	type webhookResponse struct {
		Started bool   `json:"started"`
		UUID    string `json:"uuid"`
	}
	c.JSON(http.StatusCreated, webhookResponse{Started: true, UUID: e.UUID.String()})
}
