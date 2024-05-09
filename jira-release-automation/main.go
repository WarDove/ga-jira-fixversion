package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"strconv"
)

var jiraClient *jira.Client
var jiraBaseUrl string

func main() {
	// Setup JIRA client
	if err := setupJiraClient(); err != nil {
		log.Fatalf("Error setting up Jira client: %v", err)
	}

	// Inputs from GitHub Actions
	versionName := os.Getenv("VERSION_NAME")
	projectKey := os.Getenv("PROJECT_KEY")

	// Get Project ID from Project Key
	projectID, err := getProjectID(projectKey)
	if err != nil {
		log.Fatalf("Error getting project ID: %v", err)
	}

	if versionExists, err := validateJiraVersion(projectKey, versionName); err == nil {
		if versionExists {
			fmt.Printf("The Jira version: %s already exists.", versionName)
			return
		} else {
			versionID, versionURL, err := createJiraVersion(versionName, projectID)
			if err != nil {
				log.Fatalf("Error creating Jira version: %v", err)
			}

			// Open the GitHub environment file
			envFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"), os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			defer envFile.Close()

			// Write the environment variables
			fmt.Fprintf(envFile, "VERSION_ID=%s\n", versionID)
			fmt.Fprintf(envFile, "VERSION_URL=%s\n", versionURL)
		}
	}
}

func getProjectID(projectKey string) (int, error) {
	req, err := jiraClient.NewRequest("GET", fmt.Sprintf("/rest/api/2/project/%s", projectKey), nil)
	if err != nil {
		return 0, fmt.Errorf("creating request failed: %v", err)
	}

	var project struct {
		ID string `json:"id"`
	}

	_, err = jiraClient.Do(req, &project)
	if err != nil {
		return 0, fmt.Errorf("request failed: %v", err)
	}

	// Convert project ID from string to integer
	projectID, convErr := strconv.Atoi(project.ID)
	if convErr != nil {
		return 0, fmt.Errorf("failed to convert project ID to integer: %v", convErr)
	}
	return projectID, nil
}

func setupJiraClient() error {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		return err
	}
	jiraClient = client
	return nil
}

func getJiraVersions(projectKey string) ([]jira.Version, error) {
	apiEndpoint := fmt.Sprintf("%s/rest/api/3/project/%s/versions", os.Getenv("JIRA_URL"), projectKey)
	req, err := jiraClient.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	var versions []jira.Version
	_, err = jiraClient.Do(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %v", err)
	}
	return versions, nil
}

func validateJiraVersion(projectKey string, versionName string) (bool, error) {
	versions, err := getJiraVersions(projectKey)
	if err != nil {
		log.Fatal("Jira API error, cannot get versions")
		return false, err
	}
	for _, version := range versions {
		if version.Name == versionName {
			return true, nil
		}
	}
	return false, nil
}

func createJiraVersion(versionName string, projectID int) (string, string, error) {
	released := false
	version := jira.Version{
		Name:      versionName,
		ProjectID: projectID,
		Released:  &released,
	}

	createdVersion, resp, err := jiraClient.Version.Create(&version)
	if err != nil {
		return "", "", fmt.Errorf("failed to create Jira version: %s, error: %v", resp.Status, err)
	}

	versionURL := fmt.Sprintf("%v/projects/%s/versions/%s", os.Getenv("JIRA_URL"), projectID, createdVersion.ID)
	return createdVersion.ID, versionURL, nil
}
