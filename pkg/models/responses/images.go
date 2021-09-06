package responses

//UploadResponse represents the response body sent when an image is uploaded
//successfully.
type UploadResponse struct {
	ImageId string `json:"image_id"`
}
