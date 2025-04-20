package filebrowser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
	Token    string // JWT token, optional
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

// Authenticate logs in and retrieves a JWT token.
func (c *Client) Authenticate() error {
	loginURL := fmt.Sprintf("%s/api/login", c.BaseURL)
	payload := map[string]string{
		"username": c.Username,
		"password": c.Password,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(loginURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: %s", resp.Status)
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	c.Token = result.Token
	return nil
}

// UploadFile uploads a file to the user's directory in filebrowser.
func (c *Client) UploadFile(remotePath, localPath string) error {
	uploadURL := fmt.Sprintf("%s/api/resources/%s", c.BaseURL, remotePath)
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("files", filepath.Base(localPath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadURL, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	return nil
}

// DownloadFile downloads a file from filebrowser.
func (c *Client) DownloadFile(remotePath, localPath string) error {
	dlURL := fmt.Sprintf("%s/api/raw/%s", c.BaseURL, remotePath)
	req, err := http.NewRequest("GET", dlURL, nil)
	if err != nil {
		return err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// ListFiles lists files in the user's directory.
func (c *Client) ListFiles(remoteDir string) ([]string, error) {
	listURL := fmt.Sprintf("%s/api/resources/%s", c.BaseURL, remoteDir)
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("list failed: %s", resp.Status)
	}
	var result struct {
		Items []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var files []string
	for _, item := range result.Items {
		if item.Type == "file" {
			files = append(files, item.Name)
		}
	}
	return files, nil
}
