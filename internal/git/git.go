package git

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// FileStatus represents a file's path and its current Git status.
type FileStatus struct {
	Path   string // Path of the file.
	Status string // Status of the file (e.g., Modified, Added, etc.).
}

// AssertGitRepo checks if the current directory is a Git repository.
// It returns an error if the directory is not inside a Git work tree.
func AssertGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return errors.New("not a Git repository")
	}
	return nil
}

// GetStagedFiles returns a slice of staged file names.
// These are the files that have been added to the index but not yet committed.
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return filterFiles(files), nil
}

// GetChangedFiles returns a slice of FileStatus for all the files with changes (staged or unstaged).
// It retrieves the status of each file using "git status --porcelain".
func GetChangedFiles() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting git status: %w", err)
	}

	var changedFiles []FileStatus
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue // Skip invalid lines that are too short.
		}
		status := strings.TrimSpace(line[:2])
		file := strings.TrimSpace(line[3:])
		changedFiles = append(changedFiles, FileStatus{
			Path:   file,
			Status: translateStatus(status),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning git status output: %w", err)
	}

	return changedFiles, nil
}

// translateStatus converts a Git status code (e.g., "M", "A") to a human-readable string.
func translateStatus(status string) string {
	switch status {
	case "M":
		return "Modified"
	case "A":
		return "Added"
	case "D":
		return "Deleted"
	case "R":
		return "Renamed"
	case "C":
		return "Copied"
	case "U":
		return "Updated but unmerged"
	case "??":
		return "Untracked"
	default:
		return "Unknown"
	}
}

// filterFiles removes empty file names from the list of files.
func filterFiles(files []string) []string {
	var filteredFiles []string
	for _, file := range files {
		if file != "" {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

// GetDiff returns the diff of the provided list of files.
// It runs "git diff --staged" to show the diff for the staged changes.
func GetDiff(files []string) (string, error) {
	args := append([]string{"diff", "--staged", "--"}, files...)
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GitCommit creates a new Git commit with the provided message.
// It runs "git commit -m" with the specified commit message.
func GitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// EnsureFilesAreStaged checks if there are any staged files.
// If no files are staged, it prompts the user to stage all changes.
// If there are no changes, it returns an error.
func EnsureFilesAreStaged() error {
	stagedFiles, err := GetStagedFiles()
	if err != nil {
		return err
	}
	if len(stagedFiles) == 0 {
		changedFiles, err := GetChangedFiles()
		if err != nil {
			return err
		}
		if len(changedFiles) > 0 {
			fmt.Println("The following files have changes:")
			for _, file := range changedFiles {
				fmt.Printf("%s: %s\n", file.Status, file.Path)
			}
			fmt.Print("Do you want to stage all these changes? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				cmd := exec.Command("git", "add", ".")
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("error staging files: %w", err)
				}
				fmt.Println("Changes staged successfully.")
			} else {
				return fmt.Errorf("no staged files. Please stage files before generating a commit message")
			}
		} else {
			return fmt.Errorf("no changes detected. Please make some changes before generating a commit message")
		}
	}
	return nil
}
