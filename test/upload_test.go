package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFileType(t *testing.T) {
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

func TestSaveEmptyFile(t *testing.T) {
	url := "http://localhost:8080/negative_image"
	method := "POST"

	// Создаем временный файл с расширением PNG или JPEG
	tempFile, err := ioutil.TempFile("", "*.png")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Отправляем пустой файл PNG или JPEG
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile("photo", filepath.Base(tempFile.Name()))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, tempFile)
	if err != nil {
		t.Fatal(err)
	}

	// закрываем writer после создания формы
	err = writer.Close()
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

	// Проверяем, что сервер вернул ошибку и правильное сообщение об ошибке
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "invalid request: empty file", response["error"])
}

func TestSavePhoto_Success(t *testing.T) {

	url := "http://localhost:8080/negative_image"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/home/lutik/Desktop/api/final/test/test_up/test.jpg")
	if errFile1 != nil {
		t.Fatalf("Failed to open file: %v", errFile1)
	}
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("photo", filepath.Base("/home/lutik/Desktop/api/final/test/test_up/test.jpg"))
	if errFile1 != nil {
		t.Fatalf("Failed to create form file: %v", errFile1)
	}
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		t.Fatalf("Failed to copy file to form file: %v", errFile1)
	}
	err := writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Fatalf("Failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.Status)
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

}
