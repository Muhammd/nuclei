package datasources

import (
	"context"
	"fmt"
	"regexp"

	jsoniter "github.com/json-iterator/go"
	"github.com/projectdiscovery/gologger"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type cloudlistConfiguration struct {
	Config string `json:"config"`
	Repo   string `json:"repo"`
}

type datasourcesFormConfig struct {
	Type     api.DatasourceType
	Fields   []form.Input
	Provider string
	Original map[string]interface{}
}

type cloudlistProviderName struct {
	Config string `json:"config"`
}

type cloudlistProviderMetadata struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
}

func (c *cloudlistProviderName) getProvider() (cloudlistProviderMetadata, error) {
	var list []cloudlistProviderMetadata
	if err := jsoniter.Unmarshal([]byte(c.Config), &list); err != nil {
		return cloudlistProviderMetadata{}, err
	}
	if len(list) == 0 {
		return cloudlistProviderMetadata{}, fmt.Errorf("no providers found")
	}
	return list[0], nil
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
		var metadata cloudlistProviderMetadata
		if datasource.Type == api.Cloudlist {
			var cfg cloudlistProviderName
			if err = jsoniter.Unmarshal(datasource.Metadata, &cfg); err != nil {
				gologger.Error().Msgf("Could not parse datasource metadata: %s\n", err)
				continue
			}
			metadata, err = cfg.getProvider()
		}
		if err != nil {
			gologger.Error().Msgf("Could not parse datasource metadata: %s\n", err)
			continue
		}
		items = append(items, styledlist.Item{
			ID:              datasource.Id,
			ItemTitle:       repoName,
			ItemDescription: formatDatasourceDescription(datasource.Type, metadata),
		})
	}
	return items, nil
}

func formatDatasourceDescription(Type api.DatasourceType, metadata cloudlistProviderMetadata) string {
	if Type == api.Cloudlist {
		return fmt.Sprintf("%s %s", cases.Title(language.AmericanEnglish).String(metadata.Provider), metadata.ID)
	}
	return cases.Title(language.AmericanEnglish).String(string(Type))
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

		if datasource.Type == api.Cloudlist {
			m.parseCloudlistConfigurationTemplate(datasource.Template)
			continue
		}
		fields, values := m.parseConfigurationTemplate(datasource.Type, datasource.Template)
		items = append(items, datasourcesFormConfig{
			Type:     datasource.Type,
			Fields:   fields,
			Original: values,
		})
	}
	return items, nil
}

func (m *Model) parseConfigurationTemplate(Type api.DatasourceType, data []byte) ([]form.Input, map[string]interface{}) {
	if Type == api.Cloudlist {
		return nil, nil
	}
	var values map[string]interface{}
	err := jsoniter.Unmarshal(data, &values)
	if err != nil {
		return nil, nil
	}

	var inputs []form.Input
	for key, value := range values {
		// Skip metadata keys
		if key == "metadata" || key == "id" {
			continue
		}
		inputs = append(inputs, form.Input{
			PlaceHolder: value.(string),
			CharLimit:   100,
			Width:       40,
			Prompt:      convertPromptKeyToInput(key, " "),
			Original:    key,
		})
	}
	return inputs, values
}

func (m *Model) parseCloudlistConfigurationTemplate(data []byte) {
	var values cloudlistConfiguration
	err := jsoniter.Unmarshal(data, &values)
	if err != nil {
		gologger.Error().Msgf("Could not parse cloudlist configuration template: %s\n", err)
		return
	}
	var list []map[string]interface{}
	err = jsoniter.Unmarshal([]byte(values.Config), &list)
	if err != nil {
		gologger.Error().Msgf("Could not parse cloudlist configuration template: %s\n", err)
		return
	}

	m.cloudlistFormFields = make(map[string]datasourcesFormConfig)
	var providers []string
	var inputs []form.Input
	for _, item := range list {
		item := item

		orig := item["provider"].(string)
		provider := convertPromptKeyToInput(orig, "")
		providers = append(providers, provider)

		for key, value := range item {
			// Skip metadata keys
			if key == "provider" || key == "id" {
				continue
			}
			inputs = append(inputs, form.Input{
				PlaceHolder: value.(string),
				CharLimit:   100,
				Width:       40,
				Prompt:      convertPromptKeyToInput(key, " "),
				Original:    key,
			})
		}
		m.cloudlistFormFields[provider] = datasourcesFormConfig{
			Type:     api.Cloudlist,
			Fields:   inputs,
			Provider: orig,
			Original: item,
		}
		inputs = nil
	}
	m.cloudlistProviders = providers
}

var nonWordRegex = regexp.MustCompile(`[-_]`)

func convertPromptKeyToInput(key, replace string) string {
	replaced := nonWordRegex.ReplaceAllString(key, replace)
	title := cases.Title(language.AmericanEnglish).String(replaced)
	return title
}
