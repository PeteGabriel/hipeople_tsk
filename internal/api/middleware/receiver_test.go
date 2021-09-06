package middleware

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckReceivedContent_RequestBodyIsEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/image", nil)
	w := httptest.NewRecorder()
	middleware := CheckReceivedContent(nil)
	middleware(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"title":"image not received","status":400,"detail":"An image must be provided.","instance":"/api/image"}`

	if resp.StatusCode != 400 {
		t.Errorf("expected status code: 400. Status code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("expected Content-Type: application/problem+json. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedBody {
		t.Errorf("unexpected body received: %s", string(body))
	}
}

func TestCheckReceivedContent(t *testing.T) {
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_, err := writer.CreateFormFile("upload_file", "image.png")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	writer.Close() //writes the 'boundary' content. http.Request type expected this content.
	//send an empty file. Just the boundary content is written there. For this purpose is enough.
	req := httptest.NewRequest(http.MethodPost, "/api/image", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	middleware := CheckReceivedContent(nil)
	middleware(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 { //default at this point of the middleware pipeline
		t.Errorf("expected status code: 200. Status code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "" {
		t.Errorf("Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != "" {
		t.Errorf("unexpected body received: %s", string(body))
	}
}
