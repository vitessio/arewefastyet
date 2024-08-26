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
	"github.com/vitessio/arewefastyet/go/tools/server"
	"go.uber.org/zap"
)

var slog *zap.SugaredLogger

func initLogger(mode server.Mode) (err error) {
	var logger *zap.Logger
	switch mode {
	case server.ProductionMode:
		logger, err = zap.NewProduction()
	case server.DevelopmentMode:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	slog = logger.Sugar()
	return nil
}

func (s *Server) initLogger() (err error) {
	return initLogger(s.Mode)
}

func cleanLogger() {
	_ = slog.Sync()
}

// SetSLogger sets the *zap.SugaredLogger of this package.
func SetSLogger(newSlog *zap.SugaredLogger) {
	slog = newSlog
}
