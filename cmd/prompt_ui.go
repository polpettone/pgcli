package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/polpettone/pgcli/cmd/models"
	"os"
)

func showPipelineSelectionPrompt(pipelines []*models.Pipeline) (*models.Pipeline, error) {

	templates := &promptui.SelectTemplates{
		Active:   "\U00001B61 {{ .Id | red }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{.Duration | green}} {{ .Status | green }}",
		Inactive:  "{{ .Id | blue }} {{ .CreatedAt | blue }} {{ .UpdatedAt | blue }} {{.Duration | green}} {{ .Status | blue }}",
		Selected:   "\U00001B61 {{ .Id | green }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{.Duration | green}} {{ .Status | green }}",
	}

	prompt := promptui.Select{
		Label:     "Which Pipeline",
		Items:     pipelines,
		Templates: templates,
		Size:      20,
	}

	i, _, err := prompt.Run()

	if err != nil {

		if err == promptui.ErrInterrupt {
			os.Exit(1)
		}

		return nil, err
	}

	return pipelines[i], nil
}

func showJobSelectionPrompt(jobs []models.Job) (*models.Job, error) {

	templates := &promptui.SelectTemplates{
		Active:   "\U00001B61 {{ .Status | red }} {{ .StartedAt | green }} {{ .FinishedAt | green }} {{ .Duration | green }}  {{ .Name | green }} ",
		Inactive:   "{{ .Status | blue }}  {{ .StartedAt | blue}} {{ .FinishedAt | blue }} {{ .Duration | blue }} {{ .Name | blue }} ",
		Selected:   "\U00001B61 {{ .Status | red }} {{ .StartedAt | green }} {{ .FinishedAt | green }} {{ .Duration | green }} {{ .Name | green }} ",
	}

	prompt := promptui.Select{
		Label:     "Which Job",
		Items:     jobs,
		Templates: templates,
		Size:      20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(1)
		}
		return nil, err
	}

	return &jobs[i], nil
}
