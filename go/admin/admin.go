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
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
)

var workloads = []string{"OLTP", "OLTP-READONLY", "OLTP-SET", "TPCC", "TPCC_FK", "TPCC_FK_UNMANAGED", "TPCC_UNSHARDED"}

type Admin struct {
	port   string
	router *gin.Engine

	ghAppId     string
	ghAppSecret string

	dbCfg    *psdb.Config
	dbClient *psdb.Client

	server.Mode
}

func (a *Admin) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&a.port, flagPort, "8081", "Port used for the HTTP server")
	cmd.Flags().Var(&a.Mode, flagMode, "Specify the mode on which the server will run")
	cmd.Flags().StringVar(&a.ghAppId, flagAdminAppId, "", "The ID of the GitHub App")
	cmd.Flags().StringVar(&a.ghAppSecret, flagAdminAppSecret, "", "The secret of the GitHub App")

	_ = viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	_ = viper.BindPFlag(flagMode, cmd.Flags().Lookup(flagMode))
	_ = viper.BindPFlag(flagAdminAppId, cmd.Flags().Lookup(flagAdminAppId))
	_ = viper.BindPFlag(flagAdminAppSecret, cmd.Flags().Lookup(flagAdminAppSecret))

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

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	a.Mode.SetGin()
	a.router = gin.Default()

	store := cookie.NewStore([]byte("secret"))
	a.router.Use(sessions.Sessions("mysession", store))

	a.router.Static("/assets", filepath.Join(basepath, "assets"))

	a.router.LoadHTMLGlob(filepath.Join(basepath, "templates/*"))

	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API
	a.router.GET("/admin", a.login)
	a.router.GET("/admin/login", a.handleGitHubLogin)
	a.router.GET("/admin/auth/callback", a.handleGitHubCallback)
	a.router.GET("/admin/dashboard", a.authMiddleware(), a.dashboard)

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
