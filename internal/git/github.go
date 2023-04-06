package git

import (
	"context"
	"os"
	"strings"

	giturls "github.com/whilp/git-urls"

	"github.com/google/go-github/v51/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
)

type Github struct {
	Owner   string
	Repo    string
	Version string
}

func ParseGithubInfos(source string) *Github {
	// git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0
	// https://github.com/rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement.git
	parse, err := giturls.Parse(source)
	if err != nil {
		return nil
	}
	var currentVersion string
	if len(parse.Query()["ref"]) > 0 {
		currentVersion = parse.Query()["ref"][0]
	} else {
		// no ref means it is already master
		return nil
	}
	diff := 0
	if parse.Scheme == "ssh" {
		diff = 1
	}
	owner := strings.Split(parse.Path, "/")[1-diff]
	repo := strings.Split(parse.Path, "/")[2-diff]
	repo = strings.ReplaceAll(repo, ".git", "")
	return &Github{
		Owner:   owner,
		Repo:    repo,
		Version: currentVersion,
	}
}

func CheckGitTokenIsSet() bool {
	githubToken := os.Getenv("GIT_TOKEN")
	return len(githubToken) != 0
}

func (g *Github) GetLatestVersion() (*version.Version, error) {
	ctx := context.Background()
	githubToken := os.Getenv("GIT_TOKEN")
	if len(githubToken) != 0 {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		release, _, err := client.Repositories.ListReleases(ctx, g.Owner, g.Repo, nil)
		if err != nil {
			return nil, err
		}
		latest, _ := version.NewVersion(*release[0].TagName)
		for _, repositoryRelease := range release {
			tmpRelease, _ := version.NewVersion(*repositoryRelease.TagName)
			if tmpRelease.GreaterThan(latest) {
				latest = tmpRelease
			}
		}
		return latest, nil
	}
	return nil, nil
}
