package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MessageError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Message struct {
	Success bool           `json:"success"`
	Errors  []MessageError `json:"errors"`
}

func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

func ShouldOpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	dir := name[0 : strings.LastIndex(name, "/")+1]
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, perm); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(name, flag, perm)
}

func FailPrintf(format string, a ...any) {
	json.NewEncoder(os.Stdout).Encode(Message{
		Success: false,
		Errors: []MessageError{
			{Code: 0, Message: fmt.Sprintf(format, a...)},
		},
	})
}

type SecurityConfiguration struct {
	XAuthEmail string `json:"x_auth_email"`
	XAuthKey   string `json:"x_auth_key"`
}

func (c *SecurityConfiguration) Save() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("failed to get user home dir, cause: " + err.Error())
	}
	f, err := ShouldOpenFile(filepath.Join(dir, ".cf_cli_config"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("failed to open security configuration file, cause: " + err.Error())
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}

func OpenSecurityConfiguration() (*SecurityConfiguration, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("failed to get user home dir, cause: " + err.Error())
	}
	f, err := ShouldOpenFile(filepath.Join(dir, ".cf_cli_config"), os.O_RDONLY, 0644)
	if err != nil {
		return nil, errors.New("failed to open security configuration, cause: " + err.Error())
	}
	defer f.Close()
	var c SecurityConfiguration
	if err = json.NewDecoder(f).Decode(&c); err != nil {
		return nil, errors.New("failed to parse security configuration, cause: " + err.Error())
	}
	return &c, nil
}

type RequestOption func(r *http.Request)

func UsePathParameters(name, value string) RequestOption {
	return func(r *http.Request) {
		r.URL.Path = strings.ReplaceAll(r.URL.Path, fmt.Sprintf("{%s}", name), value)
	}
}

func UseQueryParameters(key, value string) RequestOption {
	return func(r *http.Request) {
		if key == "" || value == "" {
			return
		}
		query := r.URL.Query()
		query.Set(key, value)
		r.URL.RawQuery = query.Encode()
	}
}

func UseQueryParametersWithMap(sets map[string]string) RequestOption {
	return func(r *http.Request) {
		query := r.URL.Query()
		for key, value := range sets {
			if key == "" || value == "" {
				continue
			}
			query.Set(key, value)
		}
		r.URL.RawQuery = query.Encode()
	}
}

func UseSecurity(c *SecurityConfiguration) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("X-Auth-Email", c.XAuthEmail)
		r.Header.Set("X-Auth-Key", c.XAuthKey)
	}
}

func UseBody(body io.Reader) RequestOption {
	return func(r *http.Request) {
		rc, ok := body.(io.ReadCloser)
		if !ok {
			rc = io.NopCloser(body)
		}
		r.Body = rc
	}
}

func UseJSONBody(v any) RequestOption {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(v)
	return func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		r.Body = io.NopCloser(buf)
	}
}

func Request(method string, api string, opts ...RequestOption) {
	request, err := http.NewRequest(method, path.Join(BaseAPI, api), nil)
	if err != nil {
		FailPrintf("failed to NewRequest, cause: %s", err)
		return
	}

	for _, opt := range opts {
		opt(request)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		FailPrintf("failed to request, cause: %s", err)
		return
	}
	defer response.Body.Close()

	if response.Header.Get("Content-Type") == "text/plain" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			FailPrintf("failed to read response body, cause: %s", err)
			return
		}
		json.NewEncoder(os.Stdout).Encode(map[string]any{
			"success": true,
			"errors":  []any{},
			"result":  string(body),
		})
		return
	}

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		FailPrintf("failed to copy response body, cause: %s", err)
	}
	return
}
