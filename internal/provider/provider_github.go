package provider

import (
	"context"
	"net/http"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Github provider
type Github struct {
	client *github.Client
}

// NewGithub constrcuts a new Github provider
func NewGithub(client *github.Client) (Provider, error) {
	prv := &Github{
		client: client,
	}

	return prv, nil
}

// GetUserStars returns the user's starred repositories
func (g *Github) GetUserStars(ctx context.Context, name string) ([]StarredRepository, error) {
	logger := logrus.WithFields(logrus.Fields{
		"logger":     "providers/Github.GetUserStars",
		"user.login": name,
	})

	logger.Info("getting user's starred repositories")

	stars := []StarredRepository{}

	currentPage := 1
	for currentPage != 0 {
		opts := &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    currentPage,
				PerPage: 100,
			},
		}

		moreRepos, res, err := g.client.Activity.ListStarred(ctx, name, opts)
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve user's stars")
		}

		logger.
			WithFields(logrus.Fields{
				"current_page":  currentPage,
				"count":         len(stars),
				"res.code":      res.StatusCode,
				"res.next_page": res.NextPage,
			}).
			Debug("got stars")

		for _, repo := range moreRepos {
			stars = append(stars, StarredRepository{repo.Repository.GetFullName(), repo.StarredAt.Unix()})
		}

		currentPage = res.NextPage
	}

	return stars, nil
}

// GetUserFollowers returnes the user's followers
func (g *Github) GetUserFollowers(ctx context.Context, name string) ([]string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"logger":     "providers/Github.GetUserFollowers",
		"user.login": name,
	})

	logger.Info("getting user's followers")

	followers := []string{}

	currentPage := 1
	for currentPage != 0 {
		opts := &github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		}

		moreFollowers, res, err := g.client.Users.ListFollowers(ctx, name, opts)
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve user's followers")
		}

		logger.
			WithFields(logrus.Fields{
				"current_page":  currentPage,
				"count":         len(followers),
				"res.code":      res.StatusCode,
				"res.next_page": res.NextPage,
			}).
			Debug("got followers")

		for _, follower := range moreFollowers {
			followers = append(followers, follower.GetLogin())
		}

		currentPage = res.NextPage
	}

	return followers, nil
}

// GetUserFollowees returnes the user's followees
func (g *Github) GetUserFollowees(ctx context.Context, name string) ([]string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"logger":     "providers/Github.GetUserFollowees",
		"user.login": name,
	})

	logger.Info("getting user's followers")

	followees := []string{}

	currentPage := 1
	for currentPage != 0 {
		opts := &github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		}

		moreFollowees, res, err := g.client.Users.ListFollowing(ctx, name, opts)
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve user's followees")
		}

		logger.
			WithFields(logrus.Fields{
				"current_page":  currentPage,
				"count":         len(followees),
				"res.code":      res.StatusCode,
				"res.next_page": res.NextPage,
			}).
			Debug("got followees")

		for _, follower := range moreFollowees {
			followees = append(followees, follower.GetLogin())
		}

		currentPage = res.NextPage
	}

	return followees, nil
}

// GetUserRepositories returns the user's repositories
func (g *Github) GetUserRepositories(ctx context.Context, name string) ([]string, error) {
	return nil, nil
}

// FollowUser follows a user give their login
func (g *Github) FollowUser(ctx context.Context, name string) error {
	logger := logrus.WithFields(logrus.Fields{
		"logger":     "providers/Github.FollowUser",
		"user.login": name,
	})

	logger.Info("following user")

	res, err := g.client.Users.Follow(ctx, name)
	if err != nil {
		return errors.Wrap(err, "could not follow user")
	}

	if res.StatusCode != http.StatusOK {
		logger.
			WithFields(logrus.Fields{
				"res.status":  res.StatusCode,
				"res.headers": res.Header,
			}).
			Warn("following user returned a non-ok status code")
	}

	return nil
}
