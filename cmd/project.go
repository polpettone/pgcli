package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ProjectCmd(apiClient *GitlabAPIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "project",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := handleProjectCommand(args, apiClient)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func handleProjectCommand(args []string, apiClient *GitlabAPIClient) error {

	projects := viper.GetStringSlice("projects")
	fmt.Println("Projects")
	fmt.Println(projects)

	current_project := viper.GetString("current_project")

	fmt.Println("Current Project")
	fmt.Println(current_project)

	viper.Set("current_project", "dummy")

	current_project = viper.GetString("current_project")
	fmt.Println("Current Project")
	fmt.Println(current_project)

	foo := viper.GetBool("foo")
	fmt.Println(foo)
	fmt.Println("Set Foo to True")
	viper.Set("foo", true)
	foo = viper.GetBool("foo")
	fmt.Println(foo)

	viper.WriteConfig()

	return nil
}

func init() {
	initConfig()
	projectCmd := ProjectCmd(NewGitlabAPIClient())

	rootCmd.AddCommand(projectCmd)
}
