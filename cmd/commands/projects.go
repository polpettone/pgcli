package commands

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/polpettone/pgcli/cmd/config"
	"github.com/spf13/cobra"
)

func init() {
	InitConfig()
	projectCmd := ProjectsCmd(adapter.NewGitlabAPIClient())
	rootCmd.AddCommand(projectCmd)
}

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

func handleProjectCommand(args []string, apiClient *adapter.GitlabAPIClient) (string, error) {
	if len(args) == 0 {
		return getProjects(apiClient)
	}

	state := config.State{
		CurrentProject: args[0],
	}

	err := config.WriteState(state, "/home/esteban/.config/pgcli/state.json")
	if err != nil {
		apiClient.Logging.ErrorLog.Printf("%v", err)
	}
	if err != nil {
		apiClient.Logging.ErrorLog.Printf("%s", err)
	}

	return fmt.Sprintf("Changed Project to %s", args[0]), nil
}



func getProjects(apiClient *adapter.GitlabAPIClient) (string, error) {
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


