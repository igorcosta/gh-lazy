package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputFile      string
	ignorePatterns  []string
	includeHidden   bool
	ignoreGitignore bool
	xmlOutput       bool
	systemPrompt    bool
	globalIndex     int = 1
	cfgFile         string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(codepromptCmd)

	codepromptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output to a file instead of stdout")
	codepromptCmd.Flags().StringSliceVar(&ignorePatterns, "ignore", []string{}, "List of patterns to ignore")
	codepromptCmd.Flags().BoolVar(&includeHidden, "include-hidden", false, "Include files and folders starting with .")
	codepromptCmd.Flags().BoolVar(&ignoreGitignore, "ignore-gitignore", false, "Ignore .gitignore files and include all files")
	codepromptCmd.Flags().BoolVar(&xmlOutput, "cxml", false, "Output in XML format similar to Claude's long context window")
	codepromptCmd.Flags().BoolVar(&systemPrompt, "system-prompt", false, "Include system prompt from the configured file")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Error reading config file:", err)
	}
}

var codepromptCmd = &cobra.Command{
	Use:   "codeprompt [flags] <prompt> <path>",
	Short: "Generate a code prompt from files",
	Long:  `Generate a code prompt from files, similar to the files-to-prompt Python script.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runCodeprompt,
}

func runCodeprompt(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	path := args[1]

	var output *bufio.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		output = bufio.NewWriter(file)
	} else {
		output = bufio.NewWriter(os.Stdout)
	}
	defer output.Flush()

	// Read system prompt content
	goExpertContent, err := readSystemPromptFile()
	if err != nil {
		return fmt.Errorf("failed to read the system prompt file: %w", err)
	}

	// Replace {USER_REQUESTED} with the user's prompt
	re := regexp.MustCompile(`{\s*USER_REQUESTED\s*}`)
	combinedPrompt := re.ReplaceAllString(goExpertContent, prompt)

	// Write the combined prompt as the first content
	fmt.Fprintln(output, combinedPrompt)
	fmt.Fprintln(output, "---")

	// Add project structure section
	fmt.Fprintln(output, "Here is the project structure")
	fmt.Fprintln(output, "---")

	// Generate and write project structure
	if err := writeProjectStructure(output, path); err != nil {
		return fmt.Errorf("failed to write project structure: %w", err)
	}

	fmt.Fprintln(output, "---")

	// Write the combined prompt as the first content
	if xmlOutput {
		fmt.Fprintf(output, "<documents>\n<document index=\"%d\">\n<source>combined_prompt</source>\n<document_content>\n%s\n</document_content>\n</document>\n", globalIndex, combinedPrompt)
		globalIndex++
	} else {
		fmt.Fprintln(output, combinedPrompt)
		fmt.Fprintln(output, "---")
	}

	gitignoreRules := []string{}
	if !ignoreGitignore {
		gitignoreRules = readGitignore(path)
	}

	// Handle the root path (current directory) explicitly
	fmt.Printf("Processing directory: %s\n", path)

	err = filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Special handling for the current directory `.` to avoid treating it as hidden
		if filePath == path {
			fmt.Printf("Processing root directory: %s\n", filePath)
		} else if !includeHidden && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				fmt.Printf("Skipping hidden directory: %s\n", filePath)
				return filepath.SkipDir
			}
			fmt.Printf("Skipping hidden file: %s\n", filePath)
			return nil
		}

		if d.IsDir() {
			if !ignoreGitignore {
				gitignoreRules = append(gitignoreRules, readGitignore(filePath)...)
			}
			return nil
		}

		if shouldIgnore(filePath, gitignoreRules) {
			fmt.Printf("Ignoring file: %s\n", filePath)
			return nil
		}

		if isBinaryFile(filePath) {
			fmt.Printf("Skipping binary file: %s\n", filePath)
			return nil
		}

		fmt.Printf("Processing file: %s\n", filePath)
		return appendFileContents(filePath, output, xmlOutput)
	})

	if err != nil {
		return fmt.Errorf("error walking the path %s: %w", path, err)
	}

	if xmlOutput {
		fmt.Fprintln(output, "</documents>")
	}

	return nil
}

func writeProjectStructure(w io.Writer, root string) error {
	var dirs []bool
	var dirCount, fileCount int

	fmt.Fprintf(w, ".\n")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		depth := len(strings.Split(rel, string(os.PathSeparator))) - 1

		// Adjust dirs slice
		if depth >= len(dirs) {
			dirs = append(dirs, false)
		} else {
			dirs = dirs[:depth+1]
			dirs[depth] = false
		}

		// Prepare the prefix
		var prefix string
		for i, isLast := range dirs[:depth] {
			if i == depth-1 {
				if isLast {
					prefix += "└── "
				} else {
					prefix += "├── "
				}
			} else if isLast {
				prefix += "    "
			} else {
				prefix += "│   "
			}
		}

		// Print the entry
		fmt.Fprintf(w, "%s%s\n", prefix, filepath.Base(path))

		// Update counters
		if info.IsDir() {
			dirCount++
			dirs[depth] = true
		} else {
			fileCount++
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Print summary
	fmt.Fprintf(w, "\n%d directories, %d files\n", dirCount, fileCount)
	return nil
}

func readSystemPromptFile() (string, error) {
	systemPromptFile := viper.GetString("llm.systemprompt")
	if systemPromptFile == "" {
		return "", fmt.Errorf("system prompt file not specified in config")
	}

	absPath, err := filepath.Abs(systemPromptFile)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path of system prompt file: %w", err)
	}

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("system prompt file not found at %s", absPath)
		}
		return "", fmt.Errorf("failed to read system prompt file at %s: %w", absPath, err)
	}
	return string(content), nil
}
func shouldIgnore(path string, gitignoreRules []string) bool {
	baseName := filepath.Base(path)

	// Check ignore patterns provided by the user
	for _, pattern := range ignorePatterns {
		matched, err := filepath.Match(pattern, baseName)
		if err != nil {
			fmt.Printf("Error matching pattern %s: %v\n", pattern, err)
			continue
		}
		if matched {
			fmt.Printf("Ignoring due to pattern: %s\n", pattern)
			return true
		}
	}

	// Check rules from .gitignore
	for _, rule := range gitignoreRules {
		matched, err := filepath.Match(rule, baseName)
		if err != nil {
			fmt.Printf("Error matching gitignore rule %s: %v\n", rule, err)
			continue
		}
		if matched {
			fmt.Printf("Ignoring due to gitignore rule: %s\n", rule)
			return true
		}
	}

	return false
}

func isBinaryFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return true
	}
	defer file.Close()

	// Read the first 512 bytes
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return true
	}

	// Check if the file contains null bytes (common in binary files)
	return strings.Contains(string(buffer), "\x00")
}

// Print and append file contents, supporting both default and XML-like output
func appendFileContents(path string, output *bufio.Writer, xml bool) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Only use XML if explicitly requested via the --cxml flag
	if xml {
		printAsXML(output, path, string(content))
	} else {
		printDefault(output, path, string(content))
	}

	output.Flush() // Ensure that output is flushed after writing the content
	return nil
}

func printDefault(writer *bufio.Writer, path, content string) {
	fmt.Fprintf(writer, "\n%s\n", path)
	fmt.Fprintln(writer, "---")
	fmt.Fprintln(writer, content)
	fmt.Fprintln(writer, "---")
}

func printAsXML(writer *bufio.Writer, path, content string) {
	writer.WriteString(fmt.Sprintf("<document index=\"%d\">\n", globalIndex))
	writer.WriteString(fmt.Sprintf("<source>%s</source>\n", path))
	writer.WriteString("<document_content>\n")
	writer.WriteString(content + "\n")
	writer.WriteString("</document_content>\n")
	writer.WriteString("</document>\n")
	globalIndex++
}

// Read gitignore rules from the given directory
func readGitignore(dir string) []string {
	gitignorePath := filepath.Join(dir, ".gitignore")
	content, err := ioutil.ReadFile(gitignorePath)
	if err != nil {
		return []string{}
	}

	var rules []string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			rules = append(rules, line)
		}
	}

	return rules
}
