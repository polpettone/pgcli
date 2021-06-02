package commands

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/polpettone/pgcli/cmd/config"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func init() {
	InitConfig()
	projectCmd := ProjectsCmd(adapter.NewApp())
	rootCmd.AddCommand(projectCmd)
}

func ProjectsCmd(apiClient *adapter.App) *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			 handleProjectCommand(args, apiClient)
			 return nil
		},
	}
}

func handleProjectCommand(args []string, apiClient *adapter.App) {
	if len(args) == 0 {
		getProjects(apiClient)
		return
	}

	state := config.State{
		CurrentProject: args[0],
	}

	err := config.WriteState(state, "/home/esteban/.config/pgcli/state.json")
	if err != nil {
		config.Log.ErrorLog.Printf("%v", err)
	}
	if err != nil {
		config.Log.ErrorLog.Printf("%s", err)
	}

	fmt.Printf("Changed Project to %s", args[0])
}



func getProjects(app *adapter.App) {
	projects, err := app.GetProjects()

	if err != nil {
		config.Log.ErrorLog.Printf("%v", err)
		return
	}

	value := ""
	for _, project := range projects {
		value = value + "\n" + project.NiceString()
	}

	var data [][]string
	for _, p := range projects {
		data = append(data, []string {
			strconv.Itoa(p.Id),
			p.Name,
			p.SSH_url_to_repo,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	for _, d := range data {
		table.Append(d)
	}
	table.Render()
}


