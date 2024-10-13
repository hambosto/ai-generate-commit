package service

import (
	"fmt"

	"github.com/hambosto/ai-generate-commit/internal/config"
	"github.com/hambosto/ai-generate-commit/internal/groq"
)

const (
	defaultModel  = "llama3-8b-8192"
	defaultPrompt = `
KEEP IN MIND THAT STICK TO THE POINT TO ONLY REPLY WITH MY PROMPTED MESSAGE!!! DO NOT ADD ANY ADDITIONAL INFORMATION !!!
DO NOT SAY "Here is the commit message" OR SUCH LIKE THAT. JUST REPLY ONLY THE COMMIT MESSAGE ITSELF !!!
You are an AI designed to generate concise and meaningful commit messages for code repositories, restricted to a single sentence. Craft your message based on the type of change, incorporating the appropriate prefix as follows:
  - [Add]: For new features, functions, or files.
  - [Fix]: For bug fixes or corrections.
  - [Update]: For updates or modifications to existing code.
  - [Remove]: For deletions of code or functionality.
  - [Chore]: For general tasks, maintenance, or minor changes.
  Example: [Update] (controllers/products.go, controllers/users.go) removed redundant BodyParser calls and directly used validated payload from Locals.
  Formatting Guidelines:
  1. If the combined length of the file names is 60 characters or fewer, format your message as follows:
  - '[Type] (file/s name separated by commas) $commit_message'
  2. If the combined length exceeds 60 characters, omit the file list:
  - '[Type] $commit_message'
  (do not include the prefix in the message).
  KEEP IN MIND THAT STICK TO THE POINT TO ONLY REPLY WITH MY PROMPTED MESSAGE!!! DO NOT ADD ANY ADDITIONAL INFORMATION !!!
  DO NOT SAY "Here is the commit message" OR SUCH LIKE THAT. JUST REPLY ONLY THE COMMIT MESSAGE ITSELF !!!
`
)

// CommitMessageGenerator handles the generation of commit messages.
type CommitMessageGenerator struct {
	client *groq.Client
	model  string
}

// NewCommitMessageGenerator creates a new CommitMessageGenerator.
func NewCommitMessageGenerator(model string) (*CommitMessageGenerator, error) {
	client, err := groq.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GROQ client: %w", err)
	}

	if model == "" {
		model = defaultModel
	}

	return &CommitMessageGenerator{
		client: client,
		model:  model,
	}, nil
}

// GenerateCommitMessage creates a commit message based on the provided diff.
func (g *CommitMessageGenerator) GenerateCommitMessage(diff string) (string, error) {
	commitPrompt, err := config.GetConfig("COMMIT_PROMPT")
	if err != nil {
		return "", fmt.Errorf("failed to get commit prompt: %w", err)
	}

	if commitPrompt == "" {
		commitPrompt = defaultPrompt
	}

	messages := []groq.Message{
		{Role: "system", Content: commitPrompt},
		{Role: "user", Content: fmt.Sprintf("Here's the git diff:\n%s", diff)},
	}

	return g.client.GenerateCompletion(messages, g.model)
}

