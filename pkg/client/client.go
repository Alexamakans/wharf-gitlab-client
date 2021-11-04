package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alexamakans/wharf-common-api-client/pkg/remoteprovider"
	"github.com/iver-wharf/wharf-core/pkg/problem"
	"github.com/xanzy/go-gitlab"
)

// Client implements remoteprovider.Client.
type Client struct {
	remoteprovider.BaseClient
}

func (c *Client) FetchFile(projectIdentifier remoteprovider.ProjectIdentifier, fileName string) ([]byte, error) {
	return []byte{}, nil
}

func (c *Client) FetchBranches(projectIdentifier remoteprovider.ProjectIdentifier) ([]remoteprovider.WharfBranch, error) {
	return []remoteprovider.WharfBranch{}, nil
}

func (c *Client) FetchProjectByGroupAndProjectName(groupName, projectName string) (remoteprovider.WharfProject, error) {
	gitlabClient, err := gitlab.NewClient(c.Token, gitlab.WithBaseURL(c.RemoteProviderURL))
	if err != nil {
		return remoteprovider.WharfProject{}, fmt.Errorf("failed connecting to GitLab using URL %q: %w", c.RemoteProviderURL, err)
	}

	urlEncodedProjectPath := fmt.Sprintf("%s/%s", groupName, projectName)
	gitlabProject, resp, err := gitlabClient.Projects.GetProject(urlEncodedProjectPath, nil)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return remoteprovider.WharfProject{}, problem.Response{
				Type:   "/prob/api/remote-project-not-found",
				Title:  "Project not found at remote provider.",
				Detail: fmt.Sprintf("Could not find project matching %s/%s at %s", groupName, projectName, c.RemoteProviderURL),
				Status: http.StatusNotFound,
				Errors: []string{err.Error()},
			}
		}

		return remoteprovider.WharfProject{}, err
	}

	var project remoteprovider.WharfProject
	project.RemoteProjectID = strconv.Itoa(gitlabProject.ID)
	project.GitURL = gitlabProject.SSHURLToRepo
	project.Name = gitlabProject.Name
	project.GroupName = gitlabProject.Namespace.Name

	return project, nil
}

func (c *Client) WharfProjectToIdentifier(project remoteprovider.WharfProject) remoteprovider.ProjectIdentifier {
	return remoteprovider.ProjectIdentifier{
		Values: []string{project.RemoteProjectID},
	}
}
