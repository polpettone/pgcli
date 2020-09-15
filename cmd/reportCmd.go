package cmd

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/cobra"
)

func NewReportCmd(apiClient *GitlabAPIClient) *cobra.Command{
	return &cobra.Command{
		Use:   "report",
		Short: "shows a report",
		Long:  "shows a report",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleReportCommand(cmd, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleReportCommand(cobraCommand *cobra.Command, apiClient *GitlabAPIClient) (string, error) {
	allStatus := ""

	pipelines, err := apiClient.getPipelines(allStatus, false, 20, false)

	if err != nil {
		return "", err
	}

	report := models.NewReport(pipelines)

	return report.NiceString(), nil
}


func init() {
	reportCmd := NewReportCmd(NewGitlabAPIClient())
	rootCmd.AddCommand(reportCmd)

}
