package github

import (
	"fmt"

	"github.com/shurcooL/githubv4"
)

type GHUser struct {
	User struct {
		Login         githubv4.String
		Name          githubv4.String
		Organizations struct {
			Nodes []struct {
				Name githubv4.String
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage githubv4.Boolean
			}
		} `graphql:"organizations(first: 99)"`
	} `graphql:"user(login: $login)"`
}

func (c Client) GetUser(username string) (*GHUser, error) {
	logger := c.log.With("username", username)
	logger.Debug("GetUser")
	variables := map[string]interface{}{
		"login": githubv4.String(username),
	}

	var user GHUser
	err := c.ghClient.Query(c.ctx, &user, variables)
	if err != nil {
		logger.Error("unable to fetch user", "error", err)
		return nil, fmt.Errorf("unable to fetch user: %w", err)
	}

	return &user, nil
}
