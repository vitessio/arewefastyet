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
)

type Server struct {
	port         string
	templatePath string
	staticPath   string
	apiKey       string
	router       *gin.Engine
	dbCfg        *mysql.ConfigDB
	dbClient     *mysql.Client
}

func (s *Server) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.port, flagPort, "8080", "Port used for the HTTP server")
	cmd.Flags().StringVar(&s.templatePath, flagTemplatePath, "", "Path to the template directory")
	cmd.Flags().StringVar(&s.staticPath, flagStaticPath, "", "Path to the static directory")
	cmd.Flags().StringVar(&s.apiKey, flagAPIKey, "", "API key used to authenticate requests")

	viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	viper.BindPFlag(flagTemplatePath, cmd.Flags().Lookup(flagTemplatePath))
	viper.BindPFlag(flagStaticPath, cmd.Flags().Lookup(flagStaticPath))
	viper.BindPFlag(flagAPIKey, cmd.Flags().Lookup(flagAPIKey))

	if s.dbCfg == nil {
		s.dbCfg = &mysql.ConfigDB{}
	}
	s.dbCfg.AddToCommand(cmd)
}

func (s Server) isReady() bool {
	return s.port != "" && s.templatePath != "" && s.staticPath != "" && s.apiKey != ""
}

func (s *Server) Run() error {
	if s.isReady() == false {
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
	s.router.GET("/", s.homeHanlder)

	// Search and compare page
	s.router.GET("/search_compare", s.searchCompareHandler)

	// Request benchmark page
	s.router.GET("/request_benchmark", s.requestBenchmarkHandler)

	// Request benchmark page
	s.router.GET("/microbench", s.microbenchmarkResultsHandler)

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
