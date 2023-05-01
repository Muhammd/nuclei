package config

import "testing"

func TestConfig_PopulateFromFlags(t *testing.T) {
	type fields struct {
		APIKey      string
		APIURL      string
		Workspace   string
		Project     string
		InternalIDs InternalIDs
	}
	type args struct {
		flags []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid",
			fields: fields{},
			args: args{
				flags: []string{
					"workspace=projectdiscovery",
					"project=nuclei",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIKey:      tt.fields.APIKey,
				APIURL:      tt.fields.APIURL,
				Workspace:   tt.fields.Workspace,
				Project:     tt.fields.Project,
				InternalIDs: tt.fields.InternalIDs,
			}
			if err := c.PopulateFromFlags(tt.args.flags); (err != nil) != tt.wantErr {
				t.Errorf("Config.PopulateFromFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Workspace != "projectdiscovery" {
				t.Errorf("Config.PopulateFromFlags() error = %v, wantErr %v", c.Workspace, "projectdiscovery")
			}
			if c.Project != "nuclei" {
				t.Errorf("Config.PopulateFromFlags() error = %v, wantErr %v", c.Project, "nuclei")
			}
		})
	}
}
