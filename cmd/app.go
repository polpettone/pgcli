package cmd

import (
	"github.com/spf13/viper"
)

var gitlabAPIClient APIClient

func init() {

	initConfig()

	gitlabAPIClient = NewGitlabAPIClient(
		viper.GetString("api_token"),
		viper.GetString("url"),
		viper.GetString("project_id"))
}
