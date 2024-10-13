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
	// Main entry point of the application. It calls the run() function
	// and handles any errors by logging them and terminating the program.
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Determines which command to execute based on the provided arguments.
	// Defaults to running the "generate" command if no arguments are given.
	if len(os.Args) < 2 {
		return runGenerate()
	}

	// Switches between different commands based on the first argument.
	switch os.Args[1] {
	case "setConfig":
		return runSetConfig()
	case "getConfig":
		return runGetConfig()
	case "getConfigPath":
		return runGetConfigPath()
	case "generate":
		return runGenerate()
	default:
		// Returns an error if an unknown command is provided.
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func runSetConfig() error {
	// Defines the "setConfig" command to set a configuration key-value pair.
	cmd := flag.NewFlagSet("setConfig", flag.ExitOnError)
	key := cmd.String("key", "", "Config key")
	value := cmd.String("value", "", "Config value")

	// Parses the arguments for the setConfig command.
	if err := cmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	// Ensures both key and value are provided, returns an error otherwise.
	if *key == "" || *value == "" {
		return fmt.Errorf("both key and value must be provided")
	}

	// Calls SetConfig from the config package to save the key-value pair.
	return config.SetConfig(*key, *value)
}

func runGetConfig() error {
	// Defines the "getConfig" command to retrieve a configuration value by key.
	cmd := flag.NewFlagSet("getConfig", flag.ExitOnError)
	key := cmd.String("key", "", "Config key")

	// Parses the arguments for the getConfig command.
	if err := cmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	// Ensures the key is provided, returns an error otherwise.
	if *key == "" {
		return fmt.Errorf("key must be provided")
	}

	// Retrieves the configuration value for the given key.
	value, err := config.GetConfig(*key)
	if err != nil {
		return err
	}

	// Prints the retrieved key-value pair.
	fmt.Printf("%s=%s\n", *key, value)
	return nil
}

func runGetConfigPath() error {
	// Prints the path to the configuration file.
	fmt.Printf("Configuration file path: %s\n", config.GetConfigPath())
	return nil
}

func runGenerate() error {
	// Ensures that the current directory is a valid Git repository.
	if err := git.AssertGitRepo(); err != nil {
		return err
	}

	// Checks if there are files staged for commit.
	if err := git.EnsureFilesAreStaged(); err != nil {
		return err
	}

	// Retrieves a list of staged files.
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return err
	}

	// Gets the diff (changes) for the staged files.
	diff, err := git.GetDiff(stagedFiles)
	if err != nil {
		return err
	}

	// Returns an error if no changes are detected in the staged files.
	if diff == "" {
		return fmt.Errorf("no changes detected in the staged files")
	}

	// Initializes the commit message generator.
	generator, err := service.NewCommitMessageGenerator("")
	if err != nil {
		return err
	}

	// Generates the commit message based on the diff.
	commitMessage, err := generator.GenerateCommitMessage(diff)
	if err != nil {
		return err
	}

	// Displays the generated commit message.
	fmt.Printf("Generated Commit Message:\n\n%s\n\n", commitMessage)

	// Prompts the user for confirmation to proceed with the commit.
	if confirmCommit() {
		// Commits the changes with the generated commit message if confirmed.
		if err := git.GitCommit(commitMessage); err != nil {
			return err
		}
		fmt.Println("Changes committed successfully.")
	} else {
		// Aborts the commit if the user declines.
		fmt.Println("Commit aborted.")
	}

	return nil
}

func confirmCommit() bool {
	// Prompts the user to confirm if they want to use the generated commit message.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to use this commit message? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}
		// Converts the response to lowercase and removes surrounding whitespace.
		response = strings.TrimSpace(strings.ToLower(response))
		// Checks for valid inputs (y/n) and returns a boolean value accordingly.
		switch response {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Invalid input. Please enter 'y' for yes or 'n' for no.")
		}
	}
}
