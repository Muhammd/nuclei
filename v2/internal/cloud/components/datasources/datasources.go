package datasources

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/projectdiscovery/gologger"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
)

func (m *Model) initForm() {
	choices := m.choice.Choice()
	if len(choices) == 0 {
		return
	}

	choice := choices[0]
	var provider string
	if len(choices) > 1 {
		provider = choices[1]
	}
	inputs := m.getDatasourceProviderConfig(choice, provider)

	cfg := form.Config{
		Inputs:     inputs.Fields,
		InputStyle: inputStyle,
		OnSubmit: func(fields []string) error {
			if err := m.postAddDatasource(fields, inputs); err != nil {
				gologger.Fatal().Msgf("Could not add datasource: %s\n", err)
				return err
			}
			return nil
		},
		Prompt:      "Enter the details for the datasource",
		PromptStyle: titleStyle,
	}
	m.form = form.InitialModel(cfg, m.size)
}

func (m *Model) postAddDatasource(fields []string, form datasourcesFormConfig) error {
	var cfg string
	var err error
	// Build the JSON for metadata
	if form.Type == api.Cloudlist {
		cfg, err = m.buildCloudlistConfiguration(fields, form)
	} else {
		cfg, err = m.buildDatasourceConfiguration(fields, form)
	}
	if err != nil {
		return err
	}

	resp, err := m.client.PostWorkspacesWorkspaceIdDatasourcesWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
		api.PostWorkspacesWorkspaceIdDatasourcesJSONRequestBody{
			Metadata: cfg,
			Sync:     false,
			Type:     form.Type,
		},
	)
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		return fmt.Errorf("could not add datasource: %s", string(resp.Body))
	}
	var datasourceID int64
	if resp.JSON201 != nil {
		datasourceID = resp.JSON201.Id
		gologger.Info().Msgf("Datasource added with ID: %d\n", resp.JSON201.Id)
	}
	time.Sleep(1 * time.Second) // Wait for the datasource to be added

	// Get datasource assets and templates and allow users to auto
	// create playlists from them.
	if datasourceID != 0 {
		// FIXME: Add APIs to get assets and templates by datasourceID
		// to be added to a playlist on creation.

		//	m.getDatasourceAssets(datasourceID)
		//	m.getDatasourceTemplates(datasourceID)
	}
	items, err := m.getDatasources()
	if err != nil {
		return err
	}
	m.lists.SetItems(items, nil)
	m.state = stateList
	return nil
}

func (m *Model) buildDatasourceConfiguration(fields []string, form datasourcesFormConfig) (string, error) {
	value := m.buildConfigObjectFromKV(fields, form)
	var buf bytes.Buffer
	if err := jsoniter.NewEncoder(&buf).Encode(value); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (m *Model) buildCloudlistConfiguration(fields []string, form datasourcesFormConfig) (string, error) {
	// Create configuration structure for cloudlist
	cfg := m.buildConfigObjectFromKV(fields, form)
	cfgSlice := []map[string]interface{}{cfg}

	var buf bytes.Buffer
	if err := jsoniter.NewEncoder(&buf).Encode(cfgSlice); err != nil {
		return "", err
	}
	config := cloudlistConfiguration{Config: buf.String(), Repo: "cloudlist"}
	var second bytes.Buffer
	if err := jsoniter.NewEncoder(&second).Encode(config); err != nil {
		return "", err
	}
	return second.String(), nil
}

func (m *Model) buildConfigObjectFromKV(fields []string, form datasourcesFormConfig) map[string]interface{} {
	items := make(map[string]interface{}, len(fields))
	for i, field := range fields {
		v := form.Fields[i]
		items[v.Original] = field
	}
	if form.Provider != "" {
		items["provider"] = form.Provider
	}
	return items
}

func (m *Model) getDatasourceProviderConfig(choice, provider string) datasourcesFormConfig {
	var inputs datasourcesFormConfig
	if strings.EqualFold(choice, "cloudlist") {
		inputs = m.cloudlistFormFields[provider]
	} else {
		for _, input := range m.formFields {
			if strings.EqualFold(string(input.Type), choice) {
				inputs = input
				break
			}
		}
	}
	return inputs
}
