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
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/matrix"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/types"
)

type TriggerMsg struct {
	error error
}

type ConfirmModel struct {
	mainModel       *MainModel
	viewport        viewport.Model
	spinner         spinner.Model
	help            help.Model
	keys            KeyMap
	triggered       bool
	isLoading       bool
	error           error
	showPreview     bool
	matrixGenerator *matrix.Generator
	summaryContent  string
}

func NewConfirmModel(m *MainModel) *ConfirmModel {
	vp := viewport.New(m.width, m.height-1)
	vp.HighPerformanceRendering = false

	s := spinner.New()
	s.Spinner = spinner.Meter

	h := help.New()

	generator := matrix.NewGenerator(m.config)

	generator.SetSelectedApps(m.GetSelectedApps())
	generator.SetSelectedPlatforms(m.GetSelectedPlatforms())
	generator.SetSelectedEnvironments(m.GetSelectedEnvironments())

	inputsModel, ok := m.subModels[types.InputStage].(*InputsModel)
	if ok {
		for key, value := range inputsModel.GetInputValues() {
			generator.SetInputValue(key, value)
		}
	}

	model := &ConfirmModel{
		mainModel:       m,
		viewport:        vp,
		spinner:         s,
		help:            h,
		keys:            DefaultKeyMap,
		triggered:       false,
		showPreview:     false,
		matrixGenerator: generator,
	}
	return model
}

func (m *ConfirmModel) Init() tea.Cmd {
	m.summaryContent = m.DeploymentSummary()
	m.viewport.SetContent(m.summaryContent)
	return nil
}

func (m *ConfirmModel) View() string {
	viewportContent := m.viewport.View()

	if m.showPreview {
		previewHelp := []string{
			"↑/k: scroll up",
			"↓/j: scroll down",
			"esc: return to summary",
		}

		helpText := strings.Join(previewHelp, " • ")

		return viewportContent + "\n\n" + styles.CustomHelpStyle.Render(helpText)
	}

	helpParts := strings.Split(m.help.View(m.keys), "  ")
	helpView := strings.Join(helpParts, " • ")

	if m.isLoading {
		return viewportContent + "\n\n" + styles.GitHubMessageStyle.Render(
			fmt.Sprintf("Triggering GitHub Action Workflows %s", m.spinner.View()),
		)
	} else if m.error != nil {
		return viewportContent + "\n\n" + styles.GitHubErrorStyle.Render("Error Triggering GitHub Action: "+m.error.Error())
	} else if m.triggered {
		return viewportContent + "\n\n" + styles.GitHubMessageStyle.Render("GitHub Action Workflows Triggered Successfully!")
	}
	return viewportContent + "\n\n" + styles.CustomHelpStyle.Render(helpView)
}

func (m *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4

	case tea.KeyMsg:
		if !m.isLoading {
			if m.showPreview {
				switch msg.String() {
				case "up", "k", "down", "j", "pgup", "pgdown", "home", "end":
				case "esc", "q":
					m.showPreview = false
					m.viewport.SetContent(m.summaryContent)
					m.viewport.GotoTop()
					return m, nil
				default:
					return m, nil
				}
			}

			switch msg.String() {
			case "ctrl+c":
				return m, tea.Interrupt

			case "esc", "q", "n", "N":
				return m, tea.Quit

			case "right", "enter", "ctrl+n", "y", "Y":
				if m.triggered {
					return m, nil
				}
				m.isLoading = true
				return m, tea.Batch(m.spinner.Tick, triggerAction(m.mainModel))

			case "ctrl+p", "left":
				m.mainModel.moveToPreviousStage()
				return m.mainModel, nil

			case "p", "P":
				if !m.triggered {
					m.showPreview = true

					previewContent := m.generateAllMatricesPreview()

					if previewContent == "" {
						previewContent = "No matrices to display based on your selections."
					}

					divider := styles.SummaryDividerStyle.Render(strings.Repeat("━", 55))

					content := lipgloss.JoinVertical(
						lipgloss.Left,
						divider,
						"",
						styles.SummaryHeaderStyle.Render("🔍 MATRIX PREVIEW: All Workflows"),
						"",
						divider,
						"",
						previewContent,
					)

					m.viewport.SetContent(content)
					m.viewport.GotoTop()
					return m, nil
				}
			}
		}
	case TriggerMsg:
		m.isLoading = false
		if msg.error != nil {
			m.error = msg.error
			m.triggered = false
			return m, nil
		}
		m.triggered = true
	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	var viewportCmd tea.Cmd
	m.viewport, viewportCmd = m.viewport.Update(msg)
	cmds = append(cmds, viewportCmd)

	return m, tea.Batch(cmds...)
}

func (m *ConfirmModel) generateAllMatricesPreview() string {
	groupedMatrices := m.matrixGenerator.GroupedMatricesWithMetadata()

	if len(groupedMatrices) == 0 {
		return "No matrices to display."
	}

	var preview strings.Builder

	var workflowNames []string
	for workflow := range groupedMatrices {
		workflowNames = append(workflowNames, workflow)
	}
	sort.Strings(workflowNames)

	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#90EE90"))
	workflowStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5ea1ff"))

	for _, workflow := range workflowNames {
		matrixList := groupedMatrices[workflow]
		if len(matrixList) == 0 {
			continue
		}

		workflowName := workflow
		if wf, ok := m.mainModel.config.GitHub.Workflows[workflow]; ok {
			workflowName = wf.Name
		}

		preview.WriteString(workflowStyle.Render(fmt.Sprintf("Workflow: %s", workflowName)))
		preview.WriteString("\n")
		preview.WriteString(
			styles.SummaryDividerStyle.Render(strings.Repeat("━", len(workflowName)+10)),
		)
		preview.WriteString("\n")

		for i, matrix := range matrixList {
			preview.WriteString(fmt.Sprintf("Matrix #%d:\n", i+1))

			if app, ok := matrix["app"]; ok {
				preview.WriteString(fmt.Sprintf("  • App: %s\n",
					valueStyle.Render(fmt.Sprintf("%v", app))))
			}

			if platform, ok := matrix["platform"]; ok {
				preview.WriteString(fmt.Sprintf("  • Platform: %s\n",
					valueStyle.Render(fmt.Sprintf("%v", platform))))
			}

			if env, ok := matrix["environment"]; ok {
				preview.WriteString(fmt.Sprintf("  • Environment: %s\n",
					valueStyle.Render(fmt.Sprintf("%v", env))))
			}

			preview.WriteString("  • Parameters:\n")

			var keys []string
			for k := range matrix {
				if k != "app" && k != "platform" && k != "environment" {
					keys = append(keys, k)
				}
			}

			sort.Strings(keys)

			for _, k := range keys {
				preview.WriteString(fmt.Sprintf("    - %s: %s\n",
					k, valueStyle.Render(fmt.Sprintf("%v", matrix[k]))))
			}

			preview.WriteString("\n")
		}

		preview.WriteString("\n")
	}

	return preview.String()
}

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Next    key.Binding
	Back    key.Binding
	Quit    key.Binding
	Preview key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Next, k.Back, k.Preview, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Next, k.Back, k.Preview, k.Quit},
	}
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Next: key.NewBinding(
		key.WithKeys("y", "enter", "right"),
		key.WithHelp("→/enter/y", "proceed"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+p", "left"),
		key.WithHelp("←/ctrl+p", "go back"),
	),
	Preview: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview matrices"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "n"),
		key.WithHelp("q/n", "quit"),
	),
}

func (m *ConfirmModel) DeploymentSummary() string {
	mainModel := m.mainModel

	m.matrixGenerator.SetSelectedApps(mainModel.GetSelectedApps())
	m.matrixGenerator.SetSelectedPlatforms(mainModel.GetSelectedPlatforms())
	m.matrixGenerator.SetSelectedEnvironments(mainModel.GetSelectedEnvironments())

	inputsModel, ok := mainModel.subModels[types.InputStage].(*InputsModel)
	if ok {
		for key, value := range inputsModel.GetInputValues() {
			m.matrixGenerator.SetInputValue(key, value)
		}
	}

	groupedMatrices := m.matrixGenerator.GroupedMatricesWithMetadata()
	totalCombinations := m.matrixGenerator.GetTotalCombinations()

	divider := styles.SummaryDividerStyle.Render(strings.Repeat("━", 55))

	header := styles.SummaryHeaderStyle.Render("🚀 GitHub Deployment Summary")
	space := ""

	apps := fmt.Sprintf(
		"%s\n   • %s",
		styles.SummaryTitleStyle.Render("📱 Apps Selected"),
		styles.SummaryValueStyle.Render(strings.Join(mainModel.GetSelectedApps(), "\n   • ")),
	)

	envs := fmt.Sprintf(
		"%s\n   • %s",
		styles.SummaryTitleStyle.Render("📂 Environments"),
		styles.SummaryValueStyle.Render(
			strings.Join(mainModel.GetSelectedEnvironments(), "\n   • "),
		),
	)

	platforms := fmt.Sprintf(
		"%s\n   • %s",
		styles.SummaryTitleStyle.Render("📲 Platforms"),
		styles.SummaryValueStyle.Render(strings.Join(mainModel.GetSelectedPlatforms(), "\n   • ")),
	)

	var inputsText strings.Builder
	inputsText.WriteString(styles.SummaryTitleStyle.Render("🔧 Input Values"))

	if inputsModel != nil {
		relevantInputs := getRelevantInputs(mainModel)

		for _, key := range relevantInputs {
			if value, ok := inputsModel.GetInputValues()[key]; ok {
				inputsText.WriteString(fmt.Sprintf("\n   • %s: %s",
					styles.SummaryTitleStyle.Render(key),
					styles.SummaryValueStyle.Render(value),
				))
			}
		}
	}

	var workflowsText strings.Builder
	workflowsText.WriteString(styles.SummaryTitleStyle.Render("🔄 Workflows To Trigger"))

	if totalCombinations == 0 {
		workflowsText.WriteString("\n   • No workflows will be triggered based on your selections")
	} else {
		for workflow, matrices := range groupedMatrices {
			if len(matrices) == 0 {
				continue
			}

			workflowName := workflow
			if wf, ok := mainModel.config.GitHub.Workflows[workflow]; ok {
				workflowName = wf.Name
			}

			workflowsText.WriteString(fmt.Sprintf("\n   • %s (%d matrix combinations)",
				styles.SummaryTitleStyle.Render(workflowName),
				len(matrices),
			))
		}
	}

	var changelogText string
	if inputsModel != nil {
		changelogText = fmt.Sprintf(
			"%s\n   • %s",
			styles.SummaryTitleStyle.Render("📝 Changelog"),
			styles.SummaryValueStyle.Render(
				strings.Join(strings.Split(inputsModel.changeLog, "\n"), "\n   • "),
			),
		)
	}

	var branchText string
	if inputsModel != nil {
		branchText = fmt.Sprintf(
			"%s\n   • %s",
			styles.SummaryTitleStyle.Render("🔀 Target Branch"),
			styles.SummaryValueStyle.Render(inputsModel.GetBranchName()),
		)
	}

	footer := styles.SummaryFooterStyle.Render(
		fmt.Sprintf(
			"%s\n%s",
			fmt.Sprintf(
				"Would you like to proceed with triggering %d matrix combinations?\n",
				totalCombinations,
			),
			"[Y] Yes    [N] No",
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		divider,
		space,
		header,
		space,
		divider,
		space,
		apps,
		space,
		platforms,
		space,
		envs,
		space,
		inputsText.String(),
		space,
		workflowsText.String(),
		space,
		changelogText,
		space,
		branchText,
		space,
		divider,
		space,
		footer,
		space,
	)
}
