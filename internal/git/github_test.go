package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	SSHGitRepoURL   = "git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0"
	HTTPSGitRepoURL = "https://github.com/rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement.git?ref=v0.12.0"
)

func TestGetSourceSSH(t *testing.T) {
	expected := &Github{
		Owner:   "rackspace-infrastructure-automation",
		Repo:    "aws-terraform-asg_instance_replacement",
		Version: "v0.12.0",
	}
	github := ParseGithubInfos(SSHGitRepoURL)

	require.Equal(t, expected, github)
}

func TestGetSourceHTTPS(t *testing.T) {
	expected := &Github{
		Owner:   "rackspace-infrastructure-automation",
		Repo:    "aws-terraform-asg_instance_replacement",
		Version: "v0.12.0",
	}
	github := ParseGithubInfos(HTTPSGitRepoURL)

	require.Equal(t, expected, github)
}
