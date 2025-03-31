# 🚀 Catalyst

<div align="center">
  
[![Release](https://img.shields.io/github/v/release/PraveenGongada/catalyst?style=flat-square)](https://github.com/PraveenGongada/catalyst/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/PraveenGongada/catalyst)](https://goreportcard.com/report/github.com/PraveenGongada/catalyst)
[![Go Version](https://img.shields.io/github/go-mod/go-version/PraveenGongada/catalyst?style=flat-square)](go.mod)
[![License](https://img.shields.io/github/license/PraveenGongada/catalyst?style=flat-square)](LICENSE)

</div>

Catalyst is an elegant terminal UI tool that simplifies triggering GitHub Actions workflows with matrix configurations. It's designed to streamline the deployment process for mobile applications across multiple platforms and environments.

![Catalyst Demo](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/catalyst-demo.gif)

## ✨ Features

- 📱 Select multiple apps for deployment
- 📲 Choose target platforms (iOS, Android)
- 🌿 Select environments (Debug, Production)
- 🔧 Configure dynamic input values
- 📝 Add changelog information for deployments
- 🔀 Specify target branch for workflow execution
- 🔍 Preview matrix configurations before triggering
- 🔄 Trigger GitHub Actions workflows with complex matrix configurations

## 📦 Installation

### Prerequisites

- [GitHub CLI](https://cli.github.com/) (gh) - Required for triggering workflows
  - Install from: https://cli.github.com/manual/installation
  - Authenticate with: `gh auth login`

### Using Homebrew (macOS)

```bash
brew tap PraveenGongada/tap
brew install catalyst
```

### Using Go

```bash
go install github.com/PraveenGongada/catalyst/cmd/catalyst@latest
```

### Manual Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/PraveenGongada/catalyst/releases).

#### Linux / macOS

```bash
# Download the latest release (replace X.Y.Z with actual version)
curl -L https://github.com/PraveenGongada/catalyst/releases/download/vX.Y.Z/catalyst_Linux_x86_64.tar.gz -o catalyst.tar.gz

# Extract the binary
tar -xzf catalyst.tar.gz

# Move to a directory in your PATH
sudo mv catalyst /usr/local/bin/
```

## 🔧 Configuration

Catalyst uses a YAML configuration file to define your matrix configurations. By default, it looks for `catalyst.yaml` in the current directory.

### Configuration File Location

Catalyst looks for the configuration file in the following order:

1. Path specified with the `-config` flag
2. Path set in the `CATALYST_CONFIG` environment variable
3. `./catalyst.yaml` in the current directory
4. `$USER_CONFIG_DIR/catalyst/catalyst.yaml` (e.g., `~/.config/catalyst/catalyst.yaml` on Linux)

### Adding to Your Shell Profile

Add the following to your `.zshrc`, `.bashrc`, or equivalent:

```bash
# Set the path to your catalyst config
export CATALYST_CONFIG="/path/to/your/catalyst.yaml"
```

### Sample Configuration

Here's a simplified example of a configuration file:

```yaml
# GitHub workflow metadata
github:
  repository: "your-org/mobile-apps"
  workflows:
    ios_debug:
      name: "iOS Debug Workflow"
      file: "ios_debug.yml"
    ios_prod:
      name: "iOS Production Workflow"
      file: "ios_prod.yml"

# Dynamic inputs that can be referenced in matrices
inputs:
  version:
    description: "App version"
    required: true
    default: "1.0.0"

# Matrix configurations
matrix:
  MyApp:
    ios:
      Debug:
        workflow: "ios_debug"
        matrix:
          bundle_id: "com.example.myapp.debug"
          version: "{{inputs.version}}"
      Prod:
        workflow: "ios_prod"
        matrix:
          bundle_id: "com.example.myapp"
          version: "{{inputs.version}}"
```

## 🚀 Usage

Run Catalyst from your terminal:

```bash
catalyst
```

Or specify a custom configuration path:

```bash
catalyst -config /path/to/your/catalyst.yaml
```

Check the version:

```bash
catalyst -version
```

## 📋 GitHub Actions Workflow Setup

To use Catalyst with your GitHub Actions, you need to set up a repository dispatch event in your workflow:

```yaml
# .github/workflows/example.yml
name: Example Deployment Workflow

on:
  workflow_dispatch:
    inputs:
      payload:
        description: "JSON payload containing matrix configurations"
        required: true
      change_log:
        description: "Changelog for this deployment"
        required: true

jobs:
  deploy:
    runs-on: ubuntu-latest
    # Parse the matrices JSON from the payload
    strategy:
      matrix:
        include: ${{ fromJson(fromJson(inputs.payload).matrices) }}

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      # Add your deployment steps here
```

## 📸 Interface Screenshots

Here's a visual walkthrough of the Catalyst interface:

### App Selection

![App Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/app-selection.png)

### Platform Selection

![Platform Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/platform-selection.png)

### Environment Selection

![Environment Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/env-selection.png)

### Input Configuration

![Input Configuration Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/input-configuration.png)

### Deployment Summary

![Deployment Summary Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/deployment-summary.png)

### Matrix Preview

![Matrix Preview Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/matrix-preview.png)

## 👨‍💻 Development

### Prerequisites

- Go 1.22 or higher
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/praveengongada/catalyst.git
cd catalyst

# Build the binary (main is in cmd/catalyst)
go build -o catalyst ./cmd/catalyst

# Run the application
./catalyst
```

### Making a Release

Catalyst uses GoReleaser to automate the release process:

1. Update the version in your code
2. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. The GitHub Actions workflow will automatically build and publish the release

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Third-party library attributions can be found in the [NOTICE](NOTICE) file.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
