package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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
		return "", fmt.Errorf("reading token file: %w", err)
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
Welcome, %s! Let's create some milestones and issues.
`, username)
}

func printHelp() {
	fmt.Printf(`
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

For more information, visit: https://github.com/igorcosta/gh-lazy
`)
}
