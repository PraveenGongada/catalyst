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
	"errors"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/PraveenGongada/catalyst/internal/constants"
	"github.com/PraveenGongada/catalyst/internal/styles"
)

type InputsModel struct {
	mainModel       *MainModel
	form            *huh.Form
	inputs          map[string]string
	changeLog       string
	branchName      string
	tempInputValues []tempInput
}

type tempInput struct {
	key   string
	value string
}

func (m *InputsModel) Init() tea.Cmd {
	m.initTempValues()
	m.form = createInputForms(m)
	return m.form.Init()
}

func (m *InputsModel) initTempValues() {
	relevantInputs := getRelevantInputs(m.mainModel)

	m.tempInputValues = make([]tempInput, 0, len(relevantInputs))
	for _, key := range relevantInputs {
		value := m.inputs[key]
		if value == "" && m.mainModel.config.Inputs[key].Default != "" {
			value = m.mainModel.config.Inputs[key].Default
			m.inputs[key] = value
		}
		m.tempInputValues = append(m.tempInputValues, tempInput{
			key:   key,
			value: value,
		})
	}
}

func (m *InputsModel) updateInputsFromTemp() {
	for _, temp := range m.tempInputValues {
		m.inputs[temp.key] = temp.value
	}
}

func (m *InputsModel) View() string {
	return styles.AppStyle.Render(m.form.View())
}

func (m *InputsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc":
			return m, tea.Quit
		case "ctrl+n":
			if m.hasValidInputs() {
				m.updateInputsFromTemp()
				m.mainModel.moveToNextStage()
				return m.mainModel, nil
			}
		case "ctrl+p":
			m.mainModel.moveToPreviousStage()
			return m.mainModel, nil
		}
	}

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		m.updateInputsFromTemp()
		m.mainModel.moveToNextStage()
		return m.mainModel, nil
	}

	return m, tea.Batch(cmds...)
}

func NewInputsModel(m *MainModel) *InputsModel {
	inputs := &InputsModel{
		mainModel:  m,
		inputs:     make(map[string]string),
		branchName: "main",
	}

	for key, input := range m.config.Inputs {
		inputs.inputs[key] = input.Default
	}

	inputs.initTempValues()

	return inputs
}

var inputRefPattern = regexp.MustCompile(constants.RegexInputPlaceholder)

func getRelevantInputs(m *MainModel) []string {
	inputsNeeded := make(map[string]bool)

	selectedApps := m.GetSelectedApps()
	selectedPlatforms := m.GetSelectedPlatforms()
	selectedEnvironments := m.GetSelectedEnvironments()

	for _, app := range selectedApps {
		if appConfig, ok := m.config.Matrix[app]; ok {
			for _, platform := range selectedPlatforms {
				if platformConfig, ok := appConfig[platform]; ok {
					for _, env := range selectedEnvironments {
						if envConfig, ok := platformConfig[env]; ok {
							for _, value := range envConfig.Matrix {
								if strValue, ok := value.(string); ok {
									matches := inputRefPattern.FindAllStringSubmatch(strValue, -1)
									for _, match := range matches {
										if len(match) >= 2 {
											inputName := match[1]
											inputsNeeded[inputName] = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	result := []string{}

	for key := range m.config.Inputs {
		if inputsNeeded[key] {
			result = append(result, key)
		}
	}

	return result
}

func createInputForms(inputs *InputsModel) *huh.Form {
	m := inputs.mainModel

	inputFields := []huh.Field{}

	for i, temp := range inputs.tempInputValues {
		key := temp.key
		inputConfig := m.config.Inputs[key]

		field := huh.NewInput().
			Title(key + ": ").
			Description(inputConfig.Description).
			Value(&inputs.tempInputValues[i].value).
			Validate(func(s string) error {
				if inputConfig.Required && strings.TrimSpace(s) == "" {
					return errors.New(key + " is required")
				}
				return nil
			})

		inputFields = append(inputFields, field)
	}

	branchField := huh.NewInput().
		Value(&inputs.branchName).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("branch name is required")
			}
			return nil
		}).
		Title("Branch Name: ").
		Description("Please enter the branch name to trigger workflows on").
		Placeholder("main")

	inputFields = append(inputFields, branchField)

	changelogField := huh.NewText().
		Value(&inputs.changeLog).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("changelog is required")
			}
			return nil
		}).
		Title("Changelog: ").
		Description("Please add the changes in this release").
		CharLimit(400).
		Lines(5)

	inputFields = append(inputFields, changelogField)

	form := huh.NewForm(
		huh.NewGroup(inputFields...),
	).WithHeight(m.height - 1).WithWidth(m.width - 1)

	return form
}

func (m *InputsModel) hasValidInputs() bool {
	m.updateInputsFromTemp()

	for key, value := range m.inputs {
		inputConfig := m.mainModel.config.Inputs[key]
		if inputConfig.Required && strings.TrimSpace(value) == "" {
			return false
		}
	}

	return strings.TrimSpace(m.changeLog) != "" && strings.TrimSpace(m.branchName) != ""
}

func (m *InputsModel) GetInputValues() map[string]string {
	trimmedInputs := make(map[string]string)
	for key, value := range m.inputs {
		trimmedInputs[key] = strings.TrimSpace(value)
	}
	return trimmedInputs
}

func (m *InputsModel) GetBranchName() string {
	return strings.TrimSpace(m.branchName)
}
