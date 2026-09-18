package browserrouting

import "testing"

func TestJWTFromWebSocketURL(t *testing.T) {
	const u = "wss://browser.example/browser/cdp?jwt=abc123&foo=bar"
	j, err := jwtFromWebSocketURL(u)
	if err != nil {
		t.Fatal(err)
	}
	if j != "abc123" {
		t.Fatalf("jwt: got %q want abc123", j)
	}
}

func TestJWTFromWebSocketURLMissing(t *testing.T) {
	_, err := jwtFromWebSocketURL("wss://browser.example/browser/cdp")
	if err == nil {
		t.Fatal("expected error")
	}
}
