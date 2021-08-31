package responses

//UploadResponse represents the response body sent when an image is uploaded
//succesfully.
type UploadResponse struct {
	ImageId string `json:"image_id"`
}

type GetImageResponse struct {
	Image string `json:"image"`
}
