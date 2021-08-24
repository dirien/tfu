package git

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v38/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
)

type Github struct {
	Owner   string
	Repo    string
	Version *version.Version
}

func ParseGithubInfos(source string) *Github {
	//git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0
	//https://github.com/rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement.git
	source = strings.TrimPrefix(source, "git@github.com:")
	source = strings.TrimPrefix(source, "https://github.com/")
	sources := strings.SplitAfter(source, "//")
	var ref string
	var currentVersion *version.Version
	if len(sources) > 1 {
		var err error
		ref = strings.TrimPrefix(sources[1], "?ref=")
		currentVersion, err = version.NewVersion(ref)
		if err != nil {
			return nil
		}
	} else {
		// no ref means it is already master
		return nil
	}
	//treim the hell out of it.
	source = strings.TrimSuffix(sources[0], "//")
	source = strings.TrimSuffix(source, ".git")
	return &Github{
		Owner:   strings.Split(source, "/")[0],
		Repo:    strings.Split(source, "/")[1],
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
