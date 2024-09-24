package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/igorcosta/gh-lazy/pkg/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func LogError(err error, message string) {
	log.WithError(err).Error(message)
}

func LogInfo(message string) {
	log.Info(message)
}

func WrapError(err error, message string) error {
	return errors.Wrap(err, message)
}

func ReadTokenFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("opening token file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "GH_TOKEN=") {
			return strings.TrimPrefix(line, "GH_TOKEN="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading token file: %w", err)
	}

	return "", fmt.Errorf("GH_TOKEN not found in file")
}

func LoadTasksFile(filePath string) (*models.TasksFile, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading tasks JSON file: %w", err)
	}

	var tasksFile models.TasksFile
	if err := json.Unmarshal(file, &tasksFile); err != nil {
		return nil, fmt.Errorf("parsing tasks JSON: %w", err)
	}

	return &tasksFile, nil
}

func ShowProgress(progressChan <-chan string) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for message := range progressChan {
		fmt.Printf("\r%s %s", spinner[i], message)
		i = (i + 1) % len(spinner)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
}

func PrintWelcome(username string) {
	fmt.Printf(`
██╗      █████╗ ███████╗██╗   ██╗
██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝
██║     ███████║  ███╔╝  ╚████╔╝ 
██║     ██╔══██║ ███╔╝    ╚██╔╝  
███████╗██║  ██║███████╗   ██║   
╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   
Welcome, %s! Let's automate your GitHub projects.
`, username)
}

func PrintHelp() {
	fmt.Printf(`
Lazy - GitHub Project, Milestone, and Issue Manager

Usage: gh lazy [command] [flags]

Commands:
  create    Create projects, milestones, and issues
  nuke      Delete a GitHub project and optionally all linked issues

Flags:
  -r, --repo string         The repository name (e.g., 'username/repo')
  -t, --tasks string        Path to the tasks JSON file
  -f, --token-file string   Path to the file containing the GitHub token (default ".token")

Use "gh lazy [command] --help" for more information about a command.

Examples:
  gh lazy create --repo user/repo --tasks ./tasks.json
  gh lazy nuke --projectid https://github.com/users/yourusername/projects/1 --all --dry-run

Description:
  This tool automates the creation and deletion of projects, milestones, and issues in a GitHub repository.

For more information, visit: https://github.com/igorcosta/gh-lazy
`)
}

func ParseProjectID(input string) (string, error) {
	// If the input is a URL, extract the project number
	if strings.HasPrefix(input, "http") {
		// Match patterns like /projects/1 or /projects/1/
		re := regexp.MustCompile(`/projects/(\d+)/?$`)
		matches := re.FindStringSubmatch(input)
		if len(matches) < 2 {
			return "", fmt.Errorf("invalid project URL")
		}
		return matches[1], nil
	}
	// Otherwise, assume it's a project number
	return input, nil
}

func GetGitHubCLIToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub CLI token: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetToken(tokenFile string) (string, error) {
	// First, try to get the token from the gh CLI
	token, err := GetGitHubCLIToken()
	if err == nil && token != "" {
		return token, nil
	}

	// If gh CLI token is not available and tokenFile is provided, try to read from file
	if tokenFile != "" {
		token, err := ReadTokenFromFile(tokenFile)
		if err == nil && token != "" {
			return token, nil
		}
		// Only return an error if a token file was specified but couldn't be read
		return "", fmt.Errorf("failed to read token from specified file: %w", err)
	}

	// If we've reached this point, no valid token was found
	return "", fmt.Errorf("GitHub token not found. Please authenticate with 'gh auth login' or provide a token file using the -f flag")
}

func PrintUserGuide() {
	fmt.Println(`To use gh lazy, please ensure you have:
1. Authenticated with GitHub CLI using 'gh auth login', or
2. Provided a token file using the -f or --token-file flag

Usage examples:
  gh lazy create --repo username/repo --tasks tasks.json
  gh lazy nuke --projectid <project_id_or_url> --all --dry-run

For more information, run: gh lazy --help`)
}
