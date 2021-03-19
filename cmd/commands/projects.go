package commands

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ProjectsCmd(apiClient *adapter.GitlabAPIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, err := handleProjectCommand(args, apiClient)
			if err != nil {
				return err
			}
			fmt.Println(stdout)
			return nil
		},
	}
}

func tryConfigWrite() {
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
}

func handleProjectCommand(args []string, apiClient *adapter.GitlabAPIClient) (string, error) {

	projects, err := apiClient.GetProjects()

	if err != nil {
		return "", nil
	}

	value := ""
	for _, project := range projects {
		value = value + "\n" + project.NiceString()
	}

	return value, nil
}

func init() {
	InitConfig()
	projectCmd := ProjectsCmd(adapter.NewGitlabAPIClient())
	rootCmd.AddCommand(projectCmd)
}
