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
		RedirectURL: "http://localhost:8081/admin/auth/callback",
	}
	oauthStateString = random.String(10) // A random string to protect against CSRF attacks
	client           *goGithub.Client
	orgName          = "vitessio"
	tokens           = make(map[string]oauth2.Token)

	mu sync.Mutex
)

const (
	maintainerTeamGitHub   = "maintainers"
	arewefastyetTeamGitHub = "arewefastyet"
)

type ExecutionRequest struct {
	Auth               string   `json:"auth"`
	Source             string   `json:"source"`
	SHA                string   `json:"sha"`
	Workloads          []string `json:"workloads"`
	NumberOfExecutions string   `json:"number_of_executions"`
}

func (a *Admin) login(c *gin.Context) {
	a.render(c, gin.H{}, "login.html")
}

func (a *Admin) dashboard(c *gin.Context) {
	a.render(c, gin.H{}, "base.html")
}

func CreateGhClient(token *oauth2.Token) *goGithub.Client {
	return goGithub.NewClient(oauthConf.Client(context.Background(), token))
}

func (a *Admin) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("ghtoken")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		mu.Lock()
		defer mu.Unlock()

		token, ok := tokens[cookie]

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
		slog.Error("Error checking org membership: ", err)
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

	if isMaintainer {
		mu.Lock()
		defer mu.Unlock()

		randomKey := random.String(32)
		tokens[randomKey] = *token
		c.SetCookie("ghtoken", randomKey, int(time.Hour.Seconds()), "/", "localhost", true, true)

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
	source := c.PostForm("source")
	sha := c.PostForm("sha")
	workloads := c.PostFormArray("workloads")
	numberOfExecutions := c.PostForm("numberOfExecutions")

	if source == "" || sha == "" || len(workloads) == 0 || numberOfExecutions == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: Source and/or SHA"})
		return
	}

	token, err := c.Cookie("ghtoken")
	if err != nil {
		slog.Error("Failed to get token from cookie: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	encryptedToken := server.Encrypt(token, a.auth)

	slog.Info("Encrypted token: ", encryptedToken)

	requestPayload := ExecutionRequest{
		Auth:               encryptedToken,
		Source:             source,
		SHA:                sha,
		Workloads:          workloads,
		NumberOfExecutions: numberOfExecutions,
	}

	jsonData, err := json.Marshal(requestPayload)

	slog.Infof("Request payload: %s", jsonData)

	if err != nil {
		slog.Error("Failed to marshal request payload: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request payload"})
		return
	}

	serverAPIURL := "http://localhost:8080/api/executions/add"

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

	if resp.StatusCode != http.StatusOK {
		slog.Error("Server API returned an error: ", resp.Status)
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to process request on server API"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Execution(s) added successfully"})
}
