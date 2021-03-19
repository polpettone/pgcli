package models

import (
	"encoding/json"
	"fmt"
)

type Projects struct {
	Projects []Project
}

type Project struct {
	Id              int `json:"id"`
	Name            string
	SSH_url_to_repo string `json:"ssh_url_to_repo"`
}

func (p Project) NiceString() string {
	return fmt.Sprintf("%d \t %s \t %s", p.Id, p.Name, p.SSH_url_to_repo)
}

func (p Project) View() string {
	return fmt.Sprintf("%d \t %s \t %s", p.Id, p.Name, p.SSH_url_to_repo)
}

func ConvertJsonToProjects(jsonData []byte) (*[]Project, error) {
	projects := make([]Project, 0)
	err := json.Unmarshal(jsonData, &projects)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}
