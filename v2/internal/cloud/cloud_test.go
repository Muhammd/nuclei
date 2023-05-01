package cloud

import (
	"testing"

	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

func TestCloudClient(t *testing.T) {
	cfg := &config.Config{}
	err := cfg.PopulateFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Create a context and use it to manage sessions
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	_ = client
}
