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
	Path   string
	Status string
}

// ErrNotGitRepo is returned when the current directory is not a Git repository.
var ErrNotGitRepo = errors.New("not a Git repository")

// AssertGitRepo checks if the current directory is a Git repository.
func AssertGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return ErrNotGitRepo
	}
	return nil
}

// GetStagedFiles returns a slice of staged file names.
func GetStagedFiles() ([]string, error) {
	output, err := execGitCommand("git", "diff", "--name-only", "--cached")
	if err != nil {
		return nil, err
	}
	return filterEmptyStrings(strings.Split(output, "\n")), nil
}

// GetChangedFiles returns a slice of FileStatus for all files with changes.
func GetChangedFiles() ([]FileStatus, error) {
	output, err := execGitCommand("git", "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("error getting git status: %w", err)
	}

	var changedFiles []FileStatus
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
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

// GetDiff returns the diff of the provided list of files.
func GetDiff(files []string) (string, error) {
	args := append([]string{"diff", "--staged", "--"}, files...)
	return execGitCommand("git", args...)
}

// GitCommit creates a new Git commit with the provided message.
func GitCommit(message string) error {
	_, err := execGitCommand("git", "commit", "-m", message)
	return err
}

// EnsureFilesAreStaged checks if there are any staged files and prompts to stage if necessary.
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
		if len(changedFiles) == 0 {
			return errors.New("no changes detected")
		}

		fmt.Println("The following files have changes:")
		for _, file := range changedFiles {
			fmt.Printf("%s: %s\n", file.Status, file.Path)
		}

		if !promptYesNo("Do you want to stage all these changes?") {
			return errors.New("no staged files")
		}

		if _, err := execGitCommand("git", "add", "."); err != nil {
			return fmt.Errorf("error staging files: %w", err)
		}
		fmt.Println("Changes staged successfully.")
	}
	return nil
}

// Helper functions

func execGitCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func filterEmptyStrings(slice []string) []string {
	var filtered []string
	for _, s := range slice {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func translateStatus(status string) string {
	statusMap := map[string]string{
		"M":  "Modified",
		"A":  "Added",
		"D":  "Deleted",
		"R":  "Renamed",
		"C":  "Copied",
		"U":  "Updated but unmerged",
		"??": "Untracked",
	}
	if translated, ok := statusMap[status]; ok {
		return translated
	}
	return "Unknown"
}

func promptYesNo(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

