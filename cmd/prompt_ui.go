package cmd

import (
	"github.com/manifoldco/promptui"
)

func showPipelineSelectionPrompt(pipelines []GitlabPipeline) (*GitlabPipeline, error) {

	templates := &promptui.SelectTemplates{
		Active:   "\U00001B61 {{ .Id | red }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{ .Status | green }}",
		Inactive:  "{{ .Id | blue }} {{ .CreatedAt | blue }} {{ .UpdatedAt | blue }} {{ .Status | blue }}",
		Selected:   "\U00001B61 {{ .Id | green }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{ .Status | green }}",
	}

	prompt := promptui.Select{
		Label:     "Which Pipeline",
		Items:     pipelines,
		Templates: templates,
		Size:      20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return &pipelines[i], nil
}

func showJobSelectionPrompt(jobs []GitlabJob) (*GitlabJob, error) {

	templates := &promptui.SelectTemplates{
		Active:   "\U00001B61 {{ .Name | red }} {{ .Status | green }} {{ .StartedAt | green }} {{ .FinishedAt | green }} {{ .Duration | green }} ",
		Inactive:   "{{ .Name | blue }} {{ .Status | blue }} {{ .StartedAt | blue}} {{ .FinishedAt | blue }} {{ .Duration | blue }} ",
		Selected:   "\U00001B61 {{ .Name | red }} {{ .Status | green }} {{ .StartedAt | green }} {{ .FinishedAt | green }} {{ .Duration | green }} ",
	}

	prompt := promptui.Select{
		Label:     "Which Job",
		Items:     jobs,
		Templates: templates,
		Size:      20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return &jobs[i], nil
}
