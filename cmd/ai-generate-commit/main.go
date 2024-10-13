package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hambosto/ai-generate-commit/internal/config"
	"github.com/hambosto/ai-generate-commit/internal/git"
	"github.com/hambosto/ai-generate-commit/internal/service"
)

func main() {
	// Define command line flags for setting and getting configuration
	setConfigCmd := flag.NewFlagSet("setConfig", flag.ExitOnError)
	setConfigKey := setConfigCmd.String("key", "", "Config key")
	setConfigValue := setConfigCmd.String("value", "", "Config value")

	getConfigCmd := flag.NewFlagSet("getConfig", flag.ExitOnError)
	getConfigKey := getConfigCmd.String("key", "", "Config key")

	// Check for command line arguments
	if len(os.Args) < 2 {
		runGenerate() // Run generate by default if no command is given
		return
	}

	// Parse the command line arguments
	switch os.Args[1] {
	case "setConfig":
		setConfigCmd.Parse(os.Args[2:])
		config.SetConfig(*setConfigKey, *setConfigValue)
	case "getConfig":
		getConfigCmd.Parse(os.Args[2:])
		value := config.GetConfig(*getConfigKey)
		fmt.Printf("%s=%s\n", *getConfigKey, value)
	case "getConfigPath":
		fmt.Printf("Configuration file path: %s\n", config.GetConfigPath())
	case "generate":
		runGenerate() // Directly call runGenerate for the generate command
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

// runGenerate orchestrates the steps to generate and commit a commit message.
func runGenerate() {
	// Ensure the current directory is a git repository
	if err := git.AssertGitRepo(); err != nil {
		log.Fatal(err)
	}

	// Check if files are staged for commit
	if err := git.EnsureFilesAreStaged(); err != nil {
		log.Fatal(err)
	}

	// Get the list of staged files and their differences
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		log.Fatal(err)
	}

	diff, err := git.GetDiff(stagedFiles)
	if err != nil {
		log.Fatal(err)
	}

	if diff == "" {
		log.Fatal("No changes detected in the staged files. Please make some changes before generating a commit message.")
	}

	// Generate the commit message using the service
	commitMessage, err := service.GenerateCommitMessage(diff)
	if err != nil {
		log.Fatal(err)
	}

	// Display the generated commit message
	fmt.Printf("Generated Commit Message:\n\n%s\n\n", commitMessage)

	// Prompt the user for confirmation to use the commit message
	if confirmCommit(commitMessage) {
		if err := git.GitCommit(commitMessage); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Changes committed successfully.")
	} else {
		fmt.Println("Commit aborted.")
	}
}

// confirmCommit prompts the user to confirm if they want to use the generated commit message.
func confirmCommit(commitMessage string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to use this commit message? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	return response == "y" || response == "Y"
}
