package workspace_test

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/tofutf/tofutf/internal/workspace"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	body := []byte(`
{
  "data": {
    "type": "workspaces",
    "attributes": {
      "agent-pool-id": "apool-Hg2k85ThsXfn4Td9",
      "allow-destroy-plan": true,
      "assessments-enabled": false,
      "auto-apply": true,
      "auto-apply-run-trigger": true,
      "description": "",
      "execution-mode": "agent",
      "file-triggers-enabled": true,
      "global-remote-state": true,
      "name": "organization-init",
      "queue-all-runs": false,
      "speculative-enabled": true,
      "structured-run-output-enabled": true,
      "terraform-version": "1.8.3",
      "trigger-prefixes": [
        "*"
      ],
      "vcs-repo": {
        "identifier": "pvotal-tech/pvotal-tech-lowops-organization-manifests",
        "ingress-submodules": false,
        "oauth-token-id": "vcs-ExBCXw6buRwrKnDq"
      },
      "working-directory": "./entrypoints/organization"
    },
    "relationships": {
      "tags": {
        "data": [
          {
            "type": "tags",
            "attributes": {
              "name": "pvotal-tech"
            }
          },
          {
            "type": "tags",
            "attributes": {
              "name": "lowops"
            }
          }
        ]
      }
    }
  }
}
`)
	var opts workspace.WorkspaceCreateOptions
	br := bytes.NewReader(body)
	if err := jsonapi.UnmarshalPayload(br, &opts); err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, "pvotal-tech/pvotal-tech-lowops-organization-manifests", *opts.VCSRepo.Identifier)
	assert.Equal(t, "vcs-ExBCXw6buRwrKnDq", *opts.VCSRepo.OAuthTokenID)

	for _, tag := range opts.Tags {
		fmt.Printf("\n%+v\n", tag)
	}
}

func TestMarshalMany(t *testing.T) {

	wsl := &tfe.WorkspaceList{
		Items: []*tfe.Workspace{
			{
				ID:   "test1",
				Name: "test1 name",
			},
			{
				ID:   "test1",
				Name: "test1 name",
			},
		},
		Pagination: &tfe.Pagination{
			TotalPages:   1,
			CurrentPage:  1,
			TotalCount:   2,
			PreviousPage: 1,
			NextPage:     1,
		},
	}
	//br := bytes.Buffer{}
	if p, err := jsonapi.Marshal(wsl.Items); err != nil {
		t.Error(err)
		t.FailNow()
	} else {

		switch p.(type) {
		case *jsonapi.OnePayload:
			fmt.Println("one payload")
		case *jsonapi.ManyPayload:
			fmt.Println("many payload")
		}
	}
}
