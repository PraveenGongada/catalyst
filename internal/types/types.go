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

package types

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/styles"
)

type Stage int

const (
	AppSelectStage Stage = iota
	OSSelectStage
	EnvSelectStage
	InputStage
	ConfirmStage
	GitActionStage
)

func (s Stage) String() string {
	return [...]string{"AppSelect", "OSSelect", "EnvSelect", "InputStage", "ConfirmStage"}[s]
}

func (s Stage) OutPutString() string {
	return [...]string{"Selected Apps", "Selected Platforms", "Selected Environments", "", ""}[s]
}

type DeploymentConfig struct {
	AppName     string
	OS          string
	Environment string
	Version     string
}

type MatrixConfig struct {
	AppName     string `json:"appName"`
	OS          string `json:"os"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type SelectableItem struct {
	Title    string
	Selected bool
}

func (i SelectableItem) FilterValue() string { return i.Title }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int { return 1 }

func (d ItemDelegate) Spacing() int { return 0 }

func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(SelectableItem)
	if !ok {
		return
	}

	str := item.Title
	fn := styles.ItemStyle.Render
	isSelected := index == m.Index()

	if isSelected {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render(
				">" + styles.DotStyle.Render("•") + strings.Join(s, " "),
			)
		}
	}

	if isSelected && item.Selected {
		fn = func(s ...string) string {
			return styles.ChosenSelectedStyle.Render("> ✔ " + strings.Join(s, " "))
		}
	} else if item.Selected {
		fn = func(s ...string) string {
			return styles.ChosenItemStyle.Render("✔ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
