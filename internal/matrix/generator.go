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

package matrix

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PraveenGongada/catalyst/internal/config"
)

type Generator struct {
	Config               *config.Config
	SelectedApps         []string
	SelectedPlatforms    []string
	SelectedEnvironments []string
	InputValues          map[string]string
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		Config:      cfg,
		InputValues: make(map[string]string),
	}
}

func (g *Generator) SetSelectedApps(apps []string) {
	g.SelectedApps = apps
}

func (g *Generator) SetSelectedPlatforms(platforms []string) {
	g.SelectedPlatforms = platforms
}

func (g *Generator) SetSelectedEnvironments(envs []string) {
	g.SelectedEnvironments = envs
}

func (g *Generator) SetInputValue(key, value string) {
	g.InputValues[key] = value
}

func (g *Generator) GroupedMatricesWithMetadata() map[string][]map[string]interface{} {
	result := make(map[string][]map[string]interface{})

	appMap := make(map[string]string)
	for app := range g.Config.Matrix {
		appMap[strings.ToLower(app)] = app
	}

	for _, selectedApp := range g.SelectedApps {
		configApp := selectedApp
		if actualApp, ok := appMap[strings.ToLower(selectedApp)]; ok {
			configApp = actualApp
		}

		appConfig, ok := g.Config.Matrix[configApp]
		if !ok {
			continue
		}

		platformMap := make(map[string]string)
		for platform := range appConfig {
			platformMap[strings.ToLower(platform)] = platform
		}

		for _, selectedPlatform := range g.SelectedPlatforms {
			configPlatform := selectedPlatform
			if actualPlatform, ok := platformMap[strings.ToLower(selectedPlatform)]; ok {
				configPlatform = actualPlatform
			}

			platformConfig, ok := appConfig[configPlatform]
			if !ok {
				continue
			}

			envMap := make(map[string]string)
			for env := range platformConfig {
				envMap[strings.ToLower(env)] = env
			}

			for _, selectedEnv := range g.SelectedEnvironments {
				configEnv := selectedEnv
				if actualEnv, ok := envMap[strings.ToLower(selectedEnv)]; ok {
					configEnv = actualEnv
				}

				envConfig, ok := platformConfig[configEnv]
				if !ok {
					continue
				}

				workflow := envConfig.Workflow

				matrix := map[string]interface{}{
					"app":         configApp,
					"platform":    configPlatform,
					"environment": configEnv,
				}

				for k, v := range envConfig.Matrix {
					if strVal, ok := v.(string); ok {
						matrix[k] = g.Config.SubstituteVariables(strVal, g.InputValues)
					} else {
						matrix[k] = v
					}
				}

				if _, ok := result[workflow]; !ok {
					result[workflow] = []map[string]interface{}{}
				}
				result[workflow] = append(result[workflow], matrix)
			}
		}
	}

	return result
}

func (g *Generator) GroupedMatricesPurified() map[string][]map[string]interface{} {
	result := make(map[string][]map[string]interface{})

	appMap := make(map[string]string)
	for app := range g.Config.Matrix {
		appMap[strings.ToLower(app)] = app
	}

	for _, selectedApp := range g.SelectedApps {
		configApp := selectedApp
		if actualApp, ok := appMap[strings.ToLower(selectedApp)]; ok {
			configApp = actualApp
		}

		appConfig, ok := g.Config.Matrix[configApp]
		if !ok {
			continue
		}

		platformMap := make(map[string]string)
		for platform := range appConfig {
			platformMap[strings.ToLower(platform)] = platform
		}

		for _, selectedPlatform := range g.SelectedPlatforms {
			configPlatform := selectedPlatform
			if actualPlatform, ok := platformMap[strings.ToLower(selectedPlatform)]; ok {
				configPlatform = actualPlatform
			}

			platformConfig, ok := appConfig[configPlatform]
			if !ok {
				continue
			}

			envMap := make(map[string]string)
			for env := range platformConfig {
				envMap[strings.ToLower(env)] = env
			}

			for _, selectedEnv := range g.SelectedEnvironments {
				configEnv := selectedEnv
				if actualEnv, ok := envMap[strings.ToLower(selectedEnv)]; ok {
					configEnv = actualEnv
				}

				envConfig, ok := platformConfig[configEnv]
				if !ok {
					continue
				}

				workflow := envConfig.Workflow

				matrix := make(map[string]interface{})

				for k, v := range envConfig.Matrix {
					if strVal, ok := v.(string); ok {
						matrix[k] = g.Config.SubstituteVariables(strVal, g.InputValues)
					} else {
						matrix[k] = v
					}
				}

				if _, ok := result[workflow]; !ok {
					result[workflow] = []map[string]interface{}{}
				}
				result[workflow] = append(result[workflow], matrix)
			}
		}
	}

	return result
}

func (g *Generator) GroupedMatrices() map[string][]map[string]interface{} {
	return g.GroupedMatricesWithMetadata()
}

func (g *Generator) GetTotalCombinations() int {
	matrices := g.GroupedMatricesWithMetadata()
	total := 0

	for _, matrixList := range matrices {
		total += len(matrixList)
	}

	return total
}

func (g *Generator) FormatCompleteMatrixPreview(workflowName string) string {
	matrices := g.GroupedMatricesWithMetadata()

	if matrixList, ok := matrices[workflowName]; ok {
		if len(matrixList) == 0 {
			return "No matrix entries for this workflow."
		}

		var preview strings.Builder

		for i, matrix := range matrixList {
			preview.WriteString(fmt.Sprintf("Matrix #%d:\n", i+1))

			if app, ok := matrix["app"]; ok {
				preview.WriteString(fmt.Sprintf("  • App: %v\n", app))
			}

			if platform, ok := matrix["platform"]; ok {
				preview.WriteString(fmt.Sprintf("  • Platform: %v\n", platform))
			}

			if env, ok := matrix["environment"]; ok {
				preview.WriteString(fmt.Sprintf("  • Environment: %v\n", env))
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
				preview.WriteString(fmt.Sprintf("    - %s: %v\n", k, matrix[k]))
			}

			preview.WriteString("\n")
		}

		return preview.String()
	}

	return "No matrix entries for workflow: " + workflowName
}
