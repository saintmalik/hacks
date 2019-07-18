package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getClient() (*github.Client, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		return nil, errors.New("GITHUB_TOKEN not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, nil
}

func main() {

	flag.Parse()

	user := flag.Arg(0)
	if user == "" {
		fmt.Fprintf(os.Stderr, "usage: ghrepos <username>\n")
		return
	}

	client, err := getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting client: %s\n")
		return
	}

	repos, err := getRepos(client, user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting repos: %s\n")
		// There might still be repos we can print here...
	}

	for _, repo := range repos {
		if *repo.Fork {
			continue
		}

		fmt.Println(*repo.CloneURL)
	}

}

func getRepos(client *github.Client, user string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		repos, resp, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			return allRepos, fmt.Errorf("failed to list repos: %s", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil

}
