package templates

import (
	"context"
	"fmt"
	"strings"

	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

func (m *Model) getTemplatePlaylists() ([]styledlist.Item, error) {
	resp, err := m.client.GetWorkspacesWorkspaceIdTemplatePlaylistsWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
	)
	if err != nil {
		return nil, err
	}
	var items []styledlist.Item
	for _, playlist := range *resp.JSON200 {
		playlist := playlist
		description := fmt.Sprintf("%s (%d templates)", playlist.Tags, playlist.Count)
		items = append(items, styledlist.Item{
			ID:              playlist.Id,
			ItemTitle:       playlist.Name,
			ItemDescription: strings.Trim(description, " "),
		})
	}
	return items, nil
}

func (m *Model) getTemplates() ([]styledlist.Item, error) {
	resp, err := m.client.GetWorkspacesWorkspaceIdTemplatesWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
	)
	if err != nil {
		return nil, err
	}
	var items []styledlist.Item
	for _, template := range *resp.JSON200 {
		template := template
		items = append(items, styledlist.Item{
			ID:              template.Id,
			ItemTitle:       template.Name,
			ItemDescription: formatTemplateDescription(template.Title, template.Protocol, template.Severity, template.Tags),
		})
	}
	return items, nil
}

func formatTemplateDescription(title, protocol, severity, tags string) string {
	var builder strings.Builder
	builder.WriteString(title)
	if protocol != "" {
		builder.WriteString(fmt.Sprintf(" (%s)", protocol))
	}
	if severity != "" {
		builder.WriteString(fmt.Sprintf(" [%s]", severity))
	}
	if tags != "" {
		builder.WriteString(fmt.Sprintf(" %s", tags))
	}
	return strings.Trim(builder.String(), " ")
}

func (m *Model) addTemplatePlaylist(name string, templateIDs []int64) error {
	_, err := m.client.PostWorkspacesWorkspaceIdTemplatePlaylistsWithResponse(context.Background(), m.config.InternalIDs.WorkspaceID, api.PostWorkspacesWorkspaceIdTemplatePlaylistsJSONRequestBody{
		Name:        name,
		TemplateIds: templateIDs,
	})
	return err
}
