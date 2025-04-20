package filebrowser

import (
	"bytes"
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestAuthenticate(t *testing.T) {
	// Mock server for login
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/login" {
			t.Fatalf("unexpected URL: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token": "testtoken"}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "user", "pass")
	if err := client.Authenticate(); err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if client.Token != "testtoken" {
		t.Errorf("expected token 'testtoken', got '%s'", client.Token)
	}
}

func TestUploadFile(t *testing.T) {
	// Create a temp file to upload
	tmpfile, err := os.CreateTemp("", "testupload")
	if err != nil {
		t.Fatalf("TempFile error: %v", err)
	}
	tmpfile.WriteString("hello world")
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Mock server for upload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("missing or wrong auth header")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "")
	client.Token = "testtoken"
	if err := client.UploadFile("upload.txt", tmpfile.Name()); err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}
}

func TestDownloadFile(t *testing.T) {
	// Mock server for download
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("missing or wrong auth header")
		}
		w.Write([]byte("downloaded content"))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "")
	client.Token = "testtoken"
	tmpfile, err := os.CreateTemp("", "testdownload")
	if err != nil {
		t.Fatalf("TempFile error: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	if err := client.DownloadFile("remote.txt", tmpfile.Name()); err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !bytes.Equal(data, []byte("downloaded content")) {
		t.Errorf("unexpected file content: %s", string(data))
	}
}

func TestListFiles(t *testing.T) {
	// Mock server for list
	jsonResp := `{"items": [{"name": "file1.txt", "type": "file"}, {"name": "dir1", "type": "dir"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonResp))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "")
	client.Token = "testtoken"
	files, err := client.ListFiles("/")
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}
	if len(files) != 1 || files[0] != "file1.txt" {
		t.Errorf("unexpected files: %v", files)
	}
}
