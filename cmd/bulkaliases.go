package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Alias struct {
	Name        string `mapstructure:"name" json:"name"`
	Description string `mapstructure:"description" json:"description"`
	Command     string `mapstructure:"command" json:"command"`
	Category    string `mapstructure:"category" json:"category"`
	Example     string `mapstructure:"example" json:"example"`
}

var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Manage Git and GitHub CLI aliases",
	Long: `Manage Git and GitHub CLI aliases.
			List, install, backup, and restore aliases.
			
			Important Notes:
			1. This command will modify your Git and GitHub CLI configurations.
			2. Git aliases will be set globally for all repositories.
			3. Some complex aliases may require manual addition to your .gitconfig file.
			4. Both 'git' and 'gh' (GitHub CLI) must be installed and accessible in your system PATH.`,
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %s\n", err)
			return
		}

		listAliases(aliases)
	},
}

var installCmd = &cobra.Command{
	Use:   "install [alias_name]",
	Short: "Install a specific alias or all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %s\n", err)
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		if all {
			if confirmInstallation() {
				installAliases(aliases)
			} else {
				fmt.Println("Installation cancelled.")
			}
			return
		}

		if len(args) == 0 {
			fmt.Println("Please specify an alias name or use --all to install all aliases.")
			return
		}

		aliasName := args[0]
		for _, alias := range aliases {
			if alias.Name == aliasName {
				installAlias(alias)
				return
			}
		}

		fmt.Printf("Alias '%s' not found.\n", aliasName)
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup [filename]",
	Short: "Backup current aliases to a file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "aliases_backup.json"
		if len(args) > 0 {
			filename = args[0]
		}

		aliases, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %s\n", err)
			return
		}

		data, err := json.MarshalIndent(aliases, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling aliases: %s\n", err)
			return
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("Error writing backup file: %s\n", err)
			return
		}

		fmt.Printf("Aliases backed up to %s\n", filename)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [filename]",
	Short: "Restore aliases from a backup file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "aliases_backup.json"
		if len(args) > 0 {
			filename = args[0]
		}

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading backup file: %s\n", err)
			return
		}

		var aliases []Alias
		err = json.Unmarshal(data, &aliases)
		if err != nil {
			fmt.Printf("Error unmarshaling aliases: %s\n", err)
			return
		}

		if confirmInstallation() {
			installAliases(aliases)
		} else {
			fmt.Println("Restoration cancelled.")
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new custom alias",
	Run: func(cmd *cobra.Command, args []string) {
		alias := promptForAlias()
		aliases, _ := loadAliases()
		aliases = append(aliases, alias)
		saveAliases(aliases)
		fmt.Printf("Alias '%s' added successfully.\n", alias.Name)
	},
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.AddCommand(installCmd, backupCmd, restoreCmd, addCmd)
	installCmd.Flags().Bool("all", false, "Install all aliases")
}

func loadAliases() ([]Alias, error) {
	var aliases []Alias
	err := viper.UnmarshalKey("aliases", &aliases)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal aliases from config: %w", err)
	}
	return aliases, nil
}

func saveAliases(aliases []Alias) error {
	viper.Set("aliases", aliases)
	return viper.WriteConfig()
}

func listAliases(aliases []Alias) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Name | cyan }} ({{ .Category | yellow }}): {{ .Description | faint }}",
		Inactive: "  {{ .Name | cyan }} ({{ .Category | yellow }}): {{ .Description | faint }}",
		Selected: "{{ .Name | red | cyan }}",
		Details: `
{{ "Name:" | faint }}	{{ .Name }}
{{ "Category:" | faint }}	{{ .Category }}
{{ "Description:" | faint }}	{{ .Description }}
{{ "Command:" | faint }}	{{ .Command }}
{{ "Example:" | faint }}	{{ .Example }}`,
	}

	prompt := promptui.Select{
		Label:     "Navigate the aliases (Press ESC to quit)",
		Items:     aliases,
		Size:      10,
		Templates: templates,
	}

	_, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("\nExiting...")
			return
		}
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func confirmInstallation() bool {
	fmt.Println("This will modify your Git and GitHub CLI configurations.")
	fmt.Println("Some complex aliases may require manual addition to your .gitconfig file.")
	fmt.Print("Do you want to proceed? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}

func installAliases(aliases []Alias) {
	for _, alias := range aliases {
		installAlias(alias)
	}
	fmt.Println("All aliases have been installed.")
	fmt.Println("Note: Some complex Git aliases may need to be manually added to your .gitconfig file.")
}

func installAlias(alias Alias) {
	if alias.Category == "Git" {
		gitConfigFile, err := getGitConfigPath()
		if err != nil {
			fmt.Printf("Error getting Git config file path: %s\n", err)
			return
		}
		err = addGitAlias(gitConfigFile, alias.Name, alias.Command)
		if err != nil {
			fmt.Printf("Error setting Git alias %s: %s\n", alias.Name, err)
		} else {
			fmt.Printf("Set Git alias %s\n", alias.Name)
		}
	} else if alias.Category == "GitHub" {
		cmd := exec.Command("gh", "alias", "set", alias.Name, alias.Command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error setting GitHub CLI alias %s: %s\n", alias.Name, err)
		} else {
			fmt.Printf("Set GitHub CLI alias %s: %s\n", alias.Name, strings.TrimSpace(string(output)))
		}
	} else {
		fmt.Printf("Unknown category for alias %s: %s\n", alias.Name, alias.Category)
	}
}

func getGitConfigPath() (string, error) {
	cmd := exec.Command("git", "config", "--global", "--list", "--show-origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.SplitN(lines[0], "\t", 2)
		if len(parts) == 2 {
			return strings.TrimPrefix(parts[0], "file:"), nil
		}
	}
	return "", fmt.Errorf("could not determine Git config file path")
}

func addGitAlias(configFile, name, command string) error {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	aliasSection := false
	aliasAdded := false

	for i, line := range lines {
		if line == "[alias]" {
			aliasSection = true
		} else if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if aliasSection {
				break
			}
			aliasSection = false
		}

		if aliasSection && strings.HasPrefix(line, name+"=") {
			lines[i] = name + "=" + command
			aliasAdded = true
			break
		}
	}

	if !aliasAdded {
		if !aliasSection {
			// Find the last occurrence of a section
			lastSectionIndex := -1
			for i, line := range lines {
				if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
					lastSectionIndex = i
				}
			}

			if lastSectionIndex == -1 {
				// No sections found, add [alias] at the end
				lines = append(lines, "", "[alias]")
			} else {
				// Insert [alias] after the last section
				newLines := make([]string, len(lines)+2)
				copy(newLines, lines[:lastSectionIndex+1])
				newLines[lastSectionIndex+1] = ""
				newLines[lastSectionIndex+2] = "[alias]"
				copy(newLines[lastSectionIndex+3:], lines[lastSectionIndex+1:])
				lines = newLines
			}
		}
		lines = append(lines, name+"="+command)
	}

	return ioutil.WriteFile(configFile, []byte(strings.Join(lines, "\n")), 0644)
}

func promptForAlias() Alias {
	namePrompt := promptui.Prompt{
		Label: "Alias Name",
	}
	name, _ := namePrompt.Run()

	categoryPrompt := promptui.Select{
		Label: "Category",
		Items: []string{"Git", "GitHub"},
	}
	_, category, _ := categoryPrompt.Run()

	descriptionPrompt := promptui.Prompt{
		Label: "Description",
	}
	description, _ := descriptionPrompt.Run()

	commandPrompt := promptui.Prompt{
		Label: "Command",
	}
	command, _ := commandPrompt.Run()

	examplePrompt := promptui.Prompt{
		Label: "Usage Example",
	}
	example, _ := examplePrompt.Run()

	return Alias{
		Name:        name,
		Category:    category,
		Description: description,
		Command:     command,
		Example:     example,
	}
}
