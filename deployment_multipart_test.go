package kernel_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kernel/kernel-go-sdk"
)

func TestDeploymentNewMultipartJSONParts(t *testing.T) {
	params := kernel.DeploymentNewParams{
		RuntimePolicy: kernel.DeploymentNewParamsRuntimePolicy{Statements: []kernel.DeploymentNewParamsRuntimePolicyStatement{{Action: "invocations:write", Resource: "*"}}},
		Source:        kernel.DeploymentNewParamsSource{Type: "github", URL: "https://github.com/org/repo", Ref: "main", Entrypoint: "app.py"},
	}
	body, contentType, err := params.MarshalMultipart()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/deployments", bytes.NewReader(body))
	req.Header.Set("Content-Type", contentType)
	reader, err := req.MultipartReader()
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][]any{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if ct := part.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("part %q: expected application/json, got %q", part.FormName(), ct)
		}
		data, err := io.ReadAll(part)
		if err != nil {
			t.Fatal(err)
		}
		var value any
		if err := json.Unmarshal(data, &value); err != nil {
			t.Fatal(err)
		}
		got[part.FormName()] = append(got[part.FormName()], value)
	}
	want := map[string][]any{
		"source":         {map[string]any{"type": "github", "url": "https://github.com/org/repo", "ref": "main", "entrypoint": "app.py"}},
		"runtime_policy": {map[string]any{"statements": []any{map[string]any{"action": "invocations:write", "resource": "*"}}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("multipart parts = %#v, want %#v", got, want)
	}
}
