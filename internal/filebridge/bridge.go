package filebridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Request struct {
	Operation string `json:"operation"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Directory string `json:"directory,omitempty"`
}

type Response struct {
	Kinds     []string `json:"kinds"`
	Office    string   `json:"office"`
	Path      string   `json:"path"`
	Error     string   `json:"error"`
	Directory string   `json:"directory"`
	OpenAfter bool     `json:"openAfter"`
}

type Client struct{ Executable string }

func (c Client) Call(request Request) (Response, error) {
	var response Response
	data, err := json.Marshal(request)
	if err != nil {
		return response, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.Executable)
	cmd.Stdin = bytes.NewReader(data)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return response, fmt.Errorf("文件服务不可用：%w", err)
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return response, fmt.Errorf("无法读取文件服务响应：%w", err)
	}
	if response.Error != "" {
		return response, errors.New(response.Error)
	}
	return response, nil
}

var Kinds = []string{"txt", "md", "rtf", "json", "csv", "html", "docx", "xlsx", "pptx"}

func KnownKind(kind string) bool {
	for _, value := range Kinds {
		if kind == value {
			return true
		}
	}
	return false
}

// ParseURL never performs a write. Finder requests only prefill the visible form.
func ParseURL(raw string) (Request, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "newtofile" {
		return Request{}, errors.New("无效的新建文件链接")
	}
	if u.Host == "settings" {
		return Request{Operation: "settings"}, nil
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil || u.Host != "create" || !KnownKind(q.Get("type")) ||
		!filepath.IsAbs(q.Get("directory")) || strings.ContainsRune(q.Get("directory"), 0) {
		return Request{}, errors.New("无效的新建文件链接")
	}
	return Request{Operation: "create", Kind: q.Get("type"), Directory: q.Get("directory")}, nil
}
