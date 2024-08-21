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
	"os"

	"github.com/gin-gonic/gin"
	goGithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
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

func (a *Admin) handleGitHubLogin(c *gin.Context) {
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
        c.Redirect(http.StatusSeeOther, "/dashboard")
    } else {
        log.Printf("User %s is not a maintainer in %s organization", user.GetLogin(), orgName)
        c.String(http.StatusForbidden, "You must be a maintainer in the %s organization to access this page.", orgName)
    }
}

func (a *Admin) checkUserOrgMembership(client *goGithub.Client, username, orgName string) (bool, error) {
    membership, _, err := client.Organizations.GetOrgMembership(context.Background(), username, orgName)
    if err != nil {
        log.Printf("Failed to get org membership for user %s: %v", username, err)
        return false, err
    }

    log.Printf("User %s role in %s: %s", username, orgName, membership.GetRole())
    return membership.GetRole() == "admin" || membership.GetRole() == "maintainer", nil
}
