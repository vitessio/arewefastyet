/*
 *
 * Copyright 2024 The Vitess Authora.
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

package admin

import "go.uber.org/zap"

var slog *zap.SugaredLogger

func initLogger(mode Mode) (err error) {
	var logger *zap.Logger
	switch mode {
	case ProductionMode:
		logger, err = zap.NewProduction()
	case DevelopmentMode:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	slog = logger.Sugar()
	return nil
}

func (a *Admin) initLogger() (err error) {
	return initLogger(a.Mode)
}

func cleanLogger() {
	_ = slog.Sync()
}

// SetSLogger sets the *zap.SugaredLogger of this package.
func SetSLogger(newSlog *zap.SugaredLogger) {
	slog = newSlog
}
