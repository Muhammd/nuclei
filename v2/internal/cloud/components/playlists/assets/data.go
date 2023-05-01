package assets

import (
	"context"
	"fmt"
	"strings"

	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

func (m *Model) getAssetsPlaylists() ([]styledlist.Item, error) {
	resp, err := m.client.GetWorkspacesWorkspaceIdAssetPlaylistsWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
	)
	if err != nil {
		return nil, err
	}
	var items []styledlist.Item
	for _, playlist := range *resp.JSON200 {
		playlist := playlist
		description := formatDescription(playlist.Tags, playlist.AssetCount)
		items = append(items, styledlist.Item{
			ID:              playlist.Id,
			ItemTitle:       playlist.Name,
			ItemDescription: strings.Trim(description, " "),
		})
	}
	return items, nil
}

func formatDescription(tags *[]string, count *int64) string {
	var builder strings.Builder
	if tags != nil && len(*tags) > 0 {
		builder.WriteString(fmt.Sprintf("%v", *tags))
	}
	if count != nil {
		builder.WriteString(fmt.Sprintf(" (%v assets)", *count))
	}
	return strings.Trim(builder.String(), " ")
}

func (m *Model) getAssets() ([]styledlist.Item, error) {
	size := 100
	resp, err := m.client.GetWorkspacesWorkspaceIdAssetsWithResponse(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
		&api.GetWorkspacesWorkspaceIdAssetsParams{
			Size: &size,
		},
	)
	if err != nil {
		return nil, err
	}
	var items []styledlist.Item
	for _, asset := range *resp.JSON200 {
		asset := asset
		items = append(items, styledlist.Item{
			ID:              asset.AssetId,
			ItemTitle:       asset.AssetName,
			ItemDescription: asset.DataSourceName,
		})
	}
	return items, nil
}

func (m *Model) addAssetPlaylist(name string, assetIds []int64) error {
	_, err := m.client.PostWorkspacesWorkspaceIdAssetPlaylists(
		context.Background(),
		m.config.InternalIDs.WorkspaceID,
		api.PostWorkspacesWorkspaceIdAssetPlaylistsJSONRequestBody{
			AssetIds: assetIds,
			Name:     name,
		},
	)
	return err
}
