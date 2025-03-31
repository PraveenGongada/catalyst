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
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/types"
	"github.com/PraveenGongada/catalyst/internal/util"
)

func getEnvSelections(m *MainModel) []types.SelectableItem {
	environments := m.config.GetEnvironments(
		m.GetSelectedApps(),
		m.GetSelectedPlatforms(),
	)

	selectedMap := make(map[string]bool)
	for _, env := range m.GetSelectedEnvironments() {
		selectedMap[env] = true
	}

	items := make([]types.SelectableItem, len(environments))
	for i, env := range environments {
		items[i] = types.SelectableItem{
			Title:    env,
			Selected: selectedMap[env],
		}
	}

	return items
}

type EnvSelectModel struct {
	selections []types.SelectableItem
	list       list.Model
	mainModel  *MainModel
}

func NewEnvSelectModel(m *MainModel) *EnvSelectModel {
	selections := getEnvSelections(m)
	envSelectionList := envSelectModelList(m, &selections)
	return &EnvSelectModel{
		selections: selections,
		list:       envSelectionList,
		mainModel:  m,
	}
}

func (m *EnvSelectModel) Init() tea.Cmd {
	m.selections = getEnvSelections(m.mainModel)
	m.list = envSelectModelList(m.mainModel, &m.selections)
	return nil
}

func (m *EnvSelectModel) View() string {
	return styles.AppStyle.Render(m.list.View())
}

func (m *EnvSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		case " ":
			m.toggleSelected()
			m.mainModel.UpdateSelectionsFromModel(types.EnvSelectStage, m.selections)
		case "right", "enter", "ctrl+n":
			if util.HasSelection(m.selections) {
				m.mainModel.UpdateSelectionsFromModel(types.EnvSelectStage, m.selections)
				m.mainModel.moveToNextStage()
				return m.mainModel, nil
			}
		case "ctrl+p", "left":
			m.mainModel.UpdateSelectionsFromModel(types.EnvSelectStage, m.selections)
			m.mainModel.moveToPreviousStage()
			return m.mainModel, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *EnvSelectModel) toggleSelected() {
	if selected, ok := m.list.SelectedItem().(types.SelectableItem); ok {
		index := m.list.Index()
		m.selections[index].Selected = !m.selections[index].Selected
		selected.Selected = !selected.Selected
		m.list.SetItem(index, selected)
	}
}

func envSelectModelList(m *MainModel, selections *[]types.SelectableItem) list.Model {
	items := make([]list.Item, len(*selections))
	for i, item := range *selections {
		items[i] = item
	}
	l := list.New(items, types.ItemDelegate{}, m.width, m.height)
	l.Title = "Please select the environments to be included..."
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return util.AdditionalHelpKeys()
	}
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.TitleStyle
	l.Styles.PaginationStyle = styles.PaginationStyle
	l.Styles.HelpStyle = styles.HelpStyle

	return l
}
