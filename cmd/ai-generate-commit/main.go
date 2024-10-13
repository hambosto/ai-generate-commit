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
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return runGenerate()
	}

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
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func runSetConfig() error {
	cmd := flag.NewFlagSet("setConfig", flag.ExitOnError)
	key := cmd.String("key", "", "Config key")
	value := cmd.String("value", "", "Config value")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *key == "" || *value == "" {
		return fmt.Errorf("both key and value must be provided")
	}

	return config.SetConfig(*key, *value)
}

func runGetConfig() error {
	cmd := flag.NewFlagSet("getConfig", flag.ExitOnError)
	key := cmd.String("key", "", "Config key")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *key == "" {
		return fmt.Errorf("key must be provided")
	}

	value, err := config.GetConfig(*key)
	if err != nil {
		return err
	}

	fmt.Printf("%s=%s\n", *key, value)
	return nil
}

func runGetConfigPath() error {
	fmt.Printf("Configuration file path: %s\n", config.GetConfigPath())
	return nil
}

func runGenerate() error {
	if err := git.AssertGitRepo(); err != nil {
		return err
	}

	if err := git.EnsureFilesAreStaged(); err != nil {
		return err
	}

	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return err
	}

	diff, err := git.GetDiff(stagedFiles)
	if err != nil {
		return err
	}

	if diff == "" {
		return fmt.Errorf("no changes detected in the staged files")
	}

	generator, err := service.NewCommitMessageGenerator("")
	if err != nil {
		return err
	}

	commitMessage, err := generator.GenerateCommitMessage(diff)
	if err != nil {
		return err
	}

	fmt.Printf("Generated Commit Message:\n\n%s\n\n", commitMessage)

	if confirmCommit() {
		if err := git.GitCommit(commitMessage); err != nil {
			return err
		}
		fmt.Println("Changes committed successfully.")
	} else {
		fmt.Println("Commit aborted.")
	}

	return nil
}

func confirmCommit() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to use this commit message? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}
		response = strings.TrimSpace(strings.ToLower(response))
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
