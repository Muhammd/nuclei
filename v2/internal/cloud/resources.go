package cloud

import (
	"fmt"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/playlists/assets"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/playlists/templates"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/projects"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/workspaces"
)

// preRunActions performs pre-run actions for the client.
//
// It accepts cloud resource details from the user interactively
// if not specified as part of the configuration already or previously
// in JSON files or forced by user.
func (c *Client) preRunActions() error {
	// If interactive mode is enabled, ask for resource details
	if c.config.Interactive || c.config.Workspace == "" || c.config.Project == "" {
		if err := c.askForWorkspace(); err != nil {
			return err
		}
		if err := c.askForProject(); err != nil {
			return err
		}
		if err := c.askForTemplates(); err != nil {
			return err
		}
		if err := c.askForAssets(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) askForWorkspace() error {
	if c.config.NonInteractive {
		if c.config.Workspace == "" {
			return fmt.Errorf("workspace is required")
		}
		return nil
	}

	// Try to load workspace if not present
	if c.config.Workspace == "" {
		gologger.Info().Msgf("No workspace specified, loading workspaces...\n")

		model, err := workspaces.New(c.api)
		if err != nil {
			return err
		}
		workspace, id, err := model.Run()
		if err != nil {
			return err
		}
		c.config.Workspace = workspace
		c.config.InternalIDs.WorkspaceID = id
	}
	if c.config.Workspace == "" {
		return fmt.Errorf("workspace is required")
	}
	gologger.Info().Msgf("Using workspace: %s [%d]\n", c.config.Workspace, c.config.InternalIDs.WorkspaceID)
	return nil
}

func (c *Client) askForProject() error {
	if c.config.NonInteractive {
		if c.config.Project == "" {
			return fmt.Errorf("project is required")
		}
		return nil
	}

	// Try to load project if not present
	if c.config.Project == "" {
		gologger.Info().Msgf("No project specified, loading projects...\n")

		model, err := projects.New(c.api, c.config)
		if err != nil {
			return err
		}
		project, id, err := model.Run()
		if err != nil {
			return err
		}
		c.config.Project = project
		c.config.InternalIDs.ProjectID = id
	}
	if c.config.Project == "" {
		return fmt.Errorf("project is required")
	}
	gologger.Info().Msgf("Using project: %s [%d]\n", c.config.Project, c.config.InternalIDs.ProjectID)
	return nil
}

func (c *Client) askForTemplates() error {
	if c.config.NonInteractive {
		if len(c.config.TemplatePlaylist) == 0 {
			return fmt.Errorf("template playlist is required")
		}
		return nil
	}

	// Try to load template playlist if not present
	if len(c.config.TemplatePlaylist) == 0 {
		gologger.Info().Msgf("No template playlist specified, loading playlists...\n")

		model, err := templates.New(c.api, c.config)
		if err != nil {
			return err
		}
		templates, err := model.Run()
		if err != nil {
			return err
		}
		for _, template := range templates {
			c.config.TemplatePlaylist = append(c.config.TemplatePlaylist, template.ItemTitle)
			c.config.InternalIDs.TemplatePlaylistIDs = append(c.config.InternalIDs.TemplatePlaylistIDs, template.ID)
		}
	}
	if len(c.config.TemplatePlaylist) == 0 {
		return fmt.Errorf("template playlist is required")
	}
	gologger.Info().Msgf("Using template playlist: %v [%v]\n", c.config.TemplatePlaylist, c.config.InternalIDs.TemplatePlaylistIDs)
	return nil
}

func (c *Client) askForAssets() error {
	if c.config.NonInteractive {
		if len(c.config.AssetPlaylist) == 0 {
			return fmt.Errorf("asset playlist is required")
		}
		return nil
	}

	// Try to load asset playlist if not present
	if len(c.config.AssetPlaylist) == 0 {
		gologger.Info().Msgf("No asset playlist specified, loading playlists...\n")

		model, err := assets.New(c.api, c.config)
		if err != nil {
			return err
		}
		assets, err := model.Run()
		if err != nil {
			return err
		}
		for _, asset := range assets {
			c.config.AssetPlaylist = append(c.config.AssetPlaylist, asset.ItemTitle)
			c.config.InternalIDs.AssetPlaylistIDs = append(c.config.InternalIDs.AssetPlaylistIDs, asset.ID)
		}
	}
	if len(c.config.AssetPlaylist) == 0 {
		return fmt.Errorf("asset playlist is required")
	}
	gologger.Info().Msgf("Using asset playlist: %s [%v]\n", c.config.AssetPlaylist, c.config.InternalIDs.AssetPlaylistIDs)
	return nil
}
