package api

import (
	"bytes"
	"errors"
	"fmt"
	"hipeople_task/pkg/config"
	"hipeople_task/pkg/domain"
	"hipeople_task/pkg/mocks/services"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
)

func TestApp_GetImage(t *testing.T) {
	mock := &services.ImageServiceMock{}
	mock.GetImageMock = func(id string) ([]byte, *domain.Error) {
		return []byte("image content goes in here"), nil
	}

	app := App{
		imgService: mock,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/image/123frd45tg6", nil)
	w := httptest.NewRecorder()
	getImageHandler := app.GetImage()
	getImageHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200. Code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "image/*" {
		t.Errorf("expected Content-Type: image/*. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if bytes.Compare(body, []byte("image content goes in here")) != 0 {
		t.Errorf("unexpected body content: %s", string(body))
	}
}

func TestApp_GetImage_ImageNotFound(t *testing.T) {
	mock := &services.ImageServiceMock{}
	mock.GetImageMock = func(id string) ([]byte, *domain.Error) {
		return []byte{}, &domain.Error{
			Message: "image with id 1234 not found",
			Code:    404,
		}
	}

	app := App{
		imgService: mock,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	expectedJson := `{"title":"image not found","status":404,"detail":"image with id 1234 not found","instance":"/api/image/1234"}`

	req := httptest.NewRequest(http.MethodGet, "/api/image/1234", nil)
	w := httptest.NewRecorder()
	getImageHandler := app.GetImage()
	getImageHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 404 {
		t.Errorf("expected status code: 404. Code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("expected Content-Type: application/problem+json. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedJson {
		t.Errorf("unexpected body content: %s", string(body))
	}
}

func TestApp_GetImage_HttpMethodNotAllowed(t *testing.T) {

	app := App{
		imgService: nil,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/image/1234", nil)
	w := httptest.NewRecorder()
	getImageHandler := app.GetImage()
	getImageHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != 405 {
		t.Errorf("expected status code: 405. Code received: %d", resp.StatusCode)
	}
}

func TestApp_GetImage_ImageIdIsInvalid(t *testing.T) {
	app := App{
		imgService: nil,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	expectedBody := `{"title":"image not found","status":404,"detail":"","instance":"/api/image/#"}`

	req := httptest.NewRequest(http.MethodGet, "/api/image/#", nil)
	w := httptest.NewRecorder()
	getImageHandler := app.GetImage()
	getImageHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 404 {
		t.Errorf("expected status code: 404. Code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("expected Content-Type: application/problem+json. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedBody {
		t.Errorf("unexpected body content: %s", string(body))
	}
}

func TestApp_Upload(t *testing.T) {
	mock := &services.ImageServiceMock{}
	mock.UploadImageMock = func(img *domain.ImageFile) (string, *domain.Error) {
		return "12345fga34agsf23143fasfa", nil
	}

	app := App{
		imgService: mock,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	reqBody, writer, err := buildMultipartFormRequestBody()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	req := httptest.NewRequest(http.MethodPost, "/api/image", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	uploadHandler := app.Upload()
	uploadHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := `{"image_id":"12345fga34agsf23143fasfa"}`

	if resp.StatusCode != 201 {
		t.Errorf("expected status code: 201. Code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type: application/json. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedBody {
		t.Errorf("unexpected body content: %s", string(body))
	}
}

func TestApp_Upload_ErrorSavingImage(t *testing.T) {
	mock := &services.ImageServiceMock{}
	mock.UploadImageMock = func(img *domain.ImageFile) (string, *domain.Error) {
		return "", &domain.Error{
			Message: "an error occurred while trying to save the image in storage",
			Code:    500,
			Name:    domain.ServerErr,
			Error:   errors.New("not enough space"),
		}
	}

	app := App{
		imgService: mock,
		settings:   &config.Settings{
			Host: "127.0.0.1",
			Port: "4000",
		},
	}

	reqBody, writer, err := buildMultipartFormRequestBody()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	req := httptest.NewRequest(http.MethodPost, "/api/image", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	uploadHandler := app.Upload()
	uploadHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expectedBody := ``

	if resp.StatusCode != 500 {
		t.Errorf("expected status code: 500. Code received: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("expected Content-Type: text/plain. Content-Type received: %s", resp.Header.Get("Content-Type"))
	}

	if string(body) != expectedBody {
		t.Errorf("unexpected body content: %s", string(body))
	}
}





func buildMultipartFormRequestBody() (*bytes.Buffer, *multipart.Writer, error) {
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)

	//setup multipart file. Header is necessary.
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "upload_file", "image.png"))
	h.Set("Content-Type", "image/png") //multipart file contains the correct type

	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}
	writer.Close()

	return reqBody, writer, nil
}