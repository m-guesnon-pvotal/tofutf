package run

import (
	"context"
	"testing"

	types "github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/configversion"
	"github.com/tofutf/tofutf/internal/vcs"
	"github.com/tofutf/tofutf/internal/workspace"
)

func TestReporter_HandleRun(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		run  *Run
		ws   *types.Workspace
		cv   *configversion.ConfigurationVersion
		want vcs.SetStatusOptions
	}{
		{
			name: "pending run",
			run:  &Run{ID: "run-123", Status: RunPending},
			ws: &types.Workspace{
				Name:    "dev",
				VCSRepo: &types.VCSRepo{},
			},
			cv: &configversion.ConfigurationVersion{
				IngressAttributes: &configversion.IngressAttributes{
					CommitSHA: "abc123",
					Repo:      "leg100/otf",
				},
			},
			want: vcs.SetStatusOptions{
				Workspace: "dev",
				Ref:       "abc123",
				Repo:      "leg100/otf",
				Status:    vcs.PendingStatus,
				TargetURL: "https://otf-host.org/app/runs/run-123",
			},
		},
		{
			name: "skip run with config not from a VCS repo",
			run:  &Run{ID: "run-123"},
			cv: &configversion.ConfigurationVersion{
				IngressAttributes: nil,
			},
			want: vcs.SetStatusOptions{},
		},
		{
			name: "skip UI-triggered run",
			run:  &Run{ID: "run-123", Source: SourceUI},
			want: vcs.SetStatusOptions{},
		},
		{
			name: "skip API-triggered run",
			run:  &Run{ID: "run-123", Source: SourceAPI},
			want: vcs.SetStatusOptions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got vcs.SetStatusOptions
			reporter := &Reporter{
				Workspaces:      &fakeReporterWorkspaceService{ws: tt.ws},
				Configs:         &fakeReporterConfigurationVersionService{cv: tt.cv},
				VCS:             &fakeReporterVCSProviderService{got: &got},
				HostnameService: internal.NewHostnameService("otf-host.org"),
			}
			err := reporter.handleRun(ctx, tt.run)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

type fakeReporterConfigurationVersionService struct {
	configversion.Service

	cv *configversion.ConfigurationVersion
}

func (f *fakeReporterConfigurationVersionService) Get(context.Context, string) (*configversion.ConfigurationVersion, error) {
	return f.cv, nil
}

type fakeReporterWorkspaceService struct {
	workspace.Service

	ws *types.Workspace
}

func (f *fakeReporterWorkspaceService) Get(context.Context, string) (*types.Workspace, error) {
	return f.ws, nil
}

type fakeReporterVCSProviderService struct {
	got *vcs.SetStatusOptions
}

func (f *fakeReporterVCSProviderService) GetVCSClient(context.Context, string) (vcs.Client, error) {
	return &fakeReporterCloudClient{got: f.got}, nil
}

type fakeReporterCloudClient struct {
	vcs.Client

	got *vcs.SetStatusOptions
}

func (f *fakeReporterCloudClient) SetStatus(ctx context.Context, opts vcs.SetStatusOptions) error {
	*f.got = opts
	return nil
}
