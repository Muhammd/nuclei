package cloud

import (
	"testing"

	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

func TestCloudClient(t *testing.T) {
	cfg := &config.Config{}
	err := cfg.PopulateFromFlags([]string{
		"api-key=efe61560-f223-4caa-ab6c-96fb49bd7d2e",
		"api-url=http://gcp-cloud-dev.nuclei.sh",
	})
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
