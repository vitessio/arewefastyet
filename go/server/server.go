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
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/mysql"
)

const (
	ErrorIncorrectConfiguration = "incorrect configuration"

	flagPort         = "web-port"
	flagTemplatePath = "web-template-path"
	flagStaticPath   = "web-static-path"
	flagAPIKey       = "web-api-key"
	flagMode         = "web-mode"
)

type Server struct {
	port         string
	templatePath string
	staticPath   string
	apiKey       string
	router       *gin.Engine
	dbCfg        *mysql.ConfigDB
	dbClient     *mysql.Client

	// Mode used to run the server.
	Mode
}

func (s *Server) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.port, flagPort, "8080", "Port used for the HTTP server")
	cmd.Flags().StringVar(&s.templatePath, flagTemplatePath, "", "Path to the template directory")
	cmd.Flags().StringVar(&s.staticPath, flagStaticPath, "", "Path to the static directory")
	cmd.Flags().StringVar(&s.apiKey, flagAPIKey, "", "API key used to authenticate requests")
	cmd.Flags().Var(&s.Mode, flagMode, "Specify the mode on which the server will run")

	_ = viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	_ = viper.BindPFlag(flagTemplatePath, cmd.Flags().Lookup(flagTemplatePath))
	_ = viper.BindPFlag(flagStaticPath, cmd.Flags().Lookup(flagStaticPath))
	_ = viper.BindPFlag(flagAPIKey, cmd.Flags().Lookup(flagAPIKey))
	_ = viper.BindPFlag(flagMode, cmd.Flags().Lookup(flagMode))

	if s.dbCfg == nil {
		s.dbCfg = &mysql.ConfigDB{}
	}
	s.dbCfg.AddToCommand(cmd)
}

func (s Server) isReady() bool {
	return s.port != "" && s.templatePath != "" && s.staticPath != "" && s.apiKey != ""
}

func (s *Server) Run() error {
	if s.Mode != "" && !s.Mode.correct() {
		return errors.New(ErrorIncorrectMode)
	} else if s.Mode == "" {
		s.Mode.useDefault()
	}

	if slog == nil {
		err := s.initLogger()
		if err != nil {
			return err
		}
		defer cleanLogger()
	}

	if !s.isReady() {
		return errors.New(ErrorIncorrectConfiguration)
	}

	if err := s.setupMySQL(); err != nil {
		return err
	}

	s.router = gin.Default()
	s.router.Static("/static", s.staticPath)

	s.router.LoadHTMLGlob(s.templatePath + "/*")

	// Information page
	s.router.GET("/information", s.informationHandler)

	// Home page
	s.router.GET("/", s.homeHandler)

	// Search and compare page
	s.router.GET("/search_compare", s.searchCompareHandler)

	// Request benchmark page
	s.router.GET("/request_benchmark", s.requestBenchmarkHandler)

	// Request benchmark page
	s.router.GET("/microbench", s.microbenchmarkResultsHandler)

	// MacroBench webhook
	s.router.POST("/webhook", s.webhookHandler)

	return s.router.Run(":" + s.port)
}

func Run(port, templatePath, staticPath, apiKey string) error {
	s := Server{
		port:         port,
		templatePath: templatePath,
		staticPath:   staticPath,
		apiKey:       apiKey,
	}
	return s.Run()
}
