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

func getOSSelections(m *MainModel) []types.SelectableItem {
	platforms := m.config.GetPlatforms(m.GetSelectedApps())

	selectedMap := make(map[string]bool)
	for _, platform := range m.GetSelectedPlatforms() {
		selectedMap[platform] = true
	}

	items := make([]types.SelectableItem, len(platforms))
	for i, platform := range platforms {
		items[i] = types.SelectableItem{
			Title:    platform,
			Selected: selectedMap[platform],
		}
	}

	return items
}

type OSSelectModel struct {
	selections []types.SelectableItem
	list       list.Model
	mainModel  *MainModel
}

func NewOSSelectModel(m *MainModel) *OSSelectModel {
	selections := getOSSelections(m)
	osSelectionList := createOSSelectModelList(m, &selections)
	return &OSSelectModel{
		selections: selections,
		list:       osSelectionList,
		mainModel:  m,
	}
}

func (m *OSSelectModel) Init() tea.Cmd {
	m.selections = getOSSelections(m.mainModel)
	m.list = createOSSelectModelList(m.mainModel, &m.selections)
	return nil
}

func (m *OSSelectModel) View() string {
	return styles.AppStyle.Render(m.list.View())
}

func (m *OSSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		case " ":
			m.toggleSelected()
			m.mainModel.UpdateSelectionsFromModel(types.OSSelectStage, m.selections)
		case "right", "enter", "ctrl+n":
			if util.HasSelection(m.selections) {
				m.mainModel.UpdateSelectionsFromModel(types.OSSelectStage, m.selections)
				m.mainModel.moveToNextStage()
				return m.mainModel, nil
			}
		case "ctrl+p", "left":
			m.mainModel.UpdateSelectionsFromModel(types.OSSelectStage, m.selections)
			m.mainModel.moveToPreviousStage()
			return m.mainModel, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *OSSelectModel) toggleSelected() {
	if selected, ok := m.list.SelectedItem().(types.SelectableItem); ok {
		index := m.list.Index()
		m.selections[index].Selected = !m.selections[index].Selected
		selected.Selected = !selected.Selected
		m.list.SetItem(index, selected)
	}
}

func createOSSelectModelList(m *MainModel, selections *[]types.SelectableItem) list.Model {
	items := make([]list.Item, len(*selections))
	for i, item := range *selections {
		items[i] = item
	}
	l := list.New(items, types.ItemDelegate{}, m.width, m.height)
	l.Title = "Please select the platforms to be included..."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return util.AdditionalHelpKeys()
	}
	l.Styles.Title = styles.TitleStyle
	l.Styles.PaginationStyle = styles.PaginationStyle
	l.Styles.HelpStyle = styles.HelpStyle

	return l
}
