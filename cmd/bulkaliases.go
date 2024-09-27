package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Alias struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Command     string `mapstructure:"command"`
}

var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List and optionally install useful Git and GitHub CLI aliases",
	Long: `List and optionally install useful Git and GitHub CLI aliases.
			Use --bulk to install all aliases on your machine.
			
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

		bulk, _ := cmd.Flags().GetBool("bulk")
		if bulk {
			if confirmInstallation() {
				installAliases(aliases)
			} else {
				fmt.Println("Installation cancelled.")
			}
		} else {
			listAliases(aliases)
		}
	},
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.Flags().Bool("bulk", false, "Install all aliases on your machine")
}

func loadAliases() ([]Alias, error) {
	var aliases []Alias
	err := viper.UnmarshalKey("aliases", &aliases)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal aliases from config: %w", err)
	}
	return aliases, nil
}

func listAliases(aliases []Alias) {
	// Display aliases using promptui, limiting the display to 5 commands at a time
	fmt.Println("Available aliases:")

	// Setup the promptui selector with a maximum of 5 items displayed at a time
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Name | cyan }}: {{ .Description | faint }}",
		Inactive: "  {{ .Name | cyan }}: {{ .Description | faint }}",
		Selected: "{{ .Name | red | cyan }}",
	}

	// Create the prompt using the aliases data
	prompt := promptui.Select{
		Label:     "Navigate the aliases (Press ESC to quit)",
		Items:     aliases,
		Size:      5, // Limit to displaying 5 items at a time
		Templates: templates,
	}

	// Run the prompt
	_, _, err := prompt.Run()

	// Handle error or ESC key pressed
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
	gitConfigFile, err := getGitConfigPath()
	if err != nil {
		fmt.Printf("Error getting Git config file path: %s\n", err)
		return
	}

	for _, alias := range aliases {
		if strings.HasPrefix(alias.Command, "!") {
			// Git alias
			err := addGitAlias(gitConfigFile, alias.Name, alias.Command)
			if err != nil {
				fmt.Printf("Error setting Git alias %s: %s\n", alias.Name, err)
			} else {
				fmt.Printf("Set Git alias %s\n", alias.Name)
			}
		} else {
			// GitHub CLI alias
			cmd := exec.Command("gh", "alias", "set", alias.Name, alias.Command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error setting GitHub CLI alias %s: %s\n", alias.Name, err)
			} else {
				fmt.Printf("Set GitHub CLI alias %s: %s\n", alias.Name, strings.TrimSpace(string(output)))
			}
		}
	}
	fmt.Println("All aliases have been installed.")
	fmt.Println("Note: Some complex Git aliases may need to be manually added to your .gitconfig file.")
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
	content, err := os.ReadFile(configFile)
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
			lines = append(lines, "", "[alias]")
		}
		lines = append(lines, name+"="+command)
	}

	return os.WriteFile(configFile, []byte(strings.Join(lines, "\n")), 0644)
}
