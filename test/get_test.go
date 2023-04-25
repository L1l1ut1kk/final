package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	control "rest/src/controllers"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetLatestPhotos(t *testing.T) {
	// create a mock gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// call the function being tested
	control.GetLatestPhotos(c)

	// check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

	// check the response body
	var images []gin.H
	err := json.Unmarshal(w.Body.Bytes(), &images)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err.Error())
	}

	// check that the response contains three images
	if len(images) != 3 {
		t.Errorf("Expected 3 images but got %d", len(images))
	}

	// check that each image has an "img_base64" and "id" field
	for _, image := range images {
		if _, ok := image["img_base64"]; !ok {
			t.Errorf("Expected image to have 'img_base64' field")
		}
		if _, ok := image["id"]; !ok {
			t.Errorf("Expected image to have 'id' field")
		}
	}
}
