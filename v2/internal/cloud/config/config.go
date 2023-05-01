package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	jsoniter "github.com/json-iterator/go"
)

// Config contains the configuration options for the client.
//
// The configuration can be read following ways -
//   - From JSON Configuration file
//   - From ENV variables
//   - From command line flags
//
// The iterface fields can be either string, int64 or NameAndID struct.
type Config struct {
	// CLI Flags

	// Interactive is a boolean flag to force enable interactive mode.
	Interactive bool `json:"interactive"`
	// NonInteractive is a boolean flag to force disable interactive mode.
	NonInteractive bool `json:"non-interactive"`

	// API Definitions

	// APIKey is the API key to use for authentication.
	APIKey string `json:"api-key" env:"NUCLEI_API_KEY"`
	// APIURL is the URL of the Cloud server
	APIURL string `json:"api-url" env:"NUCLEI_CLOUD_SERVER"`

	// Cloud resources

	// Workspace is the workspace to use for the client.
	Workspace string `json:"workspace"`
	// Project is the project to use for the client.
	Project string `json:"project"`

	// TemplatePlaylist is the template playlists to use for the scan.
	TemplatePlaylist []string `json:"template-playlists"`
	// AssetPlaylist is the asset playlists to use for the scan.
	AssetPlaylist []string `json:"asset-playlists"`
	// TODO: Add more fields for the configuration structure

	// Internal fields

	// InternalIDs contains the internal IDs of the workspace and project.
	// Internally maintained by the client and not exposed to the user.
	InternalIDs InternalIDs `json:"internal-ids"`
}

// InternalIDs contains the internal IDs of the workspace and project.
type InternalIDs struct {
	// WorkspaceID is the ID of the workspace to use for the client.
	WorkspaceID int64 `json:"workspace-id"`
	// ProjectID is the ID of the project to use for the client.
	ProjectID int64 `json:"project-id"`
	// TemplatePlaylistIDs is the IDs of the template playlists to use for the client.
	TemplatePlaylistIDs []int64 `json:"template-playlist-ids"`
	// AssetPlaylistIDs is the IDs of the asset playlists to use for the client.
	AssetPlaylistIDs []int64 `json:"asset-playlist-ids"`
}

// NameAndID contains the name and ID of an entity.
type NameAndID struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// PopulateFromEnv populates the configuration struct from the environment variables.
func (c *Config) PopulateFromEnv() error {
	if apikey := os.Getenv("NUCLEI_API_KEY"); apikey != "" {
		c.APIKey = apikey
	}
	if apiurl := os.Getenv("NUCLEI_CLOUD_SERVER"); apiurl != "" {
		c.APIURL = apiurl
	}
	return nil
}

// PopulateFromFlags populates the configuration struct from the command line flags.
func (c *Config) PopulateFromFlags(flags []string) error {
	for _, flag := range flags {
		keyValue := strings.SplitN(flag, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		key, userValue := keyValue[0], keyValue[1]

		var setError error
		c.iterateFields(func(name, tag string, value interface{}, field *structs.Field) bool {
			if tag == key {
				// Handle the special case of slice
				if field.Kind() == reflect.Slice {
					sliceValues := strings.Split(userValue, ",")
					if setError = field.Set(sliceValues); setError != nil {
						return false
					}
					return true
				}
				// Set the value of the field
				if setError = field.Set(userValue); setError != nil {
					return false
				}
			}
			return true
		})
	}
	return nil
}

// PopulateFromFile populates the configuration struct from the configuration file.
func (c *Config) PopulateFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jsoniter.NewDecoder(file).Decode(c); err != nil {
		return err
	}
	return nil
}

// WriteToFile writes the configuration to the configuration file.
func (c *Config) WriteToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jsoniter.NewEncoder(file).Encode(c); err != nil {
		return err
	}
	return nil
}

// iterateFields iterates over the fields of the configuration struct.
func (c *Config) iterateFields(callback func(name, tag string, value interface{}, field *structs.Field) bool) {
	configData := structs.New(c)
	for _, field := range configData.Fields() {
		field := field
		jsonTag := field.Tag("json")
		if !callback(field.Name(), jsonTag, field.Value(), field) {
			break
		}
	}
}
