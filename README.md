# AI Generate Commit

AI Generate Commit is a tool that uses artificial intelligence to automatically generate meaningful commit messages based on your staged changes. This tool aims to streamline the git commit process and ensure consistent, descriptive commit messages across your project.

## Features

- Automatically generates commit messages based on staged changes
- Uses the GROQ API for AI-powered commit message generation
- Configurable commit message prompt
- Easy-to-use command-line interface

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Installing from source

1. Clone the repository:

   ```
   git clone https://github.com/hambosto/ai-generate-commit.git
   ```

2. Change to the project directory:

   ```
   cd ai-generate-commit
   ```

3. Build the project:

   ```
   go build -o ai-generate-commit ./cmd/ai-generate-commit
   ```

4. (Optional) Move the binary to a directory in your PATH:
   ```
   sudo mv ai-generate-commit /usr/local/bin/ && sudo chmod +x /usr/local/bin/ai-generate-commit
   ```

### Installing using `go install`

1. Run following command to install the tool directly:

```
go install github.com/hambosto/ai-generate-commit/cmd/ai-generate-commit@latest
```

2. Ensure that `$GOPATH/bin` is in your `PATH` so that you can run the tool from anywhere:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

### Installing from release

1. Go to the [Releases](https://github.com/hambosto/ai-generate-commit/releases) page of the project.
2. Download the appropriate binary for your operating system.
3. (Optional) Move the binary to a directory in your PATH:
   ```
   sudo mv ai-generate-commit-* /usr/local/bin/
   ```

## Configuration

Before using the tool, you need to set up your GROQ API key and customize the commit prompt if desired.

### Setting up GROQ API Key

1. Sign up for a GROQ account and obtain your API key from [https://console.groq.com](https://console.groq.com).

2. Set your GROQ API key using the following command:
   ```
   ai-generate-commit setConfig -key GROQ_APIKEY -value your_api_key_here
   ```

### Customizing the Commit Prompt

You can customize the prompt used for generating commit messages:

```
ai-generate-commit setConfig -key COMMIT_PROMPT -value "Your custom prompt here. Use {diff} as a placeholder for the git diff."
```

Default prompt if not set:

```
Based on the following git diff, please generate a concise and informative commit message:

{diff}

The commit message should follow these guidelines:
1. Start with a short summary (50 chars or less)
2. Followed by a blank line
3. Followed by a more detailed description, if necessary

Please write the commit message now:
```

## Usage

1. Stage your changes using `git add`.

2. Run the tool:

   ```
   ai-generate-commit generate
   ```

3. Review the generated commit message and confirm if you want to use it.

## Additional Commands

- Get the current value of a configuration key:

  ```
  ai-generate-commit getConfig -key KEY_NAME
  ```

- Get the path of the configuration file:
  ```
  ai-generate-commit getConfigPath
  ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
