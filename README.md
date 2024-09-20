# Lazy

Whenever I have a new project for implementing GitHub in enterprise levels, setting up engagement issues and milestones are the first things to do. But I'm lazy, and I like to populate quickly. This project helps automate that process.

## Overview

Lazy is a tool designed to streamline the setup of GitHub projects, issues, and milestones for enterprise-level implementations. It allows you to quickly populate your GitHub repository with predefined tasks and configurations, saving you time and effort.

## How to Use


To run the Lazy tool, use the following command:

```sh
brew install gh
gh extension install lazy
gh lazy -reponame "username/repository" -token "your_github_token" -tasks "tasks/tasks.json"

### Options available

Lazy - GitHub Milestone and Issue Creator

Usage: gh lazy [flags]

Flags:
  -reponame string    The repository name (e.g., 'username/repo')
  -tasks string       Path to the tasks JSON file
  -tokenfile string   Path to the file containing the GitHub token (default ".token")

Example:
  gh lazy -reponame user/repo -tasks ./tasks.json

Description:
  This tool automates the creation of milestones and issues in a GitHub repository
  based on a JSON file, and adds the issues to a newly created GitHub Project (v2) using the gh CLI.