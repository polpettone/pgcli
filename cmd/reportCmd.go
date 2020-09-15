package cmd

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/cobra"
)

func NewReportCmd(apiClient APIClient) *cobra.Command{
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

func handleReportCommand(cobraCommand *cobra.Command, apiClient APIClient) (string, error) {
	allStatus := ""

	pipelines, err := apiClient.getPipelines(allStatus, false, 20)

	if err != nil {
		return "", err
	}

	report := models.NewReport(pipelines)

	return report.NiceString(), nil
}


func init() {
	reportCmd := NewReportCmd(gitlabAPIClient)
	rootCmd.AddCommand(reportCmd)

}
