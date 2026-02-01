package postmaster

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GitHubSender struct {
	Client *github.Client
	Owner  string
	Repo   string
}

func NewGitHubSender(token string, owner string, repo string) *GitHubSender {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubSender{
		Client: github.NewClient(tc),
		Owner:  owner,
		Repo:   repo,
	}
}

func (g *GitHubSender) Send(destination string, body string) error {
	issueNum, err := strconv.Atoi(destination)
	if err != nil {
		return fmt.Errorf("invalid issue number: %s", destination)
	}

	ctx := context.Background()
	comment := &github.IssueComment{Body: &body}
	_, _, err = g.Client.Issues.CreateComment(ctx, g.Owner, g.Repo, issueNum, comment)
	return err
}