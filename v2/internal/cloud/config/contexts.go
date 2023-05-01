package config

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Contexts are unique nuclei cloud instances for a user
// organized by workspace and project name.
//
// Each unique context is stored in the user directory.
// The context is used to cache the cloud resources locally
//
// FIXME: Integrate with nuclei cloud codebase
type Contexts struct {
	directory string
}

// NewContexts creates a new contexts instance
func NewContexts(directory string) *Contexts {
	return &Contexts{directory: directory}
}

// Create creates a new context for a user
//
// It overrides the existing context if it exists
func (c *Contexts) Create(config *Config) error {
	fileName := formatContextName(config)

	if err := config.WriteToFile(fileName); err != nil {
		return errors.Wrap(err, "could not write context file")
	}
	return nil
}

// Get gets a context for a user
func (c *Contexts) Get(config *Config) error {
	fileName := formatContextName(config)

	if err := config.PopulateFromFile(fileName); err != nil {
		return errors.Wrap(err, "could not read context file")
	}
	return nil
}

// List lists all the contexts for a user
func (c *Contexts) List() ([]string, error) {
	return listFiles(c.directory)
}

func listFiles(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, errors.Wrap(err, "could not read directory")
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func formatContextName(config *Config) string {
	var builder strings.Builder
	builder.WriteString(config.Workspace)
	builder.WriteString("-")
	builder.WriteString(config.Project)
	return builder.String()
}
