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
	setConfigCmd := flag.NewFlagSet("setConfig", flag.ExitOnError)
	setConfigKey := setConfigCmd.String("key", "", "Config key")
	setConfigValue := setConfigCmd.String("value", "", "Config value")

	getConfigCmd := flag.NewFlagSet("getConfig", flag.ExitOnError)
	getConfigKey := getConfigCmd.String("key", "", "Config key")

	if len(os.Args) < 2 {
		runGenerate()
		return
	}

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
			runGenerate()
		default:
			fmt.Println("Unknown command")
			os.Exit(1)
		}
		runGenerate()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func runGenerate() {
	err := git.AssertGitRepo()
	if err != nil {
		log.Fatal(err)
	}

	err = git.EnsureFilesAreStaged()
	if err != nil {
		log.Fatal(err)
	}

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

	commitMessage, err := service.GenerateCommitMessage(diff)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated Commit Message:\n\n%s\n\n", commitMessage)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to use this commit message? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "y" || response == "Y" {
		err = git.GitCommit(commitMessage)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Changes committed successfully.")
	} else {
		fmt.Println("Commit aborted.")
	}
}
