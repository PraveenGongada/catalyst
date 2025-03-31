/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func IsGHInstalled() error {
	_, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed or not in your PATH. " +
			"Please install it from https://cli.github.com/manual/installation " +
			"and run 'gh auth login' before using Catalyst")
	}
	return nil
}

func TriggerWorkflow(
	repository string,
	workflowID string,
	matrices []map[string]interface{},
	changeLog string,
	branchName string,
) error {
	payloadBytes, err := json.Marshal(map[string]interface{}{
		"matrices": matrices,
	})
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	args := []string{
		"workflow", "run",
		workflowID,
		"--repo", repository,
		"--ref", branchName,
	}

	// Add inputs as raw fields
	args = append(args, "--raw-field", fmt.Sprintf("payload=%s", string(payloadBytes)))
	args = append(args, "--raw-field", fmt.Sprintf("change_log=%s", changeLog))

	cmd := exec.Command("gh", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to trigger GitHub workflow: %w, output: %s", err, string(output))
	}

	return nil
}
