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
	"github.com/google/uuid"
	"github.com/vitessio/arewefastyet/go/slack"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"html/template"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ErrorIncorrectConfiguration = "incorrect configuration"

	flagPort                                 = "web-port"
	flagTemplatePath                         = "web-template-path"
	flagStaticPath                           = "web-static-path"
	flagVitessPath                           = "web-vitess-path"
	flagMode                                 = "web-mode"
	flagMicroBenchConfigFile                 = "web-microbench-config"
	flagMacroBenchConfigFileOLTP             = "web-macrobench-oltp-config"
	flagMacroBenchConfigFileTPCC             = "web-macrobench-tpcc-config"
	flagCronSchedule                         = "web-cron-schedule"
	flagCronSchedulePullRequests             = "web-cron-schedule-pull-requests"
	flagCronScheduleTags                     = "web-cron-schedule-tags"
	flagPullRequestLabelTrigger              = "web-pr-label-trigger"
	flagPullRequestLabelTriggerWithPlannerV3 = "web-pr-label-trigger-planner-v3"
	flagCronNbRetry                          = "web-cron-nb-retry"
)

type Server struct {
	port         string
	templatePath string
	staticPath   string
	router       *gin.Engine

	vitessPathMu    sync.Mutex
	localVitessPath string

	dbCfg    *psdb.Config
	dbClient *psdb.Client

	// Configuration used to send message to Slack.
	slackConfig slack.Config

	cronSchedule             string
	cronSchedulePullRequests string
	cronScheduleTags         string
	cronNbRetry              int
	microbenchConfigPath     string
	macrobenchConfigPathOLTP string
	macrobenchConfigPathTPCC string

	prLabelTrigger   string
	prLabelTriggerV3 string

	// Mode used to run the server.
	Mode
}

func (s *Server) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.port, flagPort, "8080", "Port used for the HTTP server")
	cmd.Flags().StringVar(&s.templatePath, flagTemplatePath, "", "Path to the template directory")
	cmd.Flags().StringVar(&s.staticPath, flagStaticPath, "", "Path to the static directory")
	cmd.Flags().StringVar(&s.localVitessPath, flagVitessPath, "/", "Absolute path where the vitess directory is located or where it should be cloned")
	cmd.Flags().Var(&s.Mode, flagMode, "Specify the mode on which the server will run")

	// execution configuration flags
	cmd.Flags().StringVar(&s.microbenchConfigPath, flagMicroBenchConfigFile, "", "Path to the configuration file used to execute microbenchmark.")
	cmd.Flags().StringVar(&s.macrobenchConfigPathOLTP, flagMacroBenchConfigFileOLTP, "", "Path to the configuration file used to execute OLTP macrobenchmark.")
	cmd.Flags().StringVar(&s.macrobenchConfigPathTPCC, flagMacroBenchConfigFileTPCC, "", "Path to the configuration file used to execute TPCC macrobenchmark.")
	cmd.Flags().StringVar(&s.cronSchedule, flagCronSchedule, "@midnight", "Execution CRON schedule defaults to every day at midnight. An empty string will result in no CRON.")
	cmd.Flags().StringVar(&s.cronSchedulePullRequests, flagCronSchedulePullRequests, "*/5 * * * *", "Execution CRON schedule for pull requests benchmarks. An empty string will result in no CRON. Defaults to an execution every 5 minutes.")
	cmd.Flags().StringVar(&s.cronScheduleTags, flagCronScheduleTags, "*/1 * * * *", "Execution CRON schedule for tags/releases benchmarks. An empty string will result in no CRON. Defaults to an execution every minute.")
	cmd.Flags().IntVar(&s.cronNbRetry, flagCronNbRetry, 1, "Number of retries allowed for each cron job.")
	cmd.Flags().StringVar(&s.prLabelTrigger, flagPullRequestLabelTrigger, "Benchmark me", "GitHub Pull Request label that will trigger the execution of new execution.")
	cmd.Flags().StringVar(&s.prLabelTriggerV3, flagPullRequestLabelTriggerWithPlannerV3, "Benchmark me (V3)", "GitHub Pull Request label that will trigger the execution of new execution using the V3 planner.")
	_ = cmd.MarkFlagRequired(flagMicroBenchConfigFile)
	_ = cmd.MarkFlagRequired(flagMacroBenchConfigFileOLTP)
	_ = cmd.MarkFlagRequired(flagMacroBenchConfigFileTPCC)

	_ = viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	_ = viper.BindPFlag(flagTemplatePath, cmd.Flags().Lookup(flagTemplatePath))
	_ = viper.BindPFlag(flagStaticPath, cmd.Flags().Lookup(flagStaticPath))
	_ = viper.BindPFlag(flagVitessPath, cmd.Flags().Lookup(flagVitessPath))
	_ = viper.BindPFlag(flagMode, cmd.Flags().Lookup(flagMode))
	_ = viper.BindPFlag(flagMicroBenchConfigFile, cmd.Flags().Lookup(flagMicroBenchConfigFile))
	_ = viper.BindPFlag(flagMacroBenchConfigFileOLTP, cmd.Flags().Lookup(flagMacroBenchConfigFileOLTP))
	_ = viper.BindPFlag(flagMacroBenchConfigFileTPCC, cmd.Flags().Lookup(flagMacroBenchConfigFileTPCC))
	_ = viper.BindPFlag(flagCronSchedule, cmd.Flags().Lookup(flagCronSchedule))
	_ = viper.BindPFlag(flagCronSchedulePullRequests, cmd.Flags().Lookup(flagCronSchedulePullRequests))
	_ = viper.BindPFlag(flagCronScheduleTags, cmd.Flags().Lookup(flagCronScheduleTags))
	_ = viper.BindPFlag(flagCronNbRetry, cmd.Flags().Lookup(flagCronNbRetry))
	_ = viper.BindPFlag(flagPullRequestLabelTrigger, cmd.Flags().Lookup(flagPullRequestLabelTrigger))
	_ = viper.BindPFlag(flagPullRequestLabelTriggerWithPlannerV3, cmd.Flags().Lookup(flagPullRequestLabelTriggerWithPlannerV3))

	s.slackConfig.AddToCommand(cmd)
	if s.dbCfg == nil {
		s.dbCfg = &psdb.Config{}
	}
	s.dbCfg.AddToCommand(cmd)
}

func (s *Server) isReady() bool {
	return s.port != "" && s.templatePath != "" && s.staticPath != "" &&
		s.microbenchConfigPath != "" && s.macrobenchConfigPathOLTP != "" && s.macrobenchConfigPathTPCC != "" && s.localVitessPath != ""
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

	if err := s.setupLocalVitess(); err != nil {
		return err
	}

	if err := s.createStorages(); err != nil {
		return err
	}

	err := s.createCrons()
	if err != nil {
		return err
	}

	s.prepareGin()
	s.router = gin.Default()
	s.router.SetFuncMap(template.FuncMap{
		"formatFloat": func(f float64) string { return humanize.FormatFloat("#,###.##", f) },
		"formatBytes": func(f float64) string { return humanize.Bytes(uint64(f)) },
		"toString": func(i interface{}) string {
			switch i := i.(type) {
			case string:
				return i
			case []byte:
				return string(i)
			}
			return ""
		},
		"first8Letters": func(s string) string {
			if len(s) < 8 {
				return s
			}
			return s[:8]
		},
		"uuidToString": func(u uuid.UUID) string { return u.String() },
		"timeToDateString": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format(time.RFC822)
		},
	})

	s.router.Static("/static", s.staticPath)

	s.router.LoadHTMLGlob(s.templatePath + "/*")

	// Information page
	s.router.GET("/cron", s.cronHandler)

	s.router.GET("/analytics", s.analyticsHandlerOLTP)
	s.router.GET("/analytics_tpcc", s.analyticsHandlerTPCC)
	s.router.GET("/analytics_mixed", s.analyticsHandlerMixed)

	// Home page
	s.router.GET("/", s.homeHandler)

	// Compare page
	s.router.GET("/compare", s.compareHandler)

	// Search page
	s.router.GET("/search", s.searchHandler)

	// Request benchmark page
	s.router.GET("/request_benchmark", s.requestBenchmarkHandler)

	// Microbenchmark comparison page
	s.router.GET("/microbench", s.microbenchmarkResultsHandler)

	// Single Microbenchmark page
	s.router.GET("/microbench/:name", s.microbenchmarkSingleResultsHandler)

	// Macrobenchmark comparison page
	s.router.GET("/macrobench", s.macrobenchmarkResultsHandler)

	// Macrobenchmark queries details
	s.router.GET("/macrobench/queries/:git_ref", s.macrobenchmarkQueriesDetails)
	s.router.GET("/macrobench/queries/compare", s.macrobenchmarkCompareQueriesDetails)

	// V3 VS Gen4 comparison page
	s.router.GET("/v3_VS_Gen4", s.v3VsGen4Handler)

	// status page
	s.router.GET("/status", s.statusHandler)

	return s.router.Run(":" + s.port)
}

func (s *Server) prepareGin() {
	switch s.Mode {
	case ProductionMode:
		gin.SetMode(gin.ReleaseMode)
	case DevelopmentMode:
		gin.SetMode(gin.DebugMode)
	}
}

func Run(port, templatePath, staticPath, localVitessPath string) error {
	s := Server{
		port:                     port,
		templatePath:             templatePath,
		staticPath:               staticPath,
		localVitessPath:          localVitessPath,
		microbenchConfigPath:     "/",
		macrobenchConfigPathOLTP: "/",
		macrobenchConfigPathTPCC: "/",
	}
	return s.Run()
}
