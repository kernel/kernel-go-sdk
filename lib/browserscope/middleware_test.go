package browserscope

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kernel/kernel-go-sdk/option"
)

func TestMetroKernelMiddleware(t *testing.T) {
	mw := MetroKernelMiddleware("sess1", "tok")
	var final *http.Request
	next := func(req *http.Request) (*http.Response, error) {
		final = req
		return nil, nil
	}

	u, err := url.Parse("https://host/browser/kernel/browsers/sess1/process/exec?x=1")
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/browser/kernel/browsers/sess1/process/exec"
	req := &http.Request{URL: u, Header: http.Header{"Authorization": {"Bearer sk"}}}
	_, _ = mw(req, next)

	if final.Header.Get("Authorization") != "" {
		t.Fatal("authorization should be stripped")
	}
	if final.URL.Query().Get("jwt") != "tok" {
		t.Fatalf("jwt query: got %q", final.URL.Query().Get("jwt"))
	}
	if final.URL.Path != "/browser/kernel/process/exec" {
		t.Fatalf("path rewrite: got %s", final.URL.Path)
	}
}

func TestMetroKernelMiddlewarePreservesExistingJWT(t *testing.T) {
	mw := MetroKernelMiddleware("sess1", "tok")
	var final *http.Request
	next := func(req *http.Request) (*http.Response, error) {
		final = req
		return nil, nil
	}

	u, _ := url.Parse("https://host/browser/kernel/browsers/sess1/fs/list_files?jwt=already")
	u.Path = "/browser/kernel/browsers/sess1/fs/list_files"
	req := &http.Request{URL: u}
	_, _ = mw(req, next)
	if final.URL.Query().Get("jwt") != "already" {
		t.Fatalf("jwt: got %q want already", final.URL.Query().Get("jwt"))
	}
}

func TestMetroKernelMiddlewareType(t *testing.T) {
	var _ option.Middleware = MetroKernelMiddleware("a", "b")
}
