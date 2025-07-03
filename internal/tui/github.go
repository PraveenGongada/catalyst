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

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/types"
)

/*
func generateJSON(m *MainModel) tea.Cmd {
	return func() tea.Msg {
		confirmModel, ok := m.subModels[types.ConfirmStage].(*ConfirmModel)
		if !ok {
			return TriggerMsg{error: fmt.Errorf("could not get confirm model")}
		}

		confirmModel.matrixGenerator.SetSelectedApps(m.GetSelectedApps())
		confirmModel.matrixGenerator.SetSelectedPlatforms(m.GetSelectedPlatforms())
		confirmModel.matrixGenerator.SetSelectedEnvironments(m.GetSelectedEnvironments())

		generator := confirmModel.matrixGenerator

		inputsModel, ok := m.subModels[types.InputStage].(*InputsModel)
		var changeLog string
		var branchName string
		if ok {
			changeLog = inputsModel.changeLog
			branchName = inputsModel.GetBranchName()
		}

		purifiedMatrices := generator.GroupedMatricesPurified()

		totalMatrices := 0
		for _, matrices := range purifiedMatrices {
			totalMatrices += len(matrices)
		}

		if totalMatrices == 0 {
			return TriggerMsg{error: fmt.Errorf("no matrices generated from your selections")}
		}

		timestamp := time.Now().Format("20060102-150405")
		outputDir := fmt.Sprintf("catalyst-test-payloads-%s", timestamp)
		if err := os.Mkdir(outputDir, 0755); err != nil {
			return TriggerMsg{error: fmt.Errorf("failed to create output directory: %w", err)}
		}

		repository := m.config.GitHub.Repository
		for workflow, matrices := range purifiedMatrices {
			if len(matrices) == 0 {
				continue
			}

			workflowFile := ""
			if wf, ok := m.config.GitHub.Workflows[workflow]; ok {
				workflowFile = wf.File
			} else {
				workflowFile = workflow + ".yaml"
			}

			// Create file content structure
			fileContent := map[string]interface{}{
				"workflowFile": workflowFile,
				"branch":       branchName,
				"matrices":     matrices,
				"change_log":   changeLog,
			}

			payloadBytes, err := json.MarshalIndent(fileContent, "", "  ")
			if err != nil {
				return TriggerMsg{error: fmt.Errorf("error marshaling payload: %w", err)}
			}

			filename := filepath.Join(
				outputDir,
				fmt.Sprintf("%s-%s.json", workflow, strings.ReplaceAll(repository, "/", "-")),
			)
			if err := os.WriteFile(filename, payloadBytes, 0644); err != nil {
				return TriggerMsg{error: fmt.Errorf("failed to write payload file: %w", err)}
			}
		}

		return TriggerMsg{error: nil}
	}
}
*/

func triggerAction(m *MainModel) tea.Cmd {
	return func() tea.Msg {
		confirmModel, ok := m.subModels[types.ConfirmStage].(*ConfirmModel)
		if !ok {
			return TriggerMsg{error: fmt.Errorf("could not get confirm model")}
		}

		confirmModel.matrixGenerator.SetSelectedApps(m.GetSelectedApps())
		confirmModel.matrixGenerator.SetSelectedPlatforms(m.GetSelectedPlatforms())
		confirmModel.matrixGenerator.SetSelectedEnvironments(m.GetSelectedEnvironments())

		generator := confirmModel.matrixGenerator

		inputsModel, ok := m.subModels[types.InputStage].(*InputsModel)
		var changeLog string
		branchName := "main"
		if ok {
			changeLog = strings.TrimSpace(inputsModel.changeLog)
			branchName = inputsModel.GetBranchName()
		}

		purifiedMatrices := generator.GroupedMatricesPurified()

		totalMatrices := 0
		for _, matrices := range purifiedMatrices {
			totalMatrices += len(matrices)
		}

		if totalMatrices == 0 {
			return TriggerMsg{error: fmt.Errorf("no matrices generated from your selections")}
		}

		repository := m.config.GitHub.Repository

		var errors []string

		for workflow, matrices := range purifiedMatrices {
			if len(matrices) == 0 {
				continue
			}

			if wf, ok := m.config.GitHub.Workflows[workflow]; ok {
				err := github.TriggerWorkflow(
					repository,
					wf.File,
					matrices,
					changeLog,
					branchName,
				)
				if err != nil {
					errors = append(errors, fmt.Sprintf("'%s': %v", workflow, err))
				}
			} else {
				errors = append(errors, fmt.Sprintf("workflow '%s' not found in configuration", workflow))
			}
		}

		if len(errors) > 0 {
			return TriggerMsg{error: fmt.Errorf("workflow trigger errors: %s", strings.Join(errors, "; "))}
		}

		return TriggerMsg{error: nil}
	}
}
