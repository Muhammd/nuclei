// Package cloud provides a client for the Cloud API.
package cloud

import (
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/pkg/errors"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

// Client is a client for the Cloud
type Client struct {
	config *config.Config
	api    *api.ClientWithResponses
}

// NewClient returns a new Client.
func NewClient(config *config.Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	if config.APIURL == "" {
		return nil, errors.New("API URL is required")
	}

	// Create nuclei cloud API client
	apiKeyProvider, apiKeyProviderErr := securityprovider.NewSecurityProviderApiKey("header", "X-API-Key", config.APIKey)
	if apiKeyProviderErr != nil {
		return nil, errors.Wrap(apiKeyProviderErr, "could not create user key provider")
	}
	client, err := api.NewClientWithResponses(config.APIURL, api.WithRequestEditorFn(apiKeyProvider.Intercept))
	if err != nil {
		return nil, errors.Wrap(err, "could not create user client")
	}

	cloudClient := &Client{config: config, api: client}

	if err := cloudClient.preRunActions(); err != nil {
		return nil, errors.Wrap(err, "could not run pre-run actions")
	}
	return cloudClient, nil
}
