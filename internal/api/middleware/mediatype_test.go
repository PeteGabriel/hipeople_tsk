package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
)

func TestValidateContentType(t *testing.T) {
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)

	//setup multipart file. Header is necessary.
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "upload_file", "image.png"))
	h.Set("Content-Type", "image/png") //multipart file contains the correct type

	_, err := writer.CreatePart(h)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	writer.Close() //writes the 'boundary' content. http.Request type expected this content.
	//send an empty file. Just the boundary content is written there. For this purpose is enough.
	req := httptest.NewRequest(http.MethodPost, "/api/image", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	middleware := ValidateContentType(nil)
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

func TestValidateContentType_NotAnImage(t *testing.T) {
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "upload_file", "image.png"))
	//setup as a pdf file
	h.Set("Content-Type", "application/pdf")

	_, err := writer.CreatePart(h)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	writer.Close() //writes the 'boundary' content. http.Request type expected this content.
	//send an empty file. Just the boundary content is written there. For this purpose is enough.
	req := httptest.NewRequest(http.MethodPost, "/api/image", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	middleware := ValidateContentType(nil)
	middleware(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"title":"file received is invalid","status":400,"detail":"File type not valid for upload. File received must be an image.","instance":"/api/image"}`

	if resp.StatusCode != 400 { //default at this point of the middleware pipeline
		t.Errorf("expected status code: 400. Status code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedBody {
		t.Errorf("unexpected body received: %s", string(body))
	}
}