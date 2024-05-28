package runtrigger_test

import (
	"bytes"
	"fmt"
	types "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	var opts types.RunTriggerCreateOptions

	body := []byte(`
{
  "data": {
    "type": "run-triggers",
    "relationships": {
      "sourceable": {
        "data": {
          "type": "workspaces",
          "id": "ws-xYzax6kGKXRz4Gj6"
        }
      }
    }
  }
}
`)
	br := bytes.NewReader(body)
	if err := jsonapi.UnmarshalPayload(br, &opts); err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("\n%+v\n", opts.Sourceable)
}
