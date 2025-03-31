# Catalyst Examples

This directory contains example configurations for Catalyst to help you get started.

## Example Configuration

`catalyst-example.yaml` demonstrates a basic configuration for Catalyst, including:

- GitHub repository and workflow definitions
- Input configuration with variable substitution
- Matrix configurations for multiple apps, platforms, and environments

### How to Use the Example

1. Copy the example configuration to your project:

```bash
cp catalyst-example.yaml catalyst.yaml
```

2. Edit the `catalyst.yaml` file to match your GitHub repository, workflows, and apps.

3. Run Catalyst:

```bash
catalyst
```

## Configuration Structure

The configuration file has three main sections:

### 1. GitHub Configuration

```yaml
github:
  repository: "your-org/sample-repo" # Your GitHub repository
  workflows: # Define your GitHub Actions workflows
    workflow_key: # A unique key for each workflow
      name: "Readable Workflow Name" # Human-readable name
      file: "workflow_file.yml" # The workflow file in your repository
```

### 2. Inputs Configuration

```yaml
inputs:
  input_key: # A unique key for the input
    description: "Description text" # Helps users understand the input
    required: true/false # Whether the input is required
    default: "Default value" # Default value if not provided
```

### 3. Matrix Configuration

```yaml
matrix:
  AppName: # Your app's name
    platform: # Platform (e.g., ios, android)
      Environment: # Environment (e.g., Development, Production)
        workflow: "workflow_key" # The workflow to trigger (from github.workflows)
        matrix: # Matrix parameters passed to GitHub Actions
          param1: "value1"
          param2: "{{inputs.input_key}}" # Reference to an input value
```

## Variable Substitution

You can reference input values in your matrix configurations using the `{{inputs.input_key}}` syntax. When Catalyst runs, these will be replaced with the actual values provided by the user.
