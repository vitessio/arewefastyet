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
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	goGithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"read:org"}, // Request access to read organization membership
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:9090/auth/callback",
	}
	oauthStateString = "randomStateString" // A random string to protect against CSRF attacks
)

func (a *Admin) login(c *gin.Context) {
	a.render(c, gin.H{}, "login.html")
}

func (a *Admin) dashboard(c *gin.Context) {
	a.render(c, gin.H{}, "base.html")
}

func (a *Admin) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// User not authenticated, redirect to login
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// User is authenticated, proceed to the next handler
		c.Next()
	}
}

func (a *Admin) handleGitHubLogin(c *gin.Context) {
	oauthConf.ClientID = a.ghAppId
	oauthConf.ClientSecret = a.ghAppSecret
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (a *Admin) handleGitHubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		log.Println("Invalid OAuth state")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Query("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Code exchange failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := goGithub.NewClient(oauthConf.Client(context.Background(), token))

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		log.Println("Failed to get user: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Printf("Authenticated user: %s", user.GetLogin())

	orgName := "github-go-htmx-oauth-test"
	isMaintainer, err := a.checkUserOrgMembership(client, user.GetLogin(), orgName)
	if err != nil {
		log.Printf("Error checking org membership: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if isMaintainer {
		session := sessions.Default(c)
		session.Set("user", user.GetLogin())
		session.Save()

		c.Redirect(http.StatusSeeOther, "/dashboard")
	} else {
		log.Printf("User %s is not a maintainer in %s organization", user.GetLogin(), orgName)
		c.String(http.StatusForbidden, "You must be a maintainer in the %s organization to access this page.", orgName)
	}
}
func (a *Admin) checkUserOrgMembership(client *goGithub.Client, username, orgName string) (bool, error) {
	teams, _, err := client.Teams.ListTeams(context.Background(), orgName, nil)
	if err != nil {
		return false, err
	}
	for _, team := range teams {
		if team.GetName() == "maintainers" {
			membership, _, err := client.Teams.GetTeamMembership(context.Background(), team.GetID(), username)
			if err != nil {
				if strings.Contains(err.Error(), "404 Not Found") {
					return false, nil
				}
				return false, err
			}
			return membership.GetState() == "active", nil
		}
	}
	return false, nil
}
