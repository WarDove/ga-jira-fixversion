# Jira Release Automation

This GitHub Action automates the process of creating a Jira release version from GitHub workflows, leveraging the robustness of Golang. It's designed to streamline your release management by automatically updating Jira with version details when you push to specific branches.

## Features

- **Automated Jira Releases:** Automatically create Jira release versions when pushing to release branches.
- **Docker-based:** Runs in a Docker container, ensuring consistency across environments.
- **Customizable:** Easily configure through GitHub workflow definitions.

## Getting Started

### Prerequisites

- Jira Account
- GitHub Account and Repository
- Docker installed (for local testing)

### Setup

1. **Jira API Token:** Generate an API token from your Jira account to authenticate requests.
2. **GitHub Secrets:** Store your Jira credentials (`JIRA_USER`, `JIRA_AUTOMATION_TOKEN`, `JIRA_URL`) as secrets in your GitHub repository.

### Usage

Add a workflow file to your repository (e.g., `.github/workflows/jira-release.yml`) with the following content:

```yaml
name: Jira Release Action
run-name: "Jira Release Action | triggered by @${{ github.actor }}"

on:
  push:
    branches:
      - 'release/**'

jobs:
  jira-release:
    if: github.event.created
    runs-on: ubuntu-latest
    steps:

      - name: Extract branch name
        shell: bash
        run: echo "BRANCH_NAME=${GITHUB_HEAD_REF:-${GITHUB_REF#refs/heads/}}" >> $GITHUB_ENV
        id: extract_branch

      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: Run Jira Release Action
        uses: wardove/ga-jira-fixversion@v1.0
        with:
          versionName: ${{ env.BRANCH_NAME }}
          projectKey: YOUR_JIRA_PROJECT_KEY
          jiraUser: ${{ secrets.JIRA_AUTOMATION_USER }}
          jiraToken: ${{ secrets.JIRA_AUTOMATION_TOKEN }}
          jiraUrl: https://your_jira_subdomain.atlassian.net

```

For manual test workflow:

```yaml
name: Test Jira Release Action
run-name: "Jira Release Action | triggered manually by @${{ github.actor }}"

on:
  workflow_dispatch:
    inputs:
      versionName:
        description: 'Name of the version to create'
        required: true
      projectKey:
        description: 'Jira project key'
        required: true
      jiraUser:
        description: 'Jira user email'
        required: true
      jiraToken:
        description: 'Jira API token'
        required: true
      jiraUrl:
        description: 'Jira instance URL'
        required: true

jobs:
  test-jira-release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: Run Jira Release Action
        uses: ./
        with:
          versionName: ${{ github.event.inputs.versionName }}
          projectKey: ${{ github.event.inputs.projectKey }}
          jiraUser: ${{ github.event.inputs.jiraUser }}
          jiraToken: ${{ github.event.inputs.jiraToken }}
          jiraUrl: ${{ github.event.inputs.jiraUrl }}

```

#### Source repository: https://github.com/WarDove/ga-jira-fixversion
