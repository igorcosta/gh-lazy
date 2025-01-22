# Lazy: Your Productive Git and GitHub CLI Superpowers

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome)
[![GitHub stars](https://img.shields.io/github/stars/igorcosta/gh-lazy.svg)](https://github.com/igorcosta/gh-lazy/stargazers)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

Lazy is a tool designed to enhance your productivity by automating git and GitHub CLI commands. It simplifies the process of setting up and managing GitHub projects, issues, and milestones.

## Features

- Efficiently create and delete GitHub issues and milestones
- Automate the setup of GitHub Projects (v2)
- Customize task templates using JSON
- Delete GitHub projects and issues with ease
- Seamlessly integrate with GitHub CLI
- Private LLM interface for Ollama or Llama-3.2 2b param

## Requirements

Before using Lazy, ensure you have the following:

1. Homebrew installed
2. A GitHub account
3. A valid GitHub token with appropriate permissions
4. Basic knowledge of JSON
5. Optional: If you have Ollama installed, configure an LLM for additional features

## Installation

To install Lazy, follow these steps:

1. Open your terminal.
2. Install the GitHub CLI:

   ```bash
   brew install gh
   ```

3. Install the Lazy extension:

   ```bash
   gh extension install lazy
   ```

## Usage

### Creating Projects, Milestones, and Issues

```bash
gh lazy create --repo "your-username/your-repo" --tasks "path/to/tasks.json"
```

#### Available Options for `create`

```bash
Usage: gh lazy create [flags]

Flags:
  -r, --repo string         The repository name (e.g., 'username/repo')
  -t, --tasks string        Path to the tasks JSON file
  -f, --token-file string   Path to the file containing your GitHub token (default ".token")

Example:
  gh lazy create --repo username/repo --tasks ./tasks.json
```

### Deleting a Project

To delete a GitHub project and optionally all linked issues, use the following command:

```bash
gh lazy nuke [--projectid <project_id_or_url>] [--all] [--dry-run]
```

- If you provide the `--projectid` (`-p`) flag, the command will delete the specified project.
- If you omit the `--projectid` flag, the tool will:

  1. List all your available projects and allow you to select one interactively.
  2. Ask if you want to perform a dry run first.
  3. Ask if you want to delete all associated issues.

#### Available Options for `nuke`

```bash
Usage: gh lazy nuke [flags]

Flags:
  -p, --projectid string   Project ID or URL to delete
  -a, --all                Delete all issues linked to the project
      --dry-run            Show what would happen without making changes

Example:
  gh lazy nuke --projectid https://github.com/users/username/projects/1 --all --dry-run
```

**Examples:**

- **Interactive Mode:**

  ```bash
  gh lazy nuke
  ```

  This will prompt you to select a project and configure options interactively.

- **Dry Run Without Deleting Issues:**

  ```bash
  gh lazy nuke --projectid https://github.com/users/username/projects/1 --dry-run
  ```

- **Dry Run With Deleting Issues:**

  ```bash
  gh lazy nuke --projectid 1 --all --dry-run
  ```

- **Actual Deletion:**

  ```bash
  gh lazy nuke --projectid 1 --all
  ```

- **Preparing a Prompt for LLM:**

```bash
gh lazy codeprompt "given this project, I need to modify my version.go file, help me out" --system-prompt . --ignore-gitignore --ignore "go.sum" --ignore "*.md" --ignore "gh-lazy" -o prompt.txt 
```

## Contributing

To contribute to Lazy, follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

For support, open an issue in the GitHub repository.
