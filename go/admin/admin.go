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

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"github.com/vitessio/arewefastyet/go/tools/server"
)

const (
	flagPort           = "admin-port"
	flagMode           = "admin-mode"
	flagAdminAppId     = "admin-gh-app-id"
	flagAdminAppSecret = "admin-gh-app-secret"
	flagGhAuth         = "admin-auth"
)

type Admin struct {
	port   string
	router *gin.Engine

	ghAppId     string
	ghAppSecret string
	auth        string

	dbCfg    *psdb.Config
	dbClient *psdb.Client

	server.Mode
}

func (a *Admin) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&a.port, flagPort, "8081", "Port used for the HTTP server")
	cmd.Flags().Var(&a.Mode, flagMode, "Specify the mode on which the server will run")
	cmd.Flags().StringVar(&a.ghAppId, flagAdminAppId, "", "The ID of the GitHub App")
	cmd.Flags().StringVar(&a.ghAppSecret, flagAdminAppSecret, "", "The secret of the GitHub App")
	cmd.Flags().StringVar(&a.auth, flagGhAuth, "", "The salt string to salt the GitHub Token")

	_ = viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	_ = viper.BindPFlag(flagMode, cmd.Flags().Lookup(flagMode))
	_ = viper.BindPFlag(flagAdminAppId, cmd.Flags().Lookup(flagAdminAppId))
	_ = viper.BindPFlag(flagAdminAppSecret, cmd.Flags().Lookup(flagAdminAppSecret))
	_ = viper.BindPFlag(flagGhAuth, cmd.Flags().Lookup(flagGhAuth))

	if a.dbCfg == nil {
		a.dbCfg = &psdb.Config{}
	}
	a.dbCfg.AddToCommand(cmd)
}

func (a *Admin) isReady() bool {
	return a.port != "" && a.ghAppId != "" && a.ghAppSecret != ""
}

func (a *Admin) Init() error {
	if !a.isReady() {
		return errors.New(server.ErrorIncorrectConfiguration)
	}

	if a.Mode != "" && !a.Mode.Correct() {
		return errors.New(server.ErrorIncorrectMode)
	} else if a.Mode == "" {
		a.Mode.UseDefault()
	}

	if slog == nil {
		err := a.initLogger()
		if err != nil {
			return err
		}
		defer cleanLogger()
	}

	if err := a.createStorages(); err != nil {
		return err
	}

	return nil
}

func (a *Admin) Run() error {
	if !a.isReady() {
		return errors.New(server.ErrorIncorrectConfiguration)
	}

	a.Mode.SetGin()
	a.router = gin.Default()

	a.router.Static("/admin/assets", "/go/admin/assets")

	a.router.LoadHTMLGlob("/go/admin/templates/*")

	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// OAuth
	a.router.GET("/admin", a.login)
	a.router.GET("/admin/login", a.handleGitHubLogin)
	a.router.GET("/admin/auth/callback", a.handleGitHubCallback)

	// API
	a.router.POST("/admin/executions/add", a.authMiddleware(), a.handleExecutionsAdd)
	a.router.POST("/admin/executions/clear", a.authMiddleware(), a.handleClearQueue)

	// Pages
	a.router.GET("/admin/dashboard", a.authMiddleware(), a.homePage)
	a.router.GET("/admin/dashboard/newexec", a.authMiddleware(), a.newExecutionsPage)
	a.router.GET("/admin/dashboard/clearqueue", a.authMiddleware(), a.clearQueuePage)

	return a.router.Run(":" + a.port)
}

func (a *Admin) render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
