package workspace

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	types "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/http/decode"
	"github.com/tofutf/tofutf/internal/rbac"
	"github.com/tofutf/tofutf/internal/resource"
	"github.com/tofutf/tofutf/internal/tfeapi"
)

var (
	remoteExecutionMode = "remote"
	localExecutionMode  = "local"
)

type (
	// byWorkspaceName are parameters used when looking up a workspace by
	// name
	byWorkspaceName struct {
		Name         string `schema:"workspace_name,required"`
		Organization string `schema:"organization_name,required"`
	}

	byOrganizationName struct {
		Organization string `schema:"organization_name,required"`
	}

	byWorkspaceID struct {
		WorkspaceID string `schema:"workspace_id,required"`
	}

	// tfe implements the TFC/E workspaces API:
	//
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces
	tfe struct {
		*Service
		*tfeapi.Responder
	}
)

func (a *tfe) addHandlers(r *mux.Router) {
	r = r.PathPrefix(tfeapi.APIPrefixV2).Subrouter()

	r.HandleFunc("/organizations/{organization_name}/workspaces", a.listWorkspaces).Methods("GET")
	r.HandleFunc("/organizations/{organization_name}/workspaces", a.createWorkspace).Methods("POST")
	r.HandleFunc("/organizations/{organization_name}/workspaces/{workspace_name}", a.getWorkspaceByName).Methods("GET")
	r.HandleFunc("/organizations/{organization_name}/workspaces/{workspace_name}", a.updateWorkspaceByName).Methods("PATCH")
	r.HandleFunc("/organizations/{organization_name}/workspaces/{workspace_name}", a.deleteWorkspaceByName).Methods("DELETE")

	r.HandleFunc("/workspaces/{workspace_id}", a.updateWorkspaceByID).Methods("PATCH")
	r.HandleFunc("/workspaces/{workspace_id}", a.getWorkspace).Methods("GET")
	r.HandleFunc("/workspaces/{workspace_id}", a.deleteWorkspace).Methods("DELETE")
	r.HandleFunc("/workspaces/{workspace_id}/actions/lock", a.lockWorkspace).Methods("POST")
	r.HandleFunc("/workspaces/{workspace_id}/actions/unlock", a.unlockWorkspace).Methods("POST")
	r.HandleFunc("/workspaces/{workspace_id}/actions/force-unlock", a.forceUnlockWorkspace).Methods("POST")
}

func (a *tfe) createWorkspace(w http.ResponseWriter, r *http.Request) {
	var routeParams byOrganizationName
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}
	var params WorkspaceCreateOptions
	if err := jsonapi.UnmarshalPayload(r.Body, &params); err != nil {
		tfeapi.Error(w, err)
		return
	}

	opts := CreateOptions{
		AgentPoolID:                params.AgentPoolID,
		AllowDestroyPlan:           params.AllowDestroyPlan,
		AutoApply:                  params.AutoApply,
		Description:                params.Description,
		ExecutionMode:              (*ExecutionMode)(params.ExecutionMode),
		GlobalRemoteState:          params.GlobalRemoteState,
		MigrationEnvironment:       params.MigrationEnvironment,
		Name:                       params.Name,
		Organization:               &routeParams.Organization,
		QueueAllRuns:               params.QueueAllRuns,
		SpeculativeEnabled:         params.SpeculativeEnabled,
		SourceName:                 params.SourceName,
		SourceURL:                  params.SourceURL,
		StructuredRunOutputEnabled: params.StructuredRunOutputEnabled,
		TerraformVersion:           params.TerraformVersion,
		TriggerPrefixes:            params.TriggerPrefixes,
		TriggerPatterns:            params.TriggerPatterns,
		WorkingDirectory:           params.WorkingDirectory,
		// convert from json:api structs to tag specs
		Tags: toTagSpecs(params.Tags),
	}
	// Always trigger runs if neither trigger patterns nor tags regex are set
	if len(params.TriggerPatterns) == 0 && (params.VCSRepo == nil || params.VCSRepo.TagsRegex == nil) {
		opts.AlwaysTrigger = internal.Bool(true)
	}
	if params.Operations != nil {
		if params.ExecutionMode != nil {
			err := errors.New("operations is deprecated and cannot be specified when execution mode is used")
			tfeapi.Error(w, err)
			return
		}
		if *params.Operations {
			opts.ExecutionMode = ExecutionModePtr(RemoteExecutionMode)
		} else {
			opts.ExecutionMode = ExecutionModePtr(LocalExecutionMode)
		}
	}
	if params.VCSRepo != nil {
		if params.VCSRepo.Identifier == nil || params.VCSRepo.OAuthTokenID == nil {
			tfeapi.Error(w, errors.New("must specify both oauth-token-id and identifier attributes for vcs-repo"))
			return
		}
		opts.ConnectOptions = &ConnectOptions{
			RepoPath:      params.VCSRepo.Identifier,
			VCSProviderID: params.VCSRepo.OAuthTokenID,
			Branch:        params.VCSRepo.Branch,
			TagsRegex:     params.VCSRepo.TagsRegex,
		}
	}

	ws, err := a.Create(r.Context(), opts)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.JsonAPIResponse(w, r, http.StatusOK, converted, nil)

}

func (a *tfe) getWorkspace(w http.ResponseWriter, r *http.Request) {
	var routeParams byWorkspaceID
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.Get(r.Context(), routeParams.WorkspaceID)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.JsonAPIResponse(w, r, http.StatusOK, converted, nil)

}

func (a *tfe) getWorkspaceByName(w http.ResponseWriter, r *http.Request) {
	var params byWorkspaceName
	if err := decode.All(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.GetByName(r.Context(), params.Organization, params.Name)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	err = jsonapi.MarshalPayload(w, converted)
	if err != nil {
		tfeapi.Error(w, &internal.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *tfe) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	var routeParams byOrganizationName
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}
	var params types.WorkspaceListOptions
	if err := decode.All(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	page, err := a.List(r.Context(), ListOptions{
		Search:       params.Search,
		Organization: &routeParams.Organization,
		PageOptions:  resource.PageOptions(params.ListOptions),
		Tags:         internal.SplitCSV(params.Tags),
	})
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	// convert items
	items := make([]*types.Workspace, len(page.Items))
	for i, from := range page.Items {
		a.logger.Info("found", "item", from)
		if to, err := a.convert(from, r); err != nil {
			tfeapi.Error(w, err)
			return
		} else {
			items[i] = to
		}
	}

	a.JsonAPIResponse(w, r, http.StatusOK, items, page.Pagination)

}

// updateWorkspaceByID updates a workspace using its ID.
//
// TODO: support updating workspace's vcs repo.
func (a *tfe) updateWorkspaceByID(w http.ResponseWriter, r *http.Request) {
	var routeParams byWorkspaceID
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.updateWorkspace(w, r, routeParams.WorkspaceID)
}

// updateWorkspaceByName updates a workspace using its name and organization.
//
// TODO: support updating workspace's vcs repo.
func (a *tfe) updateWorkspaceByName(w http.ResponseWriter, r *http.Request) {
	var params byWorkspaceName
	if err := decode.Route(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.GetByName(r.Context(), params.Organization, params.Name)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.updateWorkspace(w, r, ws.ID)
}

func (a *tfe) lockWorkspace(w http.ResponseWriter, r *http.Request) {
	var routeParams byWorkspaceID
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.Lock(r.Context(), routeParams.WorkspaceID, nil)
	if err != nil {
		if errors.Is(err, ErrWorkspaceAlreadyLocked) {
			tfeapi.Error(w, &internal.HTTPError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			})
		} else {
			tfeapi.Error(w, err)
		}
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.JsonAPIResponse(w, r, http.StatusOK, converted, nil)

}

func (a *tfe) unlockWorkspace(w http.ResponseWriter, r *http.Request) {
	a.unlock(w, r, false)
}

func (a *tfe) forceUnlockWorkspace(w http.ResponseWriter, r *http.Request) {
	a.unlock(w, r, true)
}

func (a *tfe) unlock(w http.ResponseWriter, r *http.Request, force bool) {
	var routeParams byWorkspaceID
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.Unlock(r.Context(), routeParams.WorkspaceID, nil, force)
	if err != nil {
		if errors.Is(err, ErrWorkspaceAlreadyUnlocked) || errors.Is(err, ErrWorkspaceLockedByRun) {
			tfeapi.Error(w, &internal.HTTPError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			})
		} else {
			tfeapi.Error(w, err)
		}
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.JsonAPIResponse(w, r, http.StatusOK, converted, nil)

}

func (a *tfe) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	var routeParams byWorkspaceID
	if err := decode.Route(&routeParams, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	_, err := a.Delete(r.Context(), routeParams.WorkspaceID)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *tfe) deleteWorkspaceByName(w http.ResponseWriter, r *http.Request) {
	var params byWorkspaceName
	if err := decode.All(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	ws, err := a.GetByName(r.Context(), params.Organization, params.Name)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	_, err = a.Delete(r.Context(), ws.ID)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *tfe) updateWorkspace(w http.ResponseWriter, r *http.Request, workspaceID string) {
	params := WorkspaceUpdateOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, &params); err != nil {
		tfeapi.Error(w, err)
		return
	}
	if err := ValidateWorkspaceUpdate(&params); err != nil {
		tfeapi.Error(w, err)
		return
	}

	opts := UpdateOptions{
		AgentPoolID:                params.AgentPoolID,
		AllowDestroyPlan:           params.AllowDestroyPlan,
		AutoApply:                  params.AutoApply,
		Description:                params.Description,
		ExecutionMode:              (*ExecutionMode)(params.ExecutionMode),
		GlobalRemoteState:          params.GlobalRemoteState,
		Name:                       params.Name,
		QueueAllRuns:               params.QueueAllRuns,
		SpeculativeEnabled:         params.SpeculativeEnabled,
		StructuredRunOutputEnabled: params.StructuredRunOutputEnabled,
		TerraformVersion:           params.TerraformVersion,
		TriggerPrefixes:            params.TriggerPrefixes,
		TriggerPatterns:            params.TriggerPatterns,
		WorkingDirectory:           params.WorkingDirectory,
	}

	// If file-triggers-enabled is set to false and tags regex is unspecified
	// then enable always trigger runs for this workspace.
	//
	// TODO: return error when client has sent incompatible combinations of
	// options:
	// (a) file-triggers-enabled=true and tags-regex=non-nil
	// (b) file-triggers-enabled=true and trigger-prefixes=empty
	// (b) trigger-prefixes=non-empty and tags-regex=non-nil
	if (params.FileTriggersEnabled != nil && !*params.FileTriggersEnabled) && (params.VCSRepo == nil || params.VCSRepo.TagsRegex == nil) {
		opts.AlwaysTrigger = internal.Bool(true)
	}

	if params.VCSRepo != nil {

		// client has provided non-null vcs options, which means they either
		// want to connect the workspace or modify the connection.
		opts.ConnectOptions = &ConnectOptions{
			RepoPath:      params.VCSRepo.Identifier,
			VCSProviderID: params.VCSRepo.OAuthTokenID,
			Branch:        params.VCSRepo.Branch,
			TagsRegex:     params.VCSRepo.TagsRegex,
		}
	} else {
		// client has explicitly set VCS options to null, which means they
		// want the workspace to be disconnected.
		opts.Disconnect = true
	}

	ws, err := a.Update(r.Context(), workspaceID, opts)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	converted, err := a.convert(ws, r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	a.JsonAPIResponse(w, r, http.StatusOK, converted, nil)

}

func (a *tfe) convert(from *Workspace, r *http.Request) (*types.Workspace, error) {
	subject, err := internal.SubjectFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	policy, err := a.GetPolicy(r.Context(), from.ID)
	if err != nil {
		return nil, err
	}
	perms := &types.WorkspacePermissions{
		CanLock:           subject.CanAccessWorkspace(rbac.LockWorkspaceAction, policy),
		CanUnlock:         subject.CanAccessWorkspace(rbac.UnlockWorkspaceAction, policy),
		CanForceUnlock:    subject.CanAccessWorkspace(rbac.UnlockWorkspaceAction, policy),
		CanQueueApply:     subject.CanAccessWorkspace(rbac.ApplyRunAction, policy),
		CanQueueDestroy:   subject.CanAccessWorkspace(rbac.ApplyRunAction, policy),
		CanQueueRun:       subject.CanAccessWorkspace(rbac.CreateRunAction, policy),
		CanDestroy:        subject.CanAccessWorkspace(rbac.DeleteWorkspaceAction, policy),
		CanReadSettings:   subject.CanAccessWorkspace(rbac.GetWorkspaceAction, policy),
		CanUpdate:         subject.CanAccessWorkspace(rbac.UpdateWorkspaceAction, policy),
		CanUpdateVariable: subject.CanAccessWorkspace(rbac.UpdateWorkspaceAction, policy),
	}

	to := &types.Workspace{
		ID: from.ID,
		Actions: &types.WorkspaceActions{
			IsDestroyable: true,
		},
		AllowDestroyPlan:     from.AllowDestroyPlan,
		AutoApply:            from.AutoApply,
		AutoApplyRunTrigger:  from.AutoApply,
		CanQueueDestroyPlan:  from.CanQueueDestroyPlan,
		CreatedAt:            from.CreatedAt,
		Description:          from.Description,
		Environment:          from.Environment,
		ExecutionMode:        string(from.ExecutionMode),
		GlobalRemoteState:    from.GlobalRemoteState,
		Locked:               from.Locked(),
		MigrationEnvironment: from.MigrationEnvironment,
		Name:                 from.Name,
		// Operations is deprecated but clients and go-tfe tests still use it
		Operations:                 from.ExecutionMode == "remote",
		Permissions:                perms,
		QueueAllRuns:               from.QueueAllRuns,
		SpeculativeEnabled:         from.SpeculativeEnabled,
		SourceName:                 from.SourceName,
		SourceURL:                  from.SourceURL,
		StructuredRunOutputEnabled: from.StructuredRunOutputEnabled,
		TerraformVersion:           from.TerraformVersion,
		TriggerPrefixes:            from.TriggerPrefixes,
		TriggerPatterns:            from.TriggerPatterns,
		WorkingDirectory:           from.WorkingDirectory,
		TagNames:                   from.Tags,
		UpdatedAt:                  from.UpdatedAt,
		Organization:               &types.Organization{Name: from.Organization},
	}
	if from.AgentPoolID != nil {
		to.AgentPool = &types.AgentPool{
			ID: *from.AgentPoolID,
		}
	}
	if len(from.TriggerPrefixes) > 0 || len(from.TriggerPatterns) > 0 {
		to.FileTriggersEnabled = true
	}
	if from.LatestRun != nil {
		to.CurrentRun = &types.Run{ID: from.LatestRun.ID}
	}

	// Add VCS repo to json:api struct if connected. NOTE: the terraform CLI
	// uses the presence of VCS repo to determine whether to allow a terraform
	// apply or not, displaying the following message if not:
	//
	//	Apply not allowed for workspaces with a VCS connection
	//
	//	A workspace that is connected to a VCS requires the VCS-driven workflow to ensure that the VCS remains the single source of truth.
	//
	// OTF permits the user to disable this behaviour by ommiting this info and
	// fool the terraform CLI into thinking its not a workspace with a VCS
	// connection.

	if from.Connection != nil {
		if !from.Connection.AllowCLIApply || !tfeapi.IsTerraformCLI(r) {
			to.VCSRepo = &types.VCSRepo{
				OAuthTokenID: from.Connection.VCSProviderID,
				Branch:       from.Connection.Branch,
				Identifier:   from.Connection.Repo,
				TagsRegex:    from.Connection.TagsRegex,
			}
		}
	}
	return to, nil
}

func (a *tfe) include(ctx context.Context, v any) ([]any, error) {
	dst := reflect.Indirect(reflect.ValueOf(v))

	// v must be a struct with a field named Workspace of type *types.Workspace
	if dst.Kind() != reflect.Struct {
		return nil, nil
	}
	field := dst.FieldByName("Workspace")
	if !field.IsValid() {
		return nil, nil
	}
	onlyID, ok := field.Interface().(*types.Workspace)
	if !ok {
		return nil, nil
	}
	// onlyID only contains the ID field, e.g. types.Workspace{ID:
	// "ws-123"}; so now retrieve the fully populated workspace, convert to a
	// tfe workspace and return.
	ws, err := a.Get(ctx, onlyID.ID)
	if err != nil {
		return nil, fmt.Errorf("retrieving workspace: %w", err)
	}
	include, err := a.convert(ws, (&http.Request{}).WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return []any{include}, nil
}

func (a *tfe) includeMany(ctx context.Context, v any) ([]any, error) {
	dst := reflect.Indirect(reflect.ValueOf(v))

	// v must be a struct with a field named Workspaces of type []*types.Workspace
	if dst.Kind() != reflect.Struct {
		return nil, nil
	}
	field := dst.FieldByName("Workspaces")
	if !field.IsValid() {
		return nil, nil
	}
	onlyIDs, ok := field.Interface().([]*types.Workspace)
	if !ok {
		return nil, nil
	}
	// onlyIDs only contains the ID field, e.g. []*types.Workspace{{ID:
	// "ws-123"}}; so now retrieve the fully populated workspaces, convert and
	// return them.
	include := make([]any, len(onlyIDs))
	for i, onlyID := range onlyIDs {
		ws, err := a.Get(ctx, onlyID.ID)
		if err != nil {
			return nil, fmt.Errorf("retrieving workspace: %w", err)
		}
		include[i], err = a.convert(ws, (&http.Request{}).WithContext(ctx))
		if err != nil {
			return nil, err
		}
	}
	return include, nil
}

// Patched types

// WorkspaceCreateOptions represents the options for creating a new workspace.
type WorkspaceCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when: execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Optional: Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Optional: Whether to enable health assessments (drift detection etc.) for the workspace.
	// Reference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#create-a-workspace
	// Requires remote execution mode, Terraform Cloud Business entitlement, and a valid agent pool to work
	AssessmentsEnabled *bool `jsonapi:"attr,assessments-enabled,omitempty"`

	// Optional: Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Optional: Whether to automatically apply changes for runs that are created by run triggers
	// from another workspace.
	AutoApplyRunTrigger *bool `jsonapi:"attr,auto-apply-run-trigger,omitempty"`

	// Optional: A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Optional: Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// Optional: The legacy TFE environment to use as the source of the migration, in the
	// form organization/environment. Omit this unless you are migrating a legacy
	// environment.
	MigrationEnvironment *string `jsonapi:"attr,migration-environment,omitempty"`

	// The name of the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization.
	Name *string `jsonapi:"attr,name"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Whether this workspace allows speculative plans. Setting this to false
	// prevents Terraform Cloud or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// BETA. A friendly name for the application or client creating this
	// workspace. If set, this will be displayed on the workspace as
	// "Created via <SOURCE NAME>".
	SourceName *string `jsonapi:"attr,source-name,omitempty"`

	// BETA. A URL for the application or client creating this workspace. This
	// can be the URL of a related resource in another app, or a link to
	// documentation or other info about the client.
	SourceURL *string `jsonapi:"attr,source-url,omitempty"`

	// BETA. Enable the experimental advanced run user interface.
	// This only applies to runs using Terraform version 0.15.2 or newer,
	// and runs executed using older versions will see the classic experience
	// regardless of this setting.
	StructuredRunOutputEnabled *bool `jsonapi:"attr,structured-run-output-enabled,omitempty"`

	// The version of Terraform to use for this workspace. Upon creating a
	// workspace, the latest version is selected unless otherwise specified.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// Optional: List of patterns used to match against changed files in order
	// to decide whether to trigger a run or not.
	TriggerPatterns []string `jsonapi:"attr,trigger-patterns,omitempty"`

	// Settings for the workspace's VCS repository. If omitted, the workspace is
	// created without a VCS repo. If included, you must specify at least the
	// oauth-token-id and identifier keys below.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching the
	// environment when multiple environments exist within the same repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`

	// A list of tags to attach to the workspace. If the tag does not already
	// exist, it is created and added to the workspace.
	Tags []*types.Tag `jsonapi:"relation,tags,omitempty"`

	// Optional: Struct of booleans, which indicate whether the workspace
	// specifies its own values for various settings. If you mark a setting as
	// `false` in this struct, it will clear the workspace's existing value for
	// that setting and defer to the default value that its project or
	// organization provides.
	//
	// In general, it's not necessary to mark a setting as `true` in this
	// struct; if you provide a literal value for a setting, Terraform Cloud will
	// automatically update its overwrites field to `true`. If you do choose to
	// manually mark a setting as overwritten, you must provide a value for that
	// setting at the same time.
	SettingOverwrites *types.WorkspaceSettingOverwritesOptions `jsonapi:"attr,setting-overwrites,omitempty"`

	// Associated Project with the workspace. If not provided, default project
	// of the organization will be assigned to the workspace.
	Project *types.Project `jsonapi:"relation,project,omitempty"`
}

// WorkspaceUpdateOptions represents the options for updating a workspace.
type WorkspaceUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when: execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Optional: Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Optional: Whether to enable health assessments (drift detection etc.) for the workspace.
	// Reference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#update-a-workspace
	// Requires remote execution mode, Terraform Cloud Business entitlement, and a valid agent pool to work
	AssessmentsEnabled *bool `jsonapi:"attr,assessments-enabled,omitempty"`

	// Optional: Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Optional: Whether to automatically apply changes for runs that are created by run triggers
	// from another workspace.
	AutoApplyRunTrigger *bool `jsonapi:"attr,auto-apply-run-trigger,omitempty"`

	// Optional: A new name for the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization. Warning: Changing a workspace's name changes its URL in the
	// API and UI.
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Optional: Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	// Optional:
	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Optional: Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Optional: Whether this workspace allows speculative plans. Setting this to false
	// prevents Terraform Cloud or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// BETA. Enable the experimental advanced run user interface.
	// This only applies to runs using Terraform version 0.15.2 or newer,
	// and runs executed using older versions will see the classic experience
	// regardless of this setting.
	StructuredRunOutputEnabled *bool `jsonapi:"attr,structured-run-output-enabled,omitempty"`

	// Optional: The version of Terraform to use for this workspace.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// Optional: List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// Optional: List of patterns used to match against changed files in order
	// to decide whether to trigger a run or not.
	TriggerPatterns []string `jsonapi:"attr,trigger-patterns,omitempty"`

	// Optional: To delete a workspace's existing VCS repo, specify null instead of an
	// object. To modify a workspace's existing VCS repo, include whichever of
	// the keys below you wish to modify. To add a new VCS repo to a workspace
	// that didn't previously have one, include at least the oauth-token-id and
	// identifier keys.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// Optional: A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching
	// the environment when multiple environments exist within the same
	// repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`

	// Optional: Struct of booleans, which indicate whether the workspace
	// specifies its own values for various settings. If you mark a setting as
	// `false` in this struct, it will clear the workspace's existing value for
	// that setting and defer to the default value that its project or
	// organization provides.
	//
	// In general, it's not necessary to mark a setting as `true` in this
	// struct; if you provide a literal value for a setting, Terraform Cloud will
	// automatically update its overwrites field to `true`. If you do choose to
	// manually mark a setting as overwritten, you must provide a value for that
	// setting at the same time.
	SettingOverwrites *types.WorkspaceSettingOverwritesOptions `jsonapi:"attr,setting-overwrites,omitempty"`

	// Associated Project with the workspace. If not provided, default project
	// of the organization will be assigned to the workspace
	Project *types.Project `jsonapi:"relation,project,omitempty"`
}

// TODO: move this struct out. VCSRepoOptions is used by workspaces, policy sets, and registry modules
// VCSRepoOptions represents the configuration options of a VCS integration.
type VCSRepoOptions struct {
	Branch            *string `jsonapi:"attr,branch"`
	Identifier        *string `jsonapi:"attr,identifier"`
	IngressSubmodules *bool   `jsonapi:"attr,ingress-submodules"`
	OAuthTokenID      *string `jsonapi:"attr,oauth-token-id"`
	TagsRegex         *string `jsonapi:"attr,tags-regex"`
	GHAInstallationID *string `jsonapi:"attr,github-app-installation-id"`
}

func ValidateWorkspaceUpdate(opts *WorkspaceUpdateOptions) error {
	if opts.Operations != nil && opts.ExecutionMode != nil {
		return errors.New("operations is deprecated and cannot be specified when execution mode is used")
	}
	if opts.Operations != nil {
		if *opts.Operations {
			opts.ExecutionMode = &remoteExecutionMode
		} else {
			opts.ExecutionMode = &localExecutionMode
		}
	}
	return nil
}
