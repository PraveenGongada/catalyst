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

package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultHeight = 44
	DefaultWidth  = 186
)

var AppStyle = lipgloss.NewStyle().Padding(1)

var TitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#5ea1ff"))

var ItemStyle = lipgloss.NewStyle().PaddingLeft(6)

var SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

var PaginationStyle = list.DefaultStyles().PaginationStyle

var HelpStyle = list.DefaultStyles().HelpStyle

var ChosenItemStyle = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#90EE90"))

var ChosenSelectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#90EE90"))

var SelectedDisplayStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#5ea1ff"))

var DotStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#808080"))

var GitHubMessageStyle = lipgloss.NewStyle().Bold(true).PaddingLeft(1).
	Foreground(lipgloss.Color("15"))

var GitHubErrorStyle = lipgloss.NewStyle().Bold(true).PaddingLeft(1).
	Foreground(lipgloss.Color("#D9534F"))

var CustomHelpStyle = lipgloss.NewStyle().Padding(0, 0, 1, 2).
	Foreground(lipgloss.Color("241"))

var SummaryHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Bold(true).
	Align(lipgloss.Center).Width(55)

var SummaryTitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Bold(true)

var SummaryValueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

var SummaryDividerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

var SummaryFooterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Align(lipgloss.Center).Width(55)
