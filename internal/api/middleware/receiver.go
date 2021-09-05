package middleware

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/models/responses"
	"log"
	"net/http"
)

//CheckReceivedContent check the existence of a file to upload
func CheckReceivedContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))
			notReceivedImgErr := responses.ErrProblem{
				Title:    "image not received",
				Detail:   "An image must be provided.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(notReceivedImgErr.Status)

			res, err := json.Marshal(notReceivedImgErr)
			if err != nil {
				// TODO handle error
			}
			w.Write(res) //todo handle error
			return
		}
		defer file.Close()

		next.ServeHTTP(w, r)
	})
}
