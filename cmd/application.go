package cmd

import "os"


var gitlabAPIClient = NewGitlabAPIClient(
	os.Getenv("GITLAB_API_TOKEN"),
	os.Getenv("GITLAB_PROJECT_URL"),
	os.Getenv("GITLAB_PROJECT_ID"))
