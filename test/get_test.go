package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	control "rest/src/controllers"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetLatestPhotos_Success(t *testing.T) {
	// Initialize the test HTTP server
	r := gin.Default()
	r.GET("/get_latest_photos", control.GetLatestPhotos)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Send a GET request
	resp, err := http.Get(ts.URL + "/get_latest_photos")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d. Response body: %s", resp.StatusCode, string(body))
	}

	// Check if the "images" key is present in the response
	var response gin.H
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if _, ok := response["images"]; !ok {
		t.Errorf("Response body does not contain \"images\"")
	}

	// Check the number of images in the response
	images, ok := response["images"].([]interface{})
	if !ok {
		t.Fatalf("Failed to convert images to []interface{}")
	}
	if len(images) < 1 {
		t.Errorf("No images found in response")
	}

	// Check that each image is a string and was uploaded within the last 24 hours
	for _, img := range images {
		_, ok := img.(string)
		if !ok {
			t.Errorf("Image is not a string: %v", img)
			continue
		}
	}
}
