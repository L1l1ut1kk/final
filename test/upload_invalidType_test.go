package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePhotoWithInvalidFileType(t *testing.T) {
	url := "http://localhost:8080/negative_image"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile := os.Open("/home/lutik/Desktop/api/final/test/test_up/test.gif")
	if errFile != nil {
		t.Fatal(errFile)
	}
	defer file.Close()
	part, errFile := writer.CreateFormFile("photo", filepath.Base(file.Name()))
	if errFile != nil {
		t.Fatal(errFile)
	}
	_, errFile = io.Copy(part, file)
	if errFile != nil {
		t.Fatal(errFile)
	}

	// закрываем writer после создания формы
	err := writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}
	var response map[string]string
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "invalid request: unsupported file type", response["error"])
}
