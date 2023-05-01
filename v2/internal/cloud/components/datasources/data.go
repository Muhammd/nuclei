package datasources

import (
	"context"
	"regexp"

	jsoniter "github.com/json-iterator/go"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type datasourcesFormConfig struct {
	Type   api.DatasourceType
	Fields []form.Input
}

func (m *Model) getDatasources() ([]styledlist.Item, error) {
	resp, err := m.client.GetWorkspacesWorkspaceIdDatasourcesWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
		&api.GetWorkspacesWorkspaceIdDatasourcesParams{},
	)
	if err != nil {
		return nil, err
	}

	var items []styledlist.Item
	for _, datasource := range *resp.JSON200 {
		datasource := datasource

		repoName := jsoniter.Get(datasource.Metadata, "repo").ToString()
		items = append(items, styledlist.Item{
			ID:              datasource.Id,
			ItemTitle:       repoName,
			ItemDescription: string(datasource.Type),
		})
	}
	return items, nil
}

func (m *Model) getDatasourcesFormConfig() ([]datasourcesFormConfig, error) {
	resp, err := m.client.GetWorkspacesIdIntegrationsDatasourcesAvailableWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
	)
	if err != nil {
		return nil, err
	}

	var items []datasourcesFormConfig
	for _, datasource := range *resp.JSON200 {
		datasource := datasource

		items = append(items, datasourcesFormConfig{
			Type:   datasource.Type,
			Fields: parseConfigurationTemplate(datasource.Template),
		})
	}
	return items, nil
}

func parseConfigurationTemplate(data []byte) []form.Input {
	var values map[string]interface{}
	err := jsoniter.Unmarshal(data, &values)
	if err != nil {
		return nil
	}

	var inputs []form.Input
	for key, value := range values {
		// Skip metadata key
		if key == "metadata" {
			continue
		}
		inputs = append(inputs, form.Input{
			PlaceHolder: value.(string),
			CharLimit:   40,
			Width:       40,
			Prompt:      convertPromptKeyToInput(key),
		})
	}
	return inputs
}

var nonWordRegex = regexp.MustCompile(`[-_]`)

func convertPromptKeyToInput(key string) string {
	replaced := nonWordRegex.ReplaceAllString(key, " ")
	title := cases.Title(language.AmericanEnglish).String(replaced)
	return title
}
