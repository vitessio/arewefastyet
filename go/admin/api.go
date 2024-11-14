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
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	goGithub "github.com/google/go-github/github"
	"github.com/labstack/gommon/random"
	"github.com/vitessio/arewefastyet/go/tools/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		Scopes:      []string{"read:org"}, // Request access to read organization membership
		Endpoint:    github.Endpoint,
		RedirectURL: "http://localhost/admin/auth/callback",
	}
	oauthStateString = random.String(10) // A random string to protect against CSRF attacks
	orgName          = "vitessio"
	tokens           = make(map[string]oauth2.Token)

	mu sync.Mutex
)

const (
	maintainerTeamGitHub   = "maintainers"
	arewefastyetTeamGitHub = "arewefastyet"
)

type (
	executionRequest struct {
		Auth               string   `json:"auth"`
		Source             string   `json:"source"`
		SHA                string   `json:"sha"`
		Workloads          []string `json:"workloads"`
		NumberOfExecutions string   `json:"number_of_executions"`
		EnableProfile      bool     `json:"enable_profile"`
		BinaryToProfile    string   `json:"binary_to_profile"`
		ProfileMode        string   `json:"profile_mode"`
	}

	clearQueueRequest struct {
		Auth                  string `json:"auth"`
		RemoveAdminExecutions bool   `json:"remove_admin_executions"`
	}
)

func (a *Admin) login(c *gin.Context) {
	a.render(c, gin.H{}, "login.html")
}

func (a *Admin) homePage(c *gin.Context) {
	a.render(c, gin.H{}, "base.html")
}

func (a *Admin) newExecutionsPage(c *gin.Context) {
	a.render(c, gin.H{"Page": "newexec"}, "base.html")
}

func (a *Admin) clearQueuePage(c *gin.Context) {
	a.render(c, gin.H{"Page": "clearqueue"}, "base.html")
}

func CreateGhClient(token *oauth2.Token) *goGithub.Client {
	return goGithub.NewClient(oauthConf.Client(context.Background(), token))
}

func (a *Admin) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("tk")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		mu.Lock()
		token, ok := tokens[cookie]
		mu.Unlock()
		if !ok {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		client := CreateGhClient(&token)

		isMaintainer, err := a.GetUser(client)

		if err != nil {
			slog.Error("Error getting user: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !isMaintainer {
			c.String(http.StatusForbidden, "You must be a maintainer in the %s organization to access this page.", orgName)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *Admin) GetUser(client *goGithub.Client) (bool, error) {
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return false, err
	}

	isMaintainer, err := a.checkUserOrgMembership(client, user.GetLogin(), orgName)
	if err != nil {
		return false, err
	}

	return isMaintainer, nil
}

func (a *Admin) handleGitHubLogin(c *gin.Context) {
	if a.Mode == server.ProductionMode {
		oauthConf.RedirectURL = "https://benchmark.vitess.io/admin/auth/callback"
	}
	oauthConf.ClientID = a.ghAppId
	oauthConf.ClientSecret = a.ghAppSecret
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *Admin) handleGitHubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Query("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		slog.Error("Code exchange failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := CreateGhClient(token)

	isMaintainer, err := a.GetUser(client)
	if err != nil {
		slog.Error("Failed to get user information: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if isMaintainer {
		mu.Lock()
		defer mu.Unlock()

		randomKey := random.String(32)
		tokens[randomKey] = *token

		domain := "localhost"

		if a.Mode == server.ProductionMode {
			domain = "benchmark.vitess.io"
		}

		c.SetCookie("tk", randomKey, int(time.Hour.Seconds()), "/", domain, true, true)

		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	} else {
		c.String(http.StatusForbidden, "You must be a maintainer in the %s organization to access this page.", orgName)
	}
}

func (a *Admin) checkUserOrgMembership(client *goGithub.Client, username, orgName string) (bool, error) {
	teams, _, err := client.Teams.ListTeams(context.Background(), orgName, nil)
	if err != nil {
		return false, err
	}

	var isMember bool
	for _, team := range teams {
		if team.GetName() == maintainerTeamGitHub || team.GetName() == arewefastyetTeamGitHub {
			membership, _, err := client.Teams.GetTeamMembership(context.Background(), team.GetID(), username)
			if err != nil {
				if strings.Contains(err.Error(), "404 Not Found") {
					continue
				}
				return false, err
			}
			if membership.GetState() == "active" {
				isMember = true
				break
			}
		}
	}
	return isMember, nil
}

func (a *Admin) handleExecutionsAdd(c *gin.Context) {
	requestPayload := executionRequest{
		Source:             c.PostForm("source"),
		SHA:                c.PostForm("sha"),
		Workloads:          c.PostFormArray("workloads"),
		NumberOfExecutions: c.PostForm("numberOfExecutions"),
		EnableProfile:      c.PostForm("enableProfiling") != "",
		BinaryToProfile:    c.PostForm("binaryToProfile"),
		ProfileMode:        c.PostForm("profileMode"),
	}

	if requestPayload.Source == "" || requestPayload.SHA == "" || len(requestPayload.Workloads) == 0 || requestPayload.NumberOfExecutions == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: Source, SHA, workflows, numberOfExecutions"})
		return
	}

	if requestPayload.EnableProfile && (requestPayload.BinaryToProfile == "" || requestPayload.ProfileMode == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "When enabling profiling, please provide a binary to profile and a mode"})
		return
	}

	tokenKey, err := c.Cookie("tk")
	if err != nil {
		slog.Error("Failed to get token from cookie: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	token, exists := tokens[tokenKey]

	if !exists {
		slog.Error("Failed to get token from map")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requestPayload.Auth, err = server.Encrypt(token.AccessToken, a.auth)

	if err != nil {
		slog.Error("Failed to encrypt token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt token"})
		return
	}

	jsonData, err := json.Marshal(requestPayload)

	if err != nil {
		slog.Error("Failed to marshal request payload: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request payload"})
		return
	}

	serverAPIURL := getAPIURL(a.Mode, "/executions/add")

	req, err := http.NewRequest("POST", serverAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to create new HTTP request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to server API"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request to server API: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to server API"})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		slog.Error("Server API returned an error: ", resp.Status)
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to process request on server API"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Execution(s) added successfully"})
}

func (a *Admin) handleClearQueue(c *gin.Context) {
	tokenKey, err := c.Cookie("tk")
	if err != nil {
		slog.Error("Failed to get token from cookie: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	mu.Lock()
	token, exists := tokens[tokenKey]
	mu.Unlock()

	if !exists {
		slog.Error("Failed to get token from map")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	encryptedToken, err := server.Encrypt(token.AccessToken, a.auth)
	if err != nil {
		slog.Error("Failed to encrypt token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt token"})
		return
	}

	requestPayload := clearQueueRequest{
		Auth:                  encryptedToken,
		RemoveAdminExecutions: c.PostForm("remove_admin") == "true",
	}

	log.Println(requestPayload)

	jsonData, err := json.Marshal(requestPayload)

	if err != nil {
		slog.Error("Failed to marshal request payload: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request payload"})
		return
	}

	serverAPIURL := getAPIURL(a.Mode, "/executions/clear")

	req, err := http.NewRequest("POST", serverAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to create new HTTP request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to server API"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request to server API: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to server API"})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		slog.Error("Server API returned an error: ", resp.Status)
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to process request on server API"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Execution(s) added successfully"})
}

func getAPIURL(mode server.Mode, endpoint string) string {
	serverAPIURL := "http://traefik/api" + endpoint
	if mode == server.ProductionMode {
		serverAPIURL = "https://benchmark.vitess.io/api" + endpoint
	}
	return serverAPIURL
}
