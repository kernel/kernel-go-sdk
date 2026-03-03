package kernel_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

// parseMultipartArrayField parses multipart form field names that encode an
// array of objects. It supports the three standard conventions:
//
//	files[0].dest_path   (bracket-dot)
//	files[0][dest_path]  (bracket-bracket)
//	files.0.dest_path    (dot-indexed)
//
// Returns (index, fieldName, true) on success.
func parseMultipartArrayField(prefix, name string) (int, string, bool) {
	if !strings.HasPrefix(name, prefix) {
		return 0, "", false
	}

	rest := name[len(prefix):]

	// bracket notation: [0].field or [0][field]
	if strings.HasPrefix(rest, "[") {
		end := strings.Index(rest, "]")
		if end == -1 {
			return 0, "", false
		}
		idxStr := rest[1:end]
		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 0 {
			return 0, "", false
		}
		after := rest[end+1:]
		after = strings.TrimPrefix(after, ".")
		if strings.HasPrefix(after, "[") && strings.HasSuffix(after, "]") {
			return idx, after[1 : len(after)-1], true
		}
		return idx, after, true
	}

	// dot notation: .0.field
	if strings.HasPrefix(rest, ".") {
		parts := strings.SplitN(rest[1:], ".", 2)
		if len(parts) != 2 {
			return 0, "", false
		}
		idx, err := strconv.Atoi(parts[0])
		if err != nil || idx < 0 {
			return 0, "", false
		}
		return idx, parts[1], true
	}

	return 0, "", false
}

// multipartUploadEntry represents one parsed file+dest_path pair from a
// multipart upload request.
type multipartUploadEntry struct {
	DestPath string
	Content  []byte
}

// parseUploadRequest reads a multipart request body and extracts the files
// array. It returns an error if field names don't use a decodable indexed
// format, which is the core of this test: the SDK must produce field names
// that a standard server can unambiguously parse.
func parseUploadRequest(r *http.Request) ([]multipartUploadEntry, error) {
	contentType := r.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("bad content-type: %w", err)
	}
	if mediaType != "multipart/form-data" {
		return nil, fmt.Errorf("expected multipart/form-data, got %s", mediaType)
	}
	boundary := params["boundary"]
	if boundary == "" {
		return nil, fmt.Errorf("missing boundary")
	}

	reader := multipart.NewReader(r.Body, boundary)
	type pending struct {
		destPath string
		content  []byte
	}
	items := map[int]*pending{}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading part: %w", err)
		}
		name := part.FormName()
		idx, field, ok := parseMultipartArrayField("files", name)
		if !ok {
			return nil, fmt.Errorf("cannot parse field name %q: expected indexed format like files[0].field or files.0.field", name)
		}
		p, exists := items[idx]
		if !exists {
			p = &pending{}
			items[idx] = p
		}
		switch field {
		case "dest_path":
			b, _ := io.ReadAll(part)
			p.destPath = strings.TrimSpace(string(b))
		case "file":
			p.content, _ = io.ReadAll(part)
		default:
			return nil, fmt.Errorf("unexpected field %q at index %d", field, idx)
		}
	}

	result := make([]multipartUploadEntry, len(items))
	for i := 0; i < len(items); i++ {
		p, ok := items[i]
		if !ok {
			return nil, fmt.Errorf("missing index %d (have indices up to %d)", i, len(items)-1)
		}
		if p.destPath == "" || p.content == nil {
			return nil, fmt.Errorf("index %d missing dest_path or file", i)
		}
		result[i] = multipartUploadEntry{DestPath: p.destPath, Content: p.content}
	}
	return result, nil
}

// multipartExtensionEntry represents one parsed name+zip_file pair from a
// multipart extensions upload request.
type multipartExtensionEntry struct {
	Name    string
	Content []byte
}

// parseExtensionsRequest is the same idea as parseUploadRequest but for the
// extensions array (name + zip_file).
func parseExtensionsRequest(r *http.Request) ([]multipartExtensionEntry, error) {
	contentType := r.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("bad content-type: %w", err)
	}
	if mediaType != "multipart/form-data" {
		return nil, fmt.Errorf("expected multipart/form-data, got %s", mediaType)
	}
	boundary := params["boundary"]
	if boundary == "" {
		return nil, fmt.Errorf("missing boundary")
	}

	reader := multipart.NewReader(r.Body, boundary)
	type pending struct {
		name    string
		content []byte
	}
	items := map[int]*pending{}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading part: %w", err)
		}
		formName := part.FormName()
		idx, field, ok := parseMultipartArrayField("extensions", formName)
		if !ok {
			return nil, fmt.Errorf("cannot parse field name %q: expected indexed format like extensions[0].field or extensions.0.field", formName)
		}
		p, exists := items[idx]
		if !exists {
			p = &pending{}
			items[idx] = p
		}
		switch field {
		case "name":
			b, _ := io.ReadAll(part)
			p.name = strings.TrimSpace(string(b))
		case "zip_file":
			p.content, _ = io.ReadAll(part)
		default:
			return nil, fmt.Errorf("unexpected field %q at index %d", field, idx)
		}
	}

	result := make([]multipartExtensionEntry, len(items))
	for i := 0; i < len(items); i++ {
		p, ok := items[i]
		if !ok {
			return nil, fmt.Errorf("missing index %d", i)
		}
		if p.name == "" || p.content == nil {
			return nil, fmt.Errorf("index %d missing name or zip_file", i)
		}
		result[i] = multipartExtensionEntry{Name: p.name, Content: p.content}
	}
	return result, nil
}

// TestUploadFilesMultipartEncoding verifies that the SDK encodes the files
// array with indexed field names that a server can unambiguously decode.
//
// The SDK's apiform.MarshalRoot uses arrayFmt:"comma" which drops array
// indices, producing field names like "files.dest_path" and "files.file"
// instead of "files.0.dest_path" and "files.0.file" (or bracket equivalents).
// Without indices, a server cannot pair file contents with their destination
// paths when multiple files are uploaded.
func TestUploadFilesMultipartEncoding(t *testing.T) {
	t.Run("single file", func(t *testing.T) {
		var parseErr error
		var entries []multipartUploadEntry

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entries, parseErr = parseUploadRequest(r)
			w.WriteHeader(201)
		}))
		defer server.Close()

		client := kernel.NewClient(
			option.WithBaseURL(server.URL),
			option.WithAPIKey("test-key"),
		)
		err := client.Browsers.Fs.Upload(context.Background(), "sess-1", kernel.BrowserFUploadParams{
			Files: []kernel.BrowserFUploadParamsFile{
				{DestPath: "/tmp/hello.txt", File: bytes.NewReader([]byte("hello world"))},
			},
		})
		if err != nil {
			t.Fatalf("SDK Upload call failed: %v", err)
		}
		if parseErr != nil {
			t.Fatalf("server could not parse multipart field names: %v", parseErr)
		}
		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}
		if entries[0].DestPath != "/tmp/hello.txt" {
			t.Errorf("dest_path = %q, want %q", entries[0].DestPath, "/tmp/hello.txt")
		}
		if !bytes.Equal(entries[0].Content, []byte("hello world")) {
			t.Errorf("content = %q, want %q", entries[0].Content, "hello world")
		}
	})

	t.Run("multiple files", func(t *testing.T) {
		var parseErr error
		var entries []multipartUploadEntry

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entries, parseErr = parseUploadRequest(r)
			w.WriteHeader(201)
		}))
		defer server.Close()

		client := kernel.NewClient(
			option.WithBaseURL(server.URL),
			option.WithAPIKey("test-key"),
		)
		err := client.Browsers.Fs.Upload(context.Background(), "sess-1", kernel.BrowserFUploadParams{
			Files: []kernel.BrowserFUploadParamsFile{
				{DestPath: "/tmp/a.txt", File: bytes.NewReader([]byte("aaa"))},
				{DestPath: "/tmp/b.txt", File: bytes.NewReader([]byte("bbb"))},
				{DestPath: "/tmp/c.txt", File: bytes.NewReader([]byte("ccc"))},
			},
		})
		if err != nil {
			t.Fatalf("SDK Upload call failed: %v", err)
		}
		if parseErr != nil {
			t.Fatalf("server could not parse multipart field names: %v", parseErr)
		}
		if len(entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(entries))
		}
		want := []struct {
			dest    string
			content string
		}{
			{"/tmp/a.txt", "aaa"},
			{"/tmp/b.txt", "bbb"},
			{"/tmp/c.txt", "ccc"},
		}
		for i, w := range want {
			if entries[i].DestPath != w.dest {
				t.Errorf("entries[%d].DestPath = %q, want %q", i, entries[i].DestPath, w.dest)
			}
			if !bytes.Equal(entries[i].Content, []byte(w.content)) {
				t.Errorf("entries[%d].Content = %q, want %q", i, entries[i].Content, w.content)
			}
		}
	})
}

// TestUploadFilesFieldNameFormat directly inspects the multipart field names
// the SDK produces, verifying they include array indices.
func TestUploadFilesFieldNameFormat(t *testing.T) {
	var fieldNames []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		_, params, _ := mime.ParseMediaType(contentType)
		reader := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := reader.NextPart()
			if err != nil {
				break
			}
			fieldNames = append(fieldNames, part.FormName())
			io.ReadAll(part)
		}
		w.WriteHeader(201)
	}))
	defer server.Close()

	client := kernel.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
	)
	_ = client.Browsers.Fs.Upload(context.Background(), "sess-1", kernel.BrowserFUploadParams{
		Files: []kernel.BrowserFUploadParamsFile{
			{DestPath: "/tmp/a.txt", File: bytes.NewReader([]byte("aaa"))},
			{DestPath: "/tmp/b.txt", File: bytes.NewReader([]byte("bbb"))},
		},
	})

	if len(fieldNames) == 0 {
		t.Fatal("no field names captured")
	}

	// Every field name should contain an index (numeric) to disambiguate
	// array elements. Acceptable patterns:
	//   files[0].dest_path, files[0][dest_path], files.0.dest_path
	// Unacceptable (current bug):
	//   files.dest_path, files.file
	indexedPattern := regexp.MustCompile(`^files[\[.](\d+)`)
	for _, name := range fieldNames {
		if !indexedPattern.MatchString(name) {
			t.Errorf("field name %q does not contain an array index; the comma arrayFmt drops indices making multi-file uploads ambiguous", name)
		}
	}

	t.Logf("Field names produced by SDK: %v", fieldNames)
}

// TestLoadExtensionsMultipartEncoding verifies that the SDK encodes the
// extensions array with indexed field names.
func TestLoadExtensionsMultipartEncoding(t *testing.T) {
	var parseErr error
	var entries []multipartExtensionEntry

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries, parseErr = parseExtensionsRequest(r)
		w.WriteHeader(201)
	}))
	defer server.Close()

	client := kernel.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
	)
	err := client.Browsers.LoadExtensions(context.Background(), "sess-1", kernel.BrowserLoadExtensionsParams{
		Extensions: []kernel.BrowserLoadExtensionsParamsExtension{
			{Name: "ext-a", ZipFile: bytes.NewReader([]byte("zip-a-data"))},
			{Name: "ext-b", ZipFile: bytes.NewReader([]byte("zip-b-data"))},
		},
	})
	if err != nil {
		t.Fatalf("SDK LoadExtensions call failed: %v", err)
	}
	if parseErr != nil {
		t.Fatalf("server could not parse multipart field names: %v", parseErr)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "ext-a" {
		t.Errorf("entries[0].Name = %q, want %q", entries[0].Name, "ext-a")
	}
	if !bytes.Equal(entries[0].Content, []byte("zip-a-data")) {
		t.Errorf("entries[0].Content = %q, want %q", entries[0].Content, "zip-a-data")
	}
	if entries[1].Name != "ext-b" {
		t.Errorf("entries[1].Name = %q, want %q", entries[1].Name, "ext-b")
	}
	if !bytes.Equal(entries[1].Content, []byte("zip-b-data")) {
		t.Errorf("entries[1].Content = %q, want %q", entries[1].Content, "zip-b-data")
	}
}
